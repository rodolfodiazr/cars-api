package repositories

import (
	"cars/models"
	e "cars/pkg/errors"
	"cars/pkg/utils"
)

// CarRepository defines methods for managing car persistence.
type CarRepository interface {
	Find(id string) (models.Car, error)
	List() (models.Cars, error)
	Create(car *models.Car) error
	Update(car *models.Car) error
}

// DefaultCarRepository is an in-memory implementation of CarRepository.
type DefaultCarRepository struct {
	cars map[string]models.Car
}

// NewCarRepository creates a new instance of DefaultCarRepository with initial data.
func NewCarRepository(initialData map[string]models.Car) CarRepository {
	repo := &DefaultCarRepository{
		cars: initialData,
	}
	return repo
}

// Find searches for a car by its ID. Returns an error if not found.
func (r *DefaultCarRepository) Find(id string) (models.Car, error) {
	car, exists := r.cars[id]
	if !exists {
		return models.Car{}, e.ErrCarNotFound
	}
	return car, nil
}

// List returns all stored cars.
func (r *DefaultCarRepository) List() (models.Cars, error) {
	var list models.Cars
	for _, car := range r.cars {
		list = append(list, car)
	}
	return list, nil
}

// Create stores a new car in the repository.
func (r *DefaultCarRepository) Create(car *models.Car) error {
	id, err := utils.GenerateID()
	if err != nil {
		return err
	}

	car.ID = id
	r.cars[car.ID] = *car
	return nil
}

// Update updates an existing car in the repository. Returns an error if not found.
func (r *DefaultCarRepository) Update(car *models.Car) error {
	if _, exists := r.cars[car.ID]; !exists {
		return e.ErrCarNotFound
	}

	r.cars[car.ID] = *car
	return nil
}
