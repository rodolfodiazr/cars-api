package services

import (
	"cars/models"
	"cars/repositories"
	"strings"
)

// CarService defines available operations for managing cars.
type CarService interface {
	Find(id string) (models.Car, error)
	List(filters models.CarFilters) (models.Cars, error)
	Create(car *models.Car) error
	Update(car *models.Car) error
}

// DefaultCarService is the default implementation of CarService.
type DefaultCarService struct {
	repo repositories.CarRepository
}

// NewCarService creates a new instance of DefaultCarService.
func NewCarService(repo repositories.CarRepository) CarService {
	return &DefaultCarService{
		repo: repo,
	}
}

// Find retrieves a car by its ID if it exists.
func (s *DefaultCarService) Find(id string) (models.Car, error) {
	return s.repo.Find(id)
}

// List retrieves all available cars.
func (s *DefaultCarService) List(f models.CarFilters) (models.Cars, error) {
	allCars, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	var filtered models.Cars
	for _, car := range allCars {
		if f.Make != "" && !strings.EqualFold(car.Make, f.Make) {
			continue
		}
		if f.Model != "" && !strings.EqualFold(car.Model, f.Model) {
			continue
		}
		if f.Year != 0 && car.Year != f.Year {
			continue
		}
		filtered = append(filtered, car)
	}

	return filtered, nil
}

// Create adds a new car to the repository.
func (s *DefaultCarService) Create(car *models.Car) error {
	return s.repo.Create(car)
}

// Update updates a car in the repository.
func (s *DefaultCarService) Update(car *models.Car) error {
	return s.repo.Update(car)
}
