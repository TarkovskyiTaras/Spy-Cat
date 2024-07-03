package missions

import (
	"github.com/google/uuid"
)

type Mission struct {
	ID       uuid.UUID
	Name     string
	CatID    uuid.UUID
	Complete bool
	Targets  []Target
}

type MissionCreateDTO struct {
	Name    string            `json:"name"`
	CatID   uuid.UUID         `json:"cat_id"`
	Targets []TargetCreateDTO `json:"targets"`
}

type MissionUpdateDTO struct {
	Name     *string    `json:"name"`
	CatID    *uuid.UUID `json:"cat_id"`
	Targets  *[]Target  `json:"targets"`
	Complete *bool      `json:"complete"`
}

type Target struct {
	ID        uuid.UUID
	MissionID uuid.UUID
	Name      string
	Country   string
	Notes     string
	Complete  bool
}

type TargetCreateDTO struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Notes   string `json:"notes"`
}

type TargetUpdateDTO struct {
	Notes    *string `json:"notes"`
	Complete *bool   `json:"complete"`
}
