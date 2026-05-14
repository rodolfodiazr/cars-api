package services

import (
	"cars/models"
	e "cars/pkg/errors"
	"errors"
	"reflect"
	"testing"
)

type FindFunc func(id string) (models.Car, error)
type ListFunc func(filters models.CarFilters) (models.Cars, error)
type CreateFunc func(car *models.Car) error
type UpdateFunc func(car *models.Car) error
type DeleteFunc func(id string) error

type MockCarRepository struct {
	FindFn   FindFunc
	ListFn   ListFunc
	CreateFn CreateFunc
	UpdateFn UpdateFunc
	DeleteFn DeleteFunc
}

func (m *MockCarRepository) Find(id string) (models.Car, error) {
	return m.FindFn(id)
}

func (m *MockCarRepository) List(filters models.CarFilters) (models.Cars, error) {
	return m.ListFn(filters)
}

func (m *MockCarRepository) Create(car *models.Car) error {
	return m.CreateFn(car)
}

func (m *MockCarRepository) Update(car *models.Car) error {
	return m.UpdateFn(car)
}

func (m *MockCarRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

func TestDefaultCarService_Find(t *testing.T) {
	t.Run("should return car when repository finds it", func(t *testing.T) {
		// Arrange
		expected := models.Car{
			ID:       "1",
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Black",
			Category: "Sedan",
			Year:     2026,
		}

		repo := &MockCarRepository{
			FindFn: func(id string) (models.Car, error) {
				return expected, nil
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		got, err := service.Find(expected.ID)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != expected {
			t.Fatalf("expected %+v, got %+v", expected, got)
		}
	})

	t.Run("should return car not found error when repository returns ErrCarNotFound", func(t *testing.T) {
		// Arrange
		repo := &MockCarRepository{
			FindFn: func(id string) (models.Car, error) {
				return models.Car{}, e.ErrCarNotFound
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		_, err := service.Find("missing-id")

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	t.Run("should return internal error for unexpected repository errors", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database unavailable")

		repo := &MockCarRepository{
			FindFn: func(id string) (models.Car, error) {
				return models.Car{}, expectedErr
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		_, err := service.Find("1")

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}

func TestDefaultCarService_List(t *testing.T) {
	t.Run("should return cars when repository succeeds", func(t *testing.T) {
		// Arrange
		expected := models.Cars{
			{
				ID:    "1",
				Make:  "Toyota",
				Model: "Corolla",
			},
			{
				ID:    "2",
				Make:  "Honda",
				Model: "Civic",
			},
		}

		filters := models.CarFilters{
			Make: "Toyota",
		}

		repo := &MockCarRepository{
			ListFn: func(f models.CarFilters) (models.Cars, error) {
				if f != filters {
					t.Fatalf("expected filters %+v, got %+v", filters, f)
				}

				return expected, nil
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		got, err := service.List(filters)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %+v, got %+v", expected, got)
		}
	})

	t.Run("should return internal error when repository fails", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database unavailable")

		repo := &MockCarRepository{
			ListFn: func(f models.CarFilters) (models.Cars, error) {
				return models.Cars{}, expectedErr
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		_, err := service.List(models.CarFilters{})

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})
}
