package repositories

import (
	"cars/models"
	e "cars/pkg/errors"
	"cars/pkg/utils"
	"strings"
	"sync"
)

// CarRepository defines methods for managing car persistence.
type CarRepository interface {
	Find(id string) (models.Car, error)
	List(filters models.CarFilters) (models.Cars, error)
	Create(car *models.Car) error
	Update(car *models.Car) error
	Delete(id string) error
}

// DefaultCarRepository is an in-memory implementation of CarRepository.
type DefaultCarRepository struct {
	cars map[string]models.Car
	mu   sync.RWMutex
}

// NewCarRepository creates a new instance of DefaultCarRepository with initial data.
func NewCarRepository(initialData map[string]models.Car) CarRepository {
	if initialData == nil {
		initialData = make(map[string]models.Car)
	}

	repo := &DefaultCarRepository{
		cars: initialData,
	}
	return repo
}

// Find searches for a car by its ID. Returns an error if not found.
func (r *DefaultCarRepository) Find(id string) (models.Car, error) {
	r.mu.RLock()
	car, exists := r.cars[id]
	r.mu.RUnlock()

	if !exists {
		return models.Car{}, e.ErrCarNotFound
	}
	return car, nil
}

// List returns all stored cars.
func (r *DefaultCarRepository) List(f models.CarFilters) (models.Cars, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make(models.Cars, 0, len(r.cars))

	for _, car := range r.cars {
		if f.Make != "" && !strings.EqualFold(car.Make, f.Make) {
			continue
		}
		if f.Model != "" && !strings.EqualFold(car.Model, f.Model) {
			continue
		}
		if f.Year != 0 && car.Year != f.Year {
			continue
		}
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

	r.mu.Lock()
	defer r.mu.Unlock()

	car.ID = id
	r.cars[car.ID] = *car

	return nil
}

// Update updates an existing car in the repository. Returns an error if not found.
func (r *DefaultCarRepository) Update(car *models.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.cars[car.ID]; !exists {
		return e.ErrCarNotFound
	}

	r.cars[car.ID] = *car
	return nil
}

// Delete removes a car identified by the given id from the repository.
func (r *DefaultCarRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.cars[id]; !ok {
		return e.ErrCarNotFound
	}

	delete(r.cars, id)
	return nil
}
