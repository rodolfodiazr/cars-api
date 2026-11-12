package services

import (
	c "cars/car"
	e "cars/pkg/errors"
	"errors"
	"net/http"
	"reflect"
	"testing"
)

type FindFunc func(id string) (c.Car, error)
type ListFunc func(filters c.CarFilters) (c.Cars, error)
type CreateFunc func(car *c.Car) error
type UpdateFunc func(car *c.Car) error
type PartialUpdateFunc func(params c.CarUpdateParams) (c.Car, error)
type DeleteFunc func(id string) error

type MockCarRepository struct {
	FindFn          FindFunc
	ListFn          ListFunc
	CreateFn        CreateFunc
	UpdateFn        UpdateFunc
	PartialUpdateFn PartialUpdateFunc
	DeleteFn        DeleteFunc
}

func (m *MockCarRepository) Find(id string) (c.Car, error) {
	return m.FindFn(id)
}

func (m *MockCarRepository) List(filters c.CarFilters) (c.Cars, error) {
	return m.ListFn(filters)
}

func (m *MockCarRepository) Create(car *c.Car) error {
	return m.CreateFn(car)
}

func (m *MockCarRepository) Update(car *c.Car) error {
	return m.UpdateFn(car)
}

func (m *MockCarRepository) PartialUpdate(params c.CarUpdateParams) (c.Car, error) {
	return m.PartialUpdateFn(params)
}

func (m *MockCarRepository) Delete(id string) error {
	return m.DeleteFn(id)
}

func assertServiceError(t *testing.T, err error, expected *e.ServiceError) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error but got nil")
	}

	var got *e.ServiceError
	if !errors.As(err, &got) {
		t.Fatalf("expected ServiceError, got %T", err)
	}

	if got.Code != expected.Code {
		t.Errorf("expected code %q, got %q", expected.Code, got.Code)
	}

	if got.Message != expected.Message {
		t.Errorf("expected message %q, got %q", expected.Message, got.Message)
	}

	if got.StatusCode != expected.StatusCode {
		t.Errorf(
			"expected status code %d, got %d",
			expected.StatusCode,
			got.StatusCode,
		)
	}

	if !errors.Is(got.Err, expected.Err) {
		t.Errorf(
			"expected underlying error %v, got %v",
			expected.Err,
			got.Err,
		)
	}
}

