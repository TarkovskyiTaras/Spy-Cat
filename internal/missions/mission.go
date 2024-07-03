package missions

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type MissionService struct {
	missionRepo MissionRepository
	targetRepo  TargetRepository
}

type MissionRepository interface {
	DB() *sql.DB
	CreateMission(ctx context.Context, tx *sql.Tx, mission *Mission) error
	GetMissionByID(ctx context.Context, id uuid.UUID) (*Mission, error)
	GetAllMissions(ctx context.Context) ([]Mission, error)
	UpdateMission(ctx context.Context, id uuid.UUID, input MissionUpdateDTO) (*Mission, error)
	DeleteMission(ctx context.Context, id uuid.UUID) error
}

type TargetRepository interface {
	CreateTarget(ctx context.Context, tx *sql.Tx, target *Target) error
	GetTargetByID(ctx context.Context, id uuid.UUID) (*Target, error)
	GetTargetsByMissionID(ctx context.Context, missionID uuid.UUID) ([]Target, error)
	UpdateTarget(ctx context.Context, id uuid.UUID, input TargetUpdateDTO) (*Target, error)
	DeleteTarget(ctx context.Context, id uuid.UUID) error
}

func NewMissionService(missionRepo MissionRepository, targetRepo TargetRepository) *MissionService {
	return &MissionService{missionRepo: missionRepo, targetRepo: targetRepo}
}

func (s *MissionService) CreateMission(ctx context.Context, input *MissionCreateDTO) (*Mission, error) {
	tx, err := s.missionRepo.DB().Begin()
	if err != nil {
		return nil, err
	}

	missionID := uuid.New()
	var targets []Target
	for _, t := range input.Targets {
		target := &Target{
			ID:        uuid.New(),
			MissionID: missionID,
			Name:      t.Name,
			Country:   t.Country,
			Notes:     t.Notes,
			Complete:  false,
		}

		targets = append(targets, *target)
	}

	mission := &Mission{
		ID:       missionID,
		Name:     input.Name,
		CatID:    input.CatID,
		Complete: false,
		Targets:  targets,
	}

	if err := s.missionRepo.CreateMission(ctx, tx, mission); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, target := range mission.Targets {
		target := target
		if err = s.targetRepo.CreateTarget(ctx, tx, &target); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return mission, nil
}

func (s *MissionService) GetMissionByID(ctx context.Context, id uuid.UUID) (*Mission, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid mission ID")
	}

	mission, err := s.missionRepo.GetMissionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mission, nil
}

func (s *MissionService) GetAllMissions(ctx context.Context) ([]Mission, error) {
	missions, err := s.missionRepo.GetAllMissions(ctx)
	if err != nil {
		return nil, err
	}

	return missions, nil
}

func (s *MissionService) UpdateMission(ctx context.Context, id uuid.UUID, input MissionUpdateDTO) (*Mission, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid mission ID")
	}

	mission, err := s.missionRepo.UpdateMission(ctx, id, input)
	if err != nil {
		return nil, err
	}

	return mission, nil
}

func (s *MissionService) DeleteMission(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid mission ID")
	}

	return s.missionRepo.DeleteMission(ctx, id)
}

func (s *MissionService) AddTarget(ctx context.Context, missionID uuid.UUID, targetDTO *TargetCreateDTO) (*Target, error) {
	tx, err := s.missionRepo.DB().Begin()
	if err != nil {
		return nil, err
	}

	target := &Target{
		ID:        uuid.New(),
		MissionID: missionID,
		Name:      targetDTO.Name,
		Country:   targetDTO.Country,
		Notes:     targetDTO.Notes,
		Complete:  false,
	}

	if err := s.targetRepo.CreateTarget(ctx, tx, target); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return target, nil
}

func (s *MissionService) GetTargetsByMissionID(ctx context.Context, missionID uuid.UUID) ([]Target, error) {
	return s.targetRepo.GetTargetsByMissionID(ctx, missionID)
}

func (s *MissionService) DeleteTarget(ctx context.Context, id uuid.UUID) error {
	if err := s.targetRepo.DeleteTarget(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *MissionService) GetTargetByID(ctx context.Context, id uuid.UUID) (*Target, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid target ID")
	}
	return s.targetRepo.GetTargetByID(ctx, id)
}

func (s *MissionService) UpdateTarget(ctx context.Context, id uuid.UUID, input TargetUpdateDTO) (*Target, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid target ID")
	}

	target, err := s.targetRepo.UpdateTarget(ctx, id, input)
	if err != nil {
		return nil, err
	}

	return target, nil
}
