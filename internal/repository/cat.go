package repository

import (
	"database/sql"
	"spycat/internal/cats"

	"github.com/google/uuid"
)

type CatsRepository struct {
	db *sql.DB
}

func NewCatsRepository(db *sql.DB) *CatsRepository {
	return &CatsRepository{db: db}
}

func (r *CatsRepository) CreateCat(cat *cats.Cat) (*cats.Cat, error) {
	query := `INSERT INTO cats (id, name, years_of_experience, breed, salary) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, cat.ID, cat.Name, cat.YearsOfExperience, cat.Breed, cat.Salary)
	if err != nil {
		return nil, err
	}

	cat, err = r.GetCatByID(cat.ID)
	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (r *CatsRepository) GetAllCats() ([]cats.Cat, error) {
	var allCats []cats.Cat
	query := `SELECT id, name, years_of_experience, breed, salary FROM cats`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cat cats.Cat
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary); err != nil {
			return nil, err
		}
		allCats = append(allCats, cat)
	}
	return allCats, nil
}

func (r *CatsRepository) GetCatByID(id uuid.UUID) (*cats.Cat, error) {
	var cat cats.Cat
	query := `SELECT id, name, years_of_experience, breed, salary FROM cats WHERE id = $1`
	row := r.db.QueryRow(query, id)
	err := row.Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary)
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CatsRepository) UpdateCat(id uuid.UUID, input *cats.CatUpdateDTO) (*cats.Cat, error) {
	query := "UPDATE cats SET salary = $1 WHERE id = $2"

	_, err := r.db.Exec(query, input.Salary, id)
	if err != nil {
		return nil, err
	}

	cat, err := r.GetCatByID(id)
	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (r *CatsRepository) DeleteCat(id uuid.UUID) error {
	query := `DELETE FROM cats WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
