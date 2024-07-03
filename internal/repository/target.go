package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	"spycat/internal/missions"
	"strings"
)

type TargetRepository struct {
	db *sql.DB
}

func NewTargetRepository(db *sql.DB) *TargetRepository {
	return &TargetRepository{db: db}
}

func (repo *TargetRepository) CreateTarget(ctx context.Context, tx *sql.Tx, target *missions.Target) error {
	query := `INSERT INTO targets (id, mission_id, name, country, notes, complete) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := tx.ExecContext(ctx, query, target.ID, target.MissionID, target.Name, target.Country, target.Notes, target.Complete)
	return err
}

func (repo *MissionRepository) GetTargetsByMissionID(ctx context.Context, tx *sql.Tx, missionID uuid.UUID) ([]missions.Target, error) {
	query := `SELECT id, mission_id, name, country, notes, complete FROM targets WHERE mission_id = $1`
	rows, err := tx.QueryContext(ctx, query, missionID)
	if err != nil {
		log.Printf("Failed to execute query for targets: %v", err)
		return nil, err
	}
	defer rows.Close()

	var missionTargets []missions.Target
	for rows.Next() {
		target := missions.Target{}
		err := rows.Scan(&target.ID, &target.MissionID, &target.Name, &target.Country, &target.Notes, &target.Complete)
		if err != nil {
			log.Printf("Failed to scan target row: %v", err)
			return nil, err
		}
		missionTargets = append(missionTargets, target)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Row iteration error for targets: %v", err)
		return nil, err
	}

	return missionTargets, nil
}

func (repo *TargetRepository) GetTargetByID(ctx context.Context, id uuid.UUID) (*missions.Target, error) {
	query := `SELECT id, mission_id, name, country, notes, complete FROM targets WHERE id = $1`
	row := repo.db.QueryRowContext(ctx, query, id)

	var target missions.Target
	err := row.Scan(&target.ID, &target.MissionID, &target.Name, &target.Country, &target.Notes, &target.Complete)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("target not found")
		}
		return nil, err
	}

	return &target, nil
}

func (repo *TargetRepository) UpdateTarget(ctx context.Context, id uuid.UUID, input missions.TargetUpdateDTO) (*missions.Target, error) {
	var setClauses []string
	var args []interface{}
	argID := 1

	if input.Notes != nil {
		setClauses = append(setClauses, fmt.Sprintf("notes = $%d", argID))
		args = append(args, *input.Notes)
		argID++
	}
	if input.Complete != nil {
		setClauses = append(setClauses, fmt.Sprintf("complete = $%d", argID))
		args = append(args, *input.Complete)
		argID++
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no set clauses found")
	}

	args = append(args, id)
	setClause := strings.Join(setClauses, ", ")
	query := fmt.Sprintf("UPDATE targets SET %s WHERE id = $%d", setClause, argID)

	_, err := repo.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	target, err := repo.GetTargetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return target, nil
}

func (repo *TargetRepository) GetTargetsByMissionID(ctx context.Context, missionID uuid.UUID) ([]missions.Target, error) {
	query := `SELECT id, mission_id, name, country, notes, complete FROM targets WHERE mission_id = $1`
	rows, err := repo.db.QueryContext(ctx, query, missionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var missionTargets []missions.Target
	for rows.Next() {
		target := missions.Target{}
		err := rows.Scan(&target.ID, &target.MissionID, &target.Name, &target.Country, &target.Notes, &target.Complete)
		if err != nil {
			return nil, err
		}
		missionTargets = append(missionTargets, target)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return missionTargets, nil
}

func (repo *TargetRepository) DeleteTarget(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.Exec("DELETE FROM targets WHERE id = $1", id)
	return err
}
