package cats

import "github.com/google/uuid"

type Cat struct {
	ID                uuid.UUID
	Name              string
	YearsOfExperience int
	Breed             string
	Salary            float64
}

type CatCreateDTO struct {
	Name              string  `json:"name" binding:"required"`
	YearsOfExperience int     `json:"years_of_experience" binding:"required"`
	Breed             string  `json:"breed" binding:"required"`
	Salary            float64 `json:"salary" binding:"required"`
}

type CatUpdateDTO struct {
	Salary *float64 `json:"salary"`
}