func TestDefaultCarService_Find(t *testing.T) {
	t.Run("should return car when repository finds it", func(t *testing.T) {
		// Arrange
		expected := c.Car{
			ID:       "1",
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Black",
			Category: "Sedan",
			Year:     2026,
		}

		repo := &MockCarRepository{
			FindFn: func(id string) (c.Car, error) {
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
		expected := &e.ServiceError{
			Code:       e.CodeCarNotFound,
			Message:    "Car not found",
			StatusCode: http.StatusNotFound,
			Err:        e.ErrCarNotFound,
		}

		repo := &MockCarRepository{
			FindFn: func(id string) (c.Car, error) {
				return c.Car{}, e.ErrCarNotFound
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		_, err := service.Find("missing-id")

		// Assert
		assertServiceError(t, err, expected)
	})

	t.Run("should return internal error for unexpected repository errors", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database unavailable")

		expected := &e.ServiceError{
			Code:       e.CodeInternalError,
			Message:    "Internal server error",
			StatusCode: http.StatusInternalServerError,
			Err:        expectedErr,
		}

		repo := &MockCarRepository{
			FindFn: func(id string) (c.Car, error) {
				return c.Car{}, expectedErr
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		_, err := service.Find("1")

		// Assert
		assertServiceError(t, err, expected)
	})
}

func TestDefaultCarService_List(t *testing.T) {
	t.Run("should return cars when repository succeeds", func(t *testing.T) {
		// Arrange
		expected := c.Cars{
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

		filters := c.CarFilters{
			Make: "Toyota",
		}

		repo := &MockCarRepository{
			ListFn: func(f c.CarFilters) (c.Cars, error) {
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

		expected := &e.ServiceError{
			Code:       e.CodeInternalError,
			Message:    "Internal server error",
			StatusCode: http.StatusInternalServerError,
			Err:        expectedErr,
		}

		repo := &MockCarRepository{
			ListFn: func(f c.CarFilters) (c.Cars, error) {
				return nil, expectedErr
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		_, err := service.List(c.CarFilters{})

		// Assert
		assertServiceError(t, err, expected)
	})
}

func TestDefaultCarService_Create(t *testing.T) {
	t.Run("should create car when validation succeeds", func(t *testing.T) {
		// Arrange
		car := &c.Car{
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Gray",
			Category: "Sedan",
			Year:     2026,
		}

		repo := &MockCarRepository{
			CreateFn: func(c *c.Car) error {
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
		expected := &e.ServiceError{
			Code:       e.CodeValidationFailed,
			Message:    "Validation failed",
			StatusCode: http.StatusBadRequest,
			Err:        c.ErrCarMakeRequired,
		}

		car := &c.Car{}

		repo := &MockCarRepository{
			CreateFn: func(c *c.Car) error {
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
		assertServiceError(t, err, expected)
	})

	t.Run("should return internal error when repository fails", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database unavailable")

		expected := &e.ServiceError{
			Code:       e.CodeInternalError,
			Message:    "Internal server error",
			StatusCode: http.StatusInternalServerError,
			Err:        expectedErr,
		}

		car := &c.Car{
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Gray",
			Category: "Sedan",
			Year:     2026,
		}

		repo := &MockCarRepository{
			CreateFn: func(c *c.Car) error {
				return expectedErr
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Create(car)

		// Assert
		assertServiceError(t, err, expected)
	})
}

func TestDefaultCarService_Update(t *testing.T) {
	t.Run("should update car when validation succeeds", func(t *testing.T) {
		// Arrange
		car := &c.Car{
			ID:       "1",
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Gray",
			Category: "Sedan",
			Year:     2026,
		}

		repo := &MockCarRepository{
			UpdateFn: func(c *c.Car) error {
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
		expected := &e.ServiceError{
			Code:       e.CodeValidationFailed,
			Message:    "Validation failed",
			StatusCode: http.StatusBadRequest,
			Err:        c.ErrCarIDRequiredForUpdate,
		}

		car := &c.Car{}

		repo := &MockCarRepository{
			UpdateFn: func(c *c.Car) error {
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
		assertServiceError(t, err, expected)
	})

	t.Run("should return car not found error when repository returns ErrCarNotFound", func(t *testing.T) {
		// Arrange
		expected := &e.ServiceError{
			Code:       e.CodeCarNotFound,
			Message:    "Car not found",
			StatusCode: http.StatusNotFound,
			Err:        e.ErrCarNotFound,
		}

		car := &c.Car{
			ID:       "1",
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Gray",
			Category: "Sedan",
			Year:     2026,
		}

		repo := &MockCarRepository{
			UpdateFn: func(c *c.Car) error {
				return e.ErrCarNotFound
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Update(car)

		// Assert
		assertServiceError(t, err, expected)
	})

	t.Run("should return internal error when repository fails unexpectedly", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database unavailable")

		expected := &e.ServiceError{
			Code:       e.CodeInternalError,
			Message:    "Internal server error",
			StatusCode: http.StatusInternalServerError,
			Err:        expectedErr,
		}

		car := &c.Car{
			ID:       "1",
			Make:     "Toyota",
			Model:    "Corolla",
			Color:    "Gray",
			Category: "Sedan",
			Year:     2026,
		}

		repo := &MockCarRepository{
			UpdateFn: func(c *c.Car) error {
				return expectedErr
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		err := service.Update(car)

		// Assert
		assertServiceError(t, err, expected)
	})
}

func TestDefaultCarService_PartialUpdate(t *testing.T) {
	t.Run("should partial update car when validation succeeds", func(t *testing.T) {
		// Arrange
		params := c.CarUpdateParams{
			ID: "1",
		}

		repo := &MockCarRepository{
			PartialUpdateFn: func(params c.CarUpdateParams) (c.Car, error) {
				return c.Car{
					ID:       "1",
					Make:     "Toyota",
					Model:    "Camry",
					Color:    "White",
					Category: "Sedan",
					Year:     2020,
				}, nil
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		_, err := service.PartialUpdate(params)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("should return validation error when car is invalid", func(t *testing.T) {
		// Arrange
		params := c.CarUpdateParams{}

		expected := &e.ServiceError{
			Code:       e.CodeValidationFailed,
			Message:    "Validation failed",
			StatusCode: http.StatusBadRequest,
			Err:        c.ErrCarIDRequiredForUpdate,
		}

		repo := &MockCarRepository{
			PartialUpdateFn: func(params c.CarUpdateParams) (c.Car, error) {
				t.Fatal("repository Update should not be called")
				return c.Car{}, nil
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		_, err := service.PartialUpdate(params)

		// Assert
		assertServiceError(t, err, expected)
	})

	t.Run("should return car not found error when repository returns ErrCarNotFound", func(t *testing.T) {
		// Arrange
		params := c.CarUpdateParams{
			ID: "1",
		}

		expected := &e.ServiceError{
			Code:       e.CodeCarNotFound,
			Message:    "Car not found",
			StatusCode: http.StatusNotFound,
			Err:        e.ErrCarNotFound,
		}

		repo := &MockCarRepository{
			PartialUpdateFn: func(params c.CarUpdateParams) (c.Car, error) {
				return c.Car{}, e.ErrCarNotFound
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		_, err := service.PartialUpdate(params)

		// Assert
		assertServiceError(t, err, expected)
	})

	t.Run("should return internal error when repository fails unexpectedly", func(t *testing.T) {
		// Arrange
		params := c.CarUpdateParams{
			ID: "1",
		}

		expectedErr := errors.New("database unavailable")

		expected := &e.ServiceError{
			Code:       e.CodeInternalError,
			Message:    "Internal server error",
			StatusCode: http.StatusInternalServerError,
			Err:        expectedErr,
		}

		repo := &MockCarRepository{
			PartialUpdateFn: func(params c.CarUpdateParams) (c.Car, error) {
				return c.Car{}, expectedErr
			},
		}

		service := &DefaultCarService{
			repo: repo,
		}

		// Act
		_, err := service.PartialUpdate(params)

		// Assert
		assertServiceError(t, err, expected)
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
		expected := &e.ServiceError{
			Code:       e.CodeCarNotFound,
			Message:    "Car not found",
			StatusCode: http.StatusNotFound,
			Err:        e.ErrCarNotFound,
		}

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
		assertServiceError(t, err, expected)
	})

	t.Run("should return internal error when repository fails unexpectedly", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("database unavailable")

		expected := &e.ServiceError{
			Code:       e.CodeInternalError,
			Message:    "Internal server error",
			StatusCode: http.StatusInternalServerError,
			Err:        expectedErr,
		}

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
		assertServiceError(t, err, expected)
	})
}
