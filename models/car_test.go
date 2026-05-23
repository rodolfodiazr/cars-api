package models

import (
	"errors"
	"testing"
	"time"
)

func TestCar_ValidateForCreate(t *testing.T) {
	validMileage := int64(10000)
	validPrice := int64(25000)
	currentYear := time.Now().Year()

	validCar := Car{
		Make:     "Toyota",
		Model:    "Corolla",
		Color:    "Black",
		Category: "Sedan",
		Year:     currentYear,
		Mileage:  &validMileage,
		Price:    &validPrice,
	}

	tests := []struct {
		name    string
		car     Car
		wantErr error
	}{
		{
			name:    "should succeed for valid car",
			car:     validCar,
			wantErr: nil,
		},
		{
			name: "should fail when ID is provided",
			car: func() Car {
				c := validCar
				c.ID = "1"
				return c
			}(),
			wantErr: ErrCarIDMustBeEmpty,
		},
		{
			name: "should fail when make is blank",
			car: func() Car {
				c := validCar
				c.Make = ""
				return c
			}(),
			wantErr: ErrCarMakeRequired,
		},
		{
			name: "should fail when model is blank",
			car: func() Car {
				c := validCar
				c.Model = ""
				return c
			}(),
			wantErr: ErrCarModelRequired,
		},
		{
			name: "should fail when color is blank",
			car: func() Car {
				c := validCar
				c.Color = ""
				return c
			}(),
			wantErr: ErrCarColorRequired,
		},
		{
			name: "should fail when category is blank",
			car: func() Car {
				c := validCar
				c.Category = ""
				return c
			}(),
			wantErr: ErrCarCategoryRequired,
		},
		{
			name: "should fail when year is zero",
			car: func() Car {
				c := validCar
				c.Year = 0
				return c
			}(),
			wantErr: ErrInvalidYear,
		},
		{
			name: "should fail when year is in the future",
			car: func() Car {
				c := validCar
				c.Year = currentYear + 1
				return c
			}(),
			wantErr: ErrInvalidYear,
		},
		{
			name: "should fail when mileage is negative",
			car: func() Car {
				c := validCar
				mileage := int64(-1)
				c.Mileage = &mileage
				return c
			}(),
			wantErr: ErrInvalidMileage,
		},
		{
			name: "should fail when price is negative",
			car: func() Car {
				c := validCar
				price := int64(-1.0)
				c.Price = &price
				return c
			}(),
			wantErr: ErrInvalidPrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := tt.car.ValidateForCreate()

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func Test_ValidateForUpdate(t *testing.T) {
	validMileage := int64(10000)
	validPrice := int64(25000)
	currentYear := time.Now().Year()

	validCar := Car{
		ID:       "1",
		Make:     "Toyota",
		Model:    "Corolla",
		Color:    "Black",
		Category: "Sedan",
		Year:     currentYear,
		Mileage:  &validMileage,
		Price:    &validPrice,
	}

	tests := []struct {
		name    string
		car     Car
		wantErr error
	}{
		{
			name:    "should succeed for valid car",
			car:     validCar,
			wantErr: nil,
		},
		{
			name:    "should fail when ID is missing",
			car:     Car{},
			wantErr: ErrCarIDRequiredForUpdate,
		},
		{
			name: "should fail when make is blank",
			car: func() Car {
				c := validCar
				c.Make = ""
				return c
			}(),
			wantErr: ErrCarMakeRequired,
		},
		{
			name: "should fail when model is blank",
			car: func() Car {
				c := validCar
				c.Model = ""
				return c
			}(),
			wantErr: ErrCarModelRequired,
		},
		{
			name: "should fail when color is blank",
			car: func() Car {
				c := validCar
				c.Color = ""
				return c
			}(),
			wantErr: ErrCarColorRequired,
		},
		{
			name: "should fail when category is blank",
			car: func() Car {
				c := validCar
				c.Category = ""
				return c
			}(),
			wantErr: ErrCarCategoryRequired,
		},
		{
			name: "should fail when year is zero",
			car: func() Car {
				c := validCar
				c.Year = 0
				return c
			}(),
			wantErr: ErrInvalidYear,
		},
		{
			name: "should fail when year is in the future",
			car: func() Car {
				c := validCar
				c.Year = currentYear + 1
				return c
			}(),
			wantErr: ErrInvalidYear,
		},
		{
			name: "should fail when mileage is negative",
			car: func() Car {
				c := validCar
				mileage := int64(-1)
				c.Mileage = &mileage
				return c
			}(),
			wantErr: ErrInvalidMileage,
		},
		{
			name: "should fail when price is negative",
			car: func() Car {
				c := validCar
				price := int64(-1)
				c.Price = &price
				return c
			}(),
			wantErr: ErrInvalidPrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := tt.car.ValidateForUpdate()

			// Assert
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func Test_Car_Validate(t *testing.T) {
	currentYear := time.Now().Year()

	tCases := []struct {
		name          string
		car           Car
		expectedError error
	}{
		{
			name: "make is not defined",
			car: Car{
				ID:       "C001",
				Model:    "Chevrolet",
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
			name: "year is not defined",
			car: Car{
				ID:       "C006",
				Make:     "Toyota",
				Model:    "Corolla",
				Category: "Sedan",
				Color:    "Silver",
				Year:     currentYear + 10,
			},
			expectedError: ErrInvalidYear,
		},
		{
			name: "All necessary properties are defined",
			car: Car{
				ID:       "C007",
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
			got := tCase.car.validate()
			expected := tCase.expectedError

			if (got == nil && expected != nil) ||
				(got != nil && expected == nil) ||
				(got != nil && expected != nil && got.Error() != expected.Error()) {
				t.Errorf("expected error %v, got %v", expected, got)
			}
		})
	}
}
