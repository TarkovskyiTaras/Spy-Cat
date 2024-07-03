package cats

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type CatService struct {
	repo catsRepository
}

type catsRepository interface {
	CreateCat(cat *Cat) (*Cat, error)
	GetAllCats() ([]Cat, error)
	GetCatByID(id uuid.UUID) (*Cat, error)
	UpdateCat(id uuid.UUID, input *CatUpdateDTO) (*Cat, error)
	DeleteCat(id uuid.UUID) error
}

func NewCatService(repo catsRepository) *CatService {
	return &CatService{repo: repo}
}

func (s *CatService) CreateCat(input *CatCreateDTO) (*Cat, error) {
	cat := &Cat{
		ID:                uuid.New(),
		Name:              input.Name,
		YearsOfExperience: input.YearsOfExperience,
		Breed:             input.Breed,
		Salary:            input.Salary,
	}

	cat, err := s.repo.CreateCat(cat)
	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *CatService) GetAllCats() ([]Cat, error) {
	return s.repo.GetAllCats()
}

func (s *CatService) GetCatByID(id uuid.UUID) (*Cat, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid ID")
	}

	return s.repo.GetCatByID(id)
}

func (s *CatService) UpdateCat(id uuid.UUID, input *CatUpdateDTO) (*Cat, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid ID")
	}

	cat, err := s.repo.UpdateCat(id, input)
	if err != nil {
		return nil, err
	}

	return cat, nil
}

func (s *CatService) DeleteCat(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid cats ID")
	}

	return s.repo.DeleteCat(id)
}

type CatBreed struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func FetchValidBreeds() ([]CatBreed, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/breeds", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch breeds")
	}

	var breeds []CatBreed
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return nil, err
	}

	return breeds, nil
}

func IsValidBreed(breed string, breeds []CatBreed) bool {
	for _, b := range breeds {
		if b.Name == breed {
			return true
		}
	}
	return false
}
