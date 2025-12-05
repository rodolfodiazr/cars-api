package repositories

import (
	"cars/models"
	"testing"
)

func TestDefaultCarRepository_List(t *testing.T) {
	// Initial seed data
	cars := map[string]models.Car{
		"JHK290XJ": {
			ID:   "JHK290XJ",
			Make: "Ford", Model: "F10", Package: "Base",
			Color: "Silver", Year: 2010, Category: "Truck",
			Mileage: 120123, Price: 1999900,
		},
		"FWL37LA": {
			ID:   "FWL37LA",
			Make: "Toyota", Model: "Camry", Package: "SE",
			Color: "White", Year: 2019, Category: "Sedan",
			Mileage: 3999, Price: 2899000,
		},
		"1I3XJRLLC": {
			ID:   "1I3XJRLLC",
			Make: "Toyota", Model: "Rav4", Package: "XSE",
			Color: "Red", Year: 2018, Category: "SUV",
			Mileage: 24001, Price: 2275000,
		},
		"DKU43920S": {
			ID:   "DKU43920S",
			Make: "Ford", Model: "Bronco", Package: "Badlands",
			Color: "Burnt Orange", Year: 2022, Category: "SUV",
			Mileage: 1, Price: 4499000,
		},
	}

	repo := NewCarRepository(cars)

	tests := []struct {
		name     string
		filters  models.CarFilters
		expected []string // list of expected IDs
	}{
		{
			name:     "no filters returns all",
			filters:  models.CarFilters{},
			expected: []string{"1", "2", "3"},
		},
		{
			name:    "filter by make",
			filters: models.CarFilters{Make: "Toyota"},
			expected: []string{
				"1", "3",
			},
		},
		{
			name:    "filter by model",
			filters: models.CarFilters{Model: "Civic"},
			expected: []string{
				"2",
			},
		},
		{
			name:    "filter by year",
			filters: models.CarFilters{Year: 2021},
			expected: []string{
				"2", "3",
			},
		},
		{
			name: "filter by make and year",
			filters: models.CarFilters{
				Make: "Toyota",
				Year: 2021,
			},
			expected: []string{"3"},
		},
		{
			name: "case-insensitive filtering",
			filters: models.CarFilters{
				Make:  "tOyOtA",
				Model: "cAmRy",
			},
			expected: []string{"3"},
		},
		{
			name:     "no match returns empty list",
			filters:  models.CarFilters{Make: "Ford"},
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
