package services

import (
	c "cars/car"
	e "cars/pkg/errors"
	"cars/repositories"
	"errors"
)

// CarService defines available operations for managing cars.
type CarService interface {
	Find(id string) (c.Car, error)
	List(filters c.CarFilters) (c.Cars, error)
	Create(car *c.Car) error
	Update(car *c.Car) error
	PartialUpdate(params c.CarUpdateParams) (c.Car, error)
	Delete(id string) error
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
func (s *DefaultCarService) Find(id string) (c.Car, error) {
	car, err := s.repo.Find(id)
	if err != nil {
		if errors.Is(err, e.ErrCarNotFound) {
			return c.Car{}, e.NewCarNotFoundError(err)
		}

		return c.Car{}, e.NewInternalError(err)
	}
	return car, nil
}

// List retrieves all available cars.
func (s *DefaultCarService) List(f c.CarFilters) (c.Cars, error) {
	cars, err := s.repo.List(f)
	if err != nil {
		return c.Cars{}, e.NewInternalError(err)
	}
	return cars, nil
}

// Create adds a new car to the repository.
// The car must contain all required fields.
func (s *DefaultCarService) Create(car *c.Car) error {
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
func (s *DefaultCarService) Update(car *c.Car) error {
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

// PartialUpdate ...
func (s *DefaultCarService) PartialUpdate(params c.CarUpdateParams) (c.Car, error) {
	if err := params.Validate(); err != nil {
		return c.Car{}, e.NewValidationError(err)
	}

	car, err := s.repo.PartialUpdate(params)
	if err != nil {
		if errors.Is(err, e.ErrCarNotFound) {
			return c.Car{}, e.NewCarNotFoundError(err)
		}
		return c.Car{}, e.NewInternalError(err)
	}
	return car, nil
}

// Delete removes a car identified by the given ID.
func (s *DefaultCarService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		if errors.Is(err, e.ErrCarNotFound) {
			return e.NewCarNotFoundError(err)
		}
		return e.NewInternalError(err)
	}
	return nil
}
