package models

import (
	"errors"
	"testing"
)

func Test_Car_Validate(t *testing.T) {
	tCases := []struct {
		name          string
		car           Car
		expectedError error
	}{
		{
			name: "make is not defined",
			car: Car{
				ID:       "C001",
				Model:    "Onix",
				Category: "Sedan",
				Color:    "Red",
				Year:     2025,
			},
			expectedError: errors.New("make is required"),
		},
		{
			name: "model is not defined",
			car: Car{
				ID:       "C002",
				Make:     "Ford",
				Category: "Sedan",
				Color:    "Black",
				Year:     2025,
			},
			expectedError: errors.New("model is required"),
		},
		{
			name: "category is not defined",
			car: Car{
				ID:    "C003",
				Make:  "Kia",
				Model: "Rio",
				Color: "White",
				Year:  2025,
			},
			expectedError: errors.New("category is required"),
		},
		{
			name: "color is not defined",
			car: Car{
				ID:       "C004",
				Make:     "Nissan",
				Model:    "Versa",
				Category: "Sedan",
				Year:     2025,
			},
			expectedError: errors.New("color is required"),
		},
		{
			name: "year is not defined",
			car: Car{
				ID:       "C005",
				Make:     "Renault",
				Model:    "Arkana",
				Category: "Sedan",
				Color:    "Blue",
			},
			expectedError: errors.New("year is not valid"),
		},
		{
			name: "All necessary properties are defined",
			car: Car{
				ID:       "C006",
				Make:     "Suzuki",
				Model:    "Swift",
				Category: "Sedan",
				Color:    "Gray",
				Year:     2025,
			},
			expectedError: nil,
		},
	}

	for _, tCase := range tCases {
		t.Run(tCase.name, func(t *testing.T) {
			got := tCase.car.Validate()
			expected := tCase.expectedError

			if (got == nil && expected != nil) ||
				(got != nil && expected == nil) ||
				(got != nil && expected != nil && got.Error() != expected.Error()) {
				t.Errorf("expected error %v, got %v", expected, got)
			}
		})
	}
}
