package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"

	"spycat/internal/missions"
)

type MissionRepository struct {
	db *sql.DB
}

func NewMissionRepository(db *sql.DB) *MissionRepository {
	return &MissionRepository{db: db}
}

func (repo *MissionRepository) DB() *sql.DB {
	return repo.db
}

func (repo *MissionRepository) CreateMission(ctx context.Context, tx *sql.Tx, mission *missions.Mission) error {
	query := `INSERT INTO missions (id, name, cat_id, complete) VALUES ($1, $2, $3, $4)`
	_, err := tx.ExecContext(ctx, query, mission.ID, mission.Name, mission.CatID, mission.Complete)
	return err
}

func (repo *MissionRepository) GetMissionByID(ctx context.Context, id uuid.UUID) (*missions.Mission, error) {
	var mission missions.Mission
	query := "SELECT id, cat_id, complete FROM missions WHERE id = $1"

	err := repo.db.QueryRowContext(ctx, query, id).Scan(
		&mission.ID,
		&mission.CatID,
		&mission.Complete,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No mission found with id: %s", id)
			return nil, nil
		}
		log.Printf("Error querying mission by id: %s, error: %v", id, err)
		return nil, err
	}

	return &mission, nil
}

func (repo *MissionRepository) GetAllMissions(ctx context.Context) ([]missions.Mission, error) {
	query := `
	SELECT 
		m.id as mission_id, 
		m.name as mission_name, 
		m.cat_id as mission_cat_id, 
		m.complete as mission_complete, 
		t.id as target_id, 
		t.mission_id as target_mission_id, 
		t.name as target_name, 
		t.country as target_country, 
		t.notes as target_notes, 
		t.complete as target_complete 
	FROM missions m
	LEFT JOIN targets t ON m.id = t.mission_id
	`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	missionMap := make(map[uuid.UUID]*missions.Mission)

	for rows.Next() {
		var missionID uuid.UUID
		var missionName string
		var missionCatID uuid.UUID
		var missionComplete bool

		var targetID sql.NullString
		var targetMissionID sql.NullString
		var targetName sql.NullString
		var targetCountry sql.NullString
		var targetNotes sql.NullString
		var targetComplete sql.NullBool

		err := rows.Scan(
			&missionID,
			&missionName,
			&missionCatID,
			&missionComplete,
			&targetID,
			&targetMissionID,
			&targetName,
			&targetCountry,
			&targetNotes,
			&targetComplete,
		)
		if err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}

		mission, exists := missionMap[missionID]
		if !exists {
			mission = &missions.Mission{
				ID:       missionID,
				Name:     missionName,
				CatID:    missionCatID,
				Complete: missionComplete,
			}
			missionMap[missionID] = mission
		}

		if targetID.Valid {
			target := missions.Target{
				ID:        uuid.MustParse(targetID.String),
				MissionID: uuid.MustParse(targetMissionID.String),
				Name:      targetName.String,
				Country:   targetCountry.String,
				Notes:     targetNotes.String,
				Complete:  targetComplete.Bool,
			}
			mission.Targets = append(mission.Targets, target)
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row iteration error: %v", err)
		return nil, err
	}

	var allMissions []missions.Mission
	for _, mission := range missionMap {
		allMissions = append(allMissions, *mission)
	}

	return allMissions, nil
}

func (repo *MissionRepository) UpdateMission(ctx context.Context, id uuid.UUID, input missions.MissionUpdateDTO) (*missions.Mission, error) {
	var setClauses []string
	var args []interface{}
	argID := 1

	if input.CatID != nil {
		setClauses = append(setClauses, fmt.Sprintf("cat_id = $%d", argID))
		args = append(args, *input.CatID)
		argID++
	}
	if input.Complete != nil {
		setClauses = append(setClauses, fmt.Sprintf("complete = $%d", argID))
		args = append(args, *input.Complete)
		argID++
	}
	if input.Name != nil {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argID))
		args = append(args, *input.Name)
		argID++
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no set clauses found")
	}

	args = append(args, id)
	setClause := strings.Join(setClauses, ", ")
	query := fmt.Sprintf("UPDATE missions SET %s WHERE id = $%d", setClause, argID)

	_, err := repo.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	mission, err := repo.GetMissionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mission, nil
}

func (repo *MissionRepository) DeleteMission(ctx context.Context, id uuid.UUID) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM targets WHERE mission_id = $1", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM missions WHERE id = $1", id)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
