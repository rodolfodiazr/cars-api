package services

import (
	"cars/models"
	e "cars/pkg/errors"
	"cars/repositories"
	"errors"
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
	car, err := s.repo.Find(id)
	if err != nil {
		if errors.Is(err, e.ErrCarNotFound) {
			return models.Car{}, e.NewCarNotFoundError(err)
		}

		return models.Car{}, e.NewInternalError(err)
	}
	return car, nil
}

// List retrieves all available cars.
func (s *DefaultCarService) List(f models.CarFilters) (models.Cars, error) {
	cars, err := s.repo.List(f)
	if err != nil {
		return models.Cars{}, e.NewInternalError(err)
	}
	return cars, nil
}

// Create adds a new car to the repository.
func (s *DefaultCarService) Create(car *models.Car) error {
	if err := car.ValidateForCreate(); err != nil {
		return e.NewValidationError(err)
	}

	if err := s.repo.Create(car); err != nil {
		return e.NewInternalError(err)
	}
	return nil
}

// Update replaces an existing car with the provided data.
// The car must contain all required fields and a valid ID.
func (s *DefaultCarService) Update(car *models.Car) error {
	if err := car.ValidateForUpdate(); err != nil {
		return e.NewValidationError(err)
	}

	if err := s.repo.Update(car); err != nil {
		if errors.Is(err, e.ErrCarNotFound) {
			return e.NewCarNotFoundError(err)
		}
		return e.NewInternalError(err)
	}
	return nil
}
