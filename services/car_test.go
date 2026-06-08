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

func TestDefaultCarService_Create(t *testing.T) {
	t.Run("should create car when validation succeeds", func(t *testing.T) {
		// Arrange
		car := &models.Car{
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Gray",
			Category: "Sedan",
			Year:     2026,
		}

		repo := &MockCarRepository{
			CreateFn: func(c *models.Car) error {
				if c != car {
					t.Fatal("expected same car pointer to be passed to repository")
				}
				return nil
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Create(car)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("should return validation error when car is invalid", func(t *testing.T) {
		// Arrange
		car := &models.Car{}

		repo := &MockCarRepository{
			CreateFn: func(c *models.Car) error {
				t.Fatal("repository Create should not be called")
				return nil
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Create(car)

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	t.Run("should return internal error when repository fails", func(t *testing.T) {
		// Arrange
		car := &models.Car{
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Gray",
			Category: "Sedan",
			Year:     2026,
		}

		expectedErr := errors.New("database unavailable")

		repo := &MockCarRepository{
			CreateFn: func(c *models.Car) error {
				return expectedErr
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Create(car)

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}

		var serviceError *e.ServiceError
		if !errors.As(err, &serviceError) {
			t.Fatalf("expected ServiceError, got %T", err)
		}

		if serviceError.Code != e.CodeInternalError {
			t.Fatalf("expected INTERNAL_ERROR, got %v", serviceError.Code)
		}
	})
}

func TestDefaultCarService_Update(t *testing.T) {
	t.Run("should update car when validation succeeds", func(t *testing.T) {
		// Arrange
		car := &models.Car{
			ID:       "1",
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Gray",
			Category: "Sedan",
			Year:     2026,
		}

		repo := &MockCarRepository{
			UpdateFn: func(c *models.Car) error {
				if c != car {
					t.Fatal("expected same car pointer to be passed to repository")
				}
				return nil
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Update(car)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("should return validation error when car is invalid", func(t *testing.T) {
		// Arrange
		car := &models.Car{}

		repo := &MockCarRepository{
			UpdateFn: func(c *models.Car) error {
				t.Fatal("repository Update should not be called")
				return nil
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Update(car)

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}

		var serviceError *e.ServiceError
		if !errors.As(err, &serviceError) {
			t.Fatalf("expected ServiceError, got %T", err)
		}

		if serviceError.Code != e.CodeValidationFailed {
			t.Fatalf("expected VALIDATION_FAILED, got %v", serviceError.Code)
		}
	})

	t.Run("should return car not found error when repository returns ErrCarNotFound", func(t *testing.T) {
		// Arrange
		car := &models.Car{
			ID:       "1",
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Gray",
			Category: "Sedan",
			Year:     2026,
		}

		repo := &MockCarRepository{
			UpdateFn: func(c *models.Car) error {
				return e.ErrCarNotFound
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Update(car)

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}

		var serviceError *e.ServiceError
		if !errors.As(err, &serviceError) {
			t.Fatalf("expected ServiceError, got %T", err)
		}

		if serviceError.Code != e.CodeCarNotFound {
			t.Fatalf("expected CAR_NOT_FOUND, got %v", serviceError.Code)
		}
	})

	t.Run("should return internal error when repository fails unexpectedly", func(t *testing.T) {
		// Arrange
		car := &models.Car{
			ID:       "1",
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Gray",
			Category: "Sedan",
			Year:     2026,
		}

		expectedErr := errors.New("database unavailable")

		repo := &MockCarRepository{
			UpdateFn: func(c *models.Car) error {
				return expectedErr
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Update(car)

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}

		var serviceError *e.ServiceError
		if !errors.As(err, &serviceError) {
			t.Fatalf("expected ServiceError, got %T", err)
		}

		if serviceError.Code != e.CodeInternalError {
			t.Fatalf("expected INTERNAL_ERROR, got %v", serviceError.Code)
		}
	})
}

func TestDefaultCarService_Delete(t *testing.T) {
	t.Run("should delete car when repository succeeds", func(t *testing.T) {
		// Arrange
		expectedID := "1"

		repo := &MockCarRepository{
			DeleteFn: func(id string) error {
				if id != expectedID {
					t.Fatalf("expected id %q, got %q", expectedID, id)
				}
				return nil
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Delete(expectedID)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("should return car not found error when repository returns ErrCarNotFound", func(t *testing.T) {
		// Arrange
		repo := &MockCarRepository{
			DeleteFn: func(id string) error {
				return e.ErrCarNotFound
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Delete("missing-id")

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	t.Run("should return internal error when repository fails unexpectedly", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database unavailable")

		repo := &MockCarRepository{
			DeleteFn: func(id string) error {
				return expectedErr
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Delete("1")

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}

		var serviceError *e.ServiceError
		if !errors.As(err, &serviceError) {
			t.Fatalf("expected ServiceError, got %T", err)
		}

		if serviceError.Code != e.CodeInternalError {
			t.Fatalf("expected INTERNAL_ERROR, got %v", serviceError.Code)
		}
	})
}
