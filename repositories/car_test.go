package repositories

import (
	"cars/models"
	e "cars/pkg/errors"
	u "cars/pkg/utils"
	"errors"
	"testing"
)

func TestDefaultCarRepository_Find(t *testing.T) {
	t.Run("should return car when it exists", func(t *testing.T) {
		// Arrange
		expected := models.Car{
			ID:    "1",
			Make:  "Toyota",
			Model: "Corolla",
		}

		repo := &DefaultCarRepository{
			cars: map[string]models.Car{
				expected.ID: expected,
			},
		}

		// Act
		got, err := repo.Find(expected.ID)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if expected != got {
			t.Errorf("expected %v car, got %v", expected, got)
		}
	})

	t.Run("should return error when car does not exist", func(t *testing.T) {
		// Arrange
		repo := &DefaultCarRepository{
			cars: map[string]models.Car{},
		}

		// Act
		_, err := repo.Find("missing-id")

		// Assert
		if err == nil {
			t.Fatalf("expected error, got nil")
		}

		if !errors.Is(err, e.ErrCarNotFound) {
			t.Fatalf("expected %v, got %v", e.ErrCarNotFound, err)
		}
	})
}

func TestDefaultCarRepository_List(t *testing.T) {
	// Initial seed data
	cars := map[string]models.Car{
		"1": {
			ID:   "1",
			Make: "Ford", Model: "F10", Package: u.Ptr("Base"),
			Color: "Silver", Year: 2010, Category: "Truck",
			Mileage: u.Ptr(int64(120123)), Price: u.Ptr(int64(1999900)),
		},
		"2": {
			ID:   "2",
			Make: "Toyota", Model: "Camry", Package: u.Ptr("SE"),
			Color: "White", Year: 2019, Category: "Sedan",
			Mileage: u.Ptr(int64(3999)), Price: u.Ptr(int64(2899000)),
		},
		"3": {
			ID:   "3",
			Make: "Toyota", Model: "Rav4", Package: u.Ptr("XSE"),
			Color: "Red", Year: 2018, Category: "SUV",
			Mileage: u.Ptr(int64(24001)), Price: u.Ptr(int64(2275000)),
		},
		"4": {
			ID:   "4",
			Make: "Ford", Model: "Bronco", Package: u.Ptr("Badlands"),
			Color: "Burnt Orange", Year: 2022, Category: "SUV",
			Mileage: u.Ptr(int64(0)), Price: u.Ptr(int64(4499000)),
		},
	}

	repo := NewCarRepository(cars)

	tests := []struct {
		name     string
		filters  models.CarFilters
		expected []string
	}{
		{
			name:     "No filters, returns all",
			filters:  models.CarFilters{},
			expected: []string{"1", "2", "3", "4"},
		},
		{
			name:     "Filter by make",
			filters:  models.CarFilters{Make: "Toyota"},
			expected: []string{"2", "3"},
		},
		{
			name:     "Filter by model",
			filters:  models.CarFilters{Model: "Bronco"},
			expected: []string{"4"},
		},
		{
			name:     "Filter by year",
			filters:  models.CarFilters{Year: u.Ptr(2010)},
			expected: []string{"1"},
		},
		{
			name: "Filter by make and model",
			filters: models.CarFilters{
				Make:  "Ford",
				Model: "F10",
			},
			expected: []string{"1"},
		},
		{
			name: "Filter by make and year",
			filters: models.CarFilters{
				Make: "Toyota",
				Year: u.Ptr(2019),
			},
			expected: []string{"2"},
		},
		{
			name: "Filter by model and year",
			filters: models.CarFilters{
				Model: "Rav4",
				Year:  u.Ptr(2018),
			},
			expected: []string{"3"},
		},
		{
			name: "Filter by make, model and year",
			filters: models.CarFilters{
				Make:  "Ford",
				Model: "Bronco",
				Year:  u.Ptr(2022),
			},
			expected: []string{"4"},
		},
		{
			name: "Case-insensitive filtering",
			filters: models.CarFilters{
				Make:  "tOyOtA",
				Model: "cAmRy",
			},
			expected: []string{"2"},
		},
		{
			name:     "No match, returns empty list",
			filters:  models.CarFilters{Make: "BMW"},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.List(tt.filters)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d cars, got %d", len(tt.expected), len(got))
			}

			// Verify IDs
			for _, id := range tt.expected {
				found := false
				for _, car := range got {
					if car.ID == id {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected car ID %s not found in result", id)
				}
			}
		})
	}
}

func TestDefaultCarRepository_Create(t *testing.T) {
	// Arrange
	repo := &DefaultCarRepository{
		cars: map[string]models.Car{},
	}

	car := &models.Car{
		Make:  "Honda",
		Model: "Civic",
	}

	// Act
	err := repo.Create(car)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if car.ID == "" {
		t.Fatal("expected generated car ID")
	}

	stored, exists := repo.cars[car.ID]

	if !exists {
		t.Fatal("expected car to be stored in repository")
	}

	if stored != *car {
		t.Fatalf("expected stored car %+v, got %+v", *car, stored)
	}
}

func TestDefaultCarRepository_Update(t *testing.T) {
	t.Run("should update existing car", func(t *testing.T) {
		// Arrange
		existing := models.Car{
			ID:    "1",
			Make:  "Toyota",
			Model: "Corolla",
		}

		repo := &DefaultCarRepository{
			cars: map[string]models.Car{
				existing.ID: existing,
			},
		}

		updated := &models.Car{
			ID:    "1",
			Make:  "Honda",
			Model: "Civic",
		}

		// Act
		err := repo.Update(updated)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		stored, exists := repo.cars[updated.ID]

		if !exists {
			t.Fatal("expected updated car to exist in repository")
		}

		if stored != *updated {
			t.Fatalf("expected stored car %+v, got %+v", *updated, stored)
		}
	})

	t.Run("should return error when car does not exist", func(t *testing.T) {
		// Arrange
		repo := &DefaultCarRepository{
			cars: map[string]models.Car{},
		}

		car := &models.Car{
			ID:    "missing-id",
			Make:  "Mazda",
			Model: "3",
		}

		// Act
		err := repo.Update(car)

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}

		if !errors.Is(err, e.ErrCarNotFound) {
			t.Fatalf("expected ErrCarNotFound, got %v", err)
		}
	})
}

func TestDefaultCarRepository_Delete(t *testing.T) {
	t.Run("should delete existing car", func(t *testing.T) {
		// Arrange
		car := models.Car{
			ID:    "1",
			Make:  "Toyota",
			Model: "Corolla",
		}

		repo := &DefaultCarRepository{
			cars: map[string]models.Car{
				car.ID: car,
			},
		}

		// Act
		err := repo.Delete(car.ID)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if _, exists := repo.cars[car.ID]; exists {
			t.Fatal("expected car to be deleted from repository")
		}

		if len(repo.cars) != 0 {
			t.Fatalf("expected repository to be empty, got %d items", len(repo.cars))
		}
	})

	t.Run("should return error when car does not exist", func(t *testing.T) {
		// Arrange
		repo := &DefaultCarRepository{
			cars: map[string]models.Car{},
		}

		// Act
		err := repo.Delete("missing-id")

		// Assert
		if err == nil {
			t.Fatal("expected error but got nil")
		}

		if !errors.Is(err, e.ErrCarNotFound) {
			t.Fatalf("expected ErrCarNotFound, got %v", err)
		}
	})
}
