/*
I have these tests, but I am thinking that is not good idea to have the same implementation in my mock interfaces, that does not make sense I think, while I believe It would be good to add mock response or something in every test case, but that is only an idea, do you know a better approach?
*/

package controllers

import (
	"cars/models"
	e "cars/pkg/errors"
	"cars/services"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockCarRepository struct {
	FindFn   func(id string) (models.Car, error)
	ListFn   func() (models.Cars, error)
	CreateFn func(car *models.Car) error
	UpdateFn func(car *models.Car) error
}

func (m *MockCarRepository) Find(id string) (models.Car, error) {
	return m.FindFn(id)
}
func (m *MockCarRepository) List() (models.Cars, error) {
	return m.ListFn()
}
func (m *MockCarRepository) Create(car *models.Car) error {
	return m.CreateFn(car)
}
func (m *MockCarRepository) Update(car *models.Car) error {
	return m.UpdateFn(car)
}

func Test_Car_Get(t *testing.T) {
	tCases := []struct {
		name           string
		idParam        string
		findFn         func(id string) (models.Car, error)
		expectedStatus int
		expectedData   *models.Car
	}{
		{
			name:           "missing param",
			idParam:        "",
			findFn:         nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:    "repository error",
			idParam: "ABC123",
			findFn: func(id string) (models.Car, error) {
				return models.Car{}, errors.New("repository error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedData:   nil,
		},
		{
			name:    "car found",
			idParam: "ABC123",
			findFn: func(id string) (models.Car, error) {
				return models.Car{
					ID:       "ABC123",
					Make:     "Chevrolet",
					Model:    "Onix",
					Color:    "Black",
					Category: "Sedan",
					Year:     2025,
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedData: &models.Car{
				ID:       "ABC123",
				Make:     "Chevrolet",
				Model:    "Onix",
				Color:    "Black",
				Category: "Sedan",
				Year:     2025,
			},
		},
		{
			name:    "car not found",
			idParam: "DEF456",
			findFn: func(id string) (models.Car, error) {
				return models.Car{}, e.ErrCarNotFound
			},
			expectedStatus: http.StatusNotFound,
			expectedData:   nil,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewCarController(
				services.NewCarService(
					&MockCarRepository{
						FindFn: tc.findFn,
					},
				),
			)

			req, err := http.NewRequest("GET", "/cars/"+tc.idParam, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp := httptest.NewRecorder()
			http.HandlerFunc(controller.Get).ServeHTTP(resp, req)
			if status := resp.Code; status != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, status)
			}

			if tc.expectedData != nil {
				var car models.Car
				if err := json.NewDecoder(resp.Body).Decode(&car); err != nil {
					t.Fatal(err)
				}

				if car.ID != tc.expectedData.ID {
					t.Errorf("expected ID %v, got %v", tc.expectedData.ID, car.ID)
				}

				if car.Make != tc.expectedData.Make {
					t.Errorf("expected make %v, got %v", tc.expectedData.Make, car.Make)
				}

				if car.Model != tc.expectedData.Model {
					t.Errorf("expected model %v, got %v", tc.expectedData.Model, car.Model)
				}

				if car.Color != tc.expectedData.Color {
					t.Errorf("expected color %v, got %v", tc.expectedData.Color, car.Color)
				}

				if car.Category != tc.expectedData.Category {
					t.Errorf("expected category %v, got %v", tc.expectedData.Category, car.Category)
				}

				if car.Year != tc.expectedData.Year {
					t.Errorf("expected year %v, got %v", tc.expectedData.Year, car.Year)
				}
			}
		})
	}
}

func Test_Car_List(t *testing.T) {
	tCases := []struct {
		name           string
		listFn         func() (models.Cars, error)
		expectedStatus int
		expectedData   models.Cars
	}{
		{
			name: "list with 3 cars",
			listFn: func() (models.Cars, error) {
				return models.Cars{
					{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Package: "ABC", Color: "Black", Category: "Sedan", Year: 2025},
					{ID: "DEF456", Make: "Toyota", Model: "Yaris", Package: "DEF", Color: "Red", Category: "Sedan", Year: 2025},
					{ID: "GHI789", Make: "Renault", Model: "Arkana", Package: "GHI", Color: "White", Category: "Sedan", Year: 2025},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedData: models.Cars{
				{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Package: "ABC", Color: "Black", Category: "Sedan", Year: 2025},
				{ID: "DEF456", Make: "Toyota", Model: "Yaris", Package: "DEF", Color: "Red", Category: "Sedan", Year: 2025},
				{ID: "GHI789", Make: "Renault", Model: "Arkana", Package: "GHI", Color: "White", Category: "Sedan", Year: 2025},
			},
		},
		{
			name: "repository error",
			listFn: func() (models.Cars, error) {
				return nil, errors.New("repository error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedData:   nil,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewCarController(
				services.NewCarService(
					&MockCarRepository{
						ListFn: tc.listFn,
					},
				),
			)

			req, _ := http.NewRequest("GET", "/cars", nil)
			resp := httptest.NewRecorder()
			http.HandlerFunc(controller.List).ServeHTTP(resp, req)

			if resp.Code != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, resp.Code)
			}

			if tc.expectedData != nil {
				var cars models.Cars
				if err := json.NewDecoder(resp.Body).Decode(&cars); err != nil {
					t.Fatal(err)
				}

				if len(cars) != len(tc.expectedData) {
					t.Errorf("expected %d cars, got %d", len(tc.expectedData), len(cars))
				}

				actualCars := make(map[string]models.Car)
				for _, c := range cars {
					actualCars[c.ID] = c
				}

				for _, expected := range tc.expectedData {
					got, ok := actualCars[expected.ID]
					if !ok {
						t.Errorf("expected car ID %s not found", expected.ID)
						continue
					}

					if got.Make != expected.Make {
						t.Errorf("expected make %v, got %v", expected.Make, got.Make)
					}

					if got.Model != expected.Model {
						t.Errorf("expected model %v, got %v", expected.Model, got.Model)
					}

					if got.Color != expected.Color {
						t.Errorf("expected color %v, got %v", expected.Color, got.Color)
					}

					if got.Category != expected.Category {
						t.Errorf("expected category %v, got %v", expected.Category, got.Category)
					}

					if got.Year != expected.Year {
						t.Errorf("expected year %v, got %v", expected.Year, got.Year)
					}
				}
			}
		})
	}
}

func Test_Car_Create(t *testing.T) {
	tCases := []struct {
		name           string
		body           string
		createFn       func(car *models.Car) error
		expectedStatus int
		expectedData   *models.Car
	}{
		{
			name:           "invalid request body: missing comma",
			body:           `{"make": "Chevrolet" "model": "Onix"}`,
			createFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:           "year is missing",
			body:           `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan"}`,
			createFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:           "category is missing",
			body:           `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "year":2025}`,
			createFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:           "color is missing",
			body:           `{"make":"Chevrolet", "model":"Onix", "category":"Sedan", "year":2025}`,
			createFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:           "model is missing",
			body:           `{"make":"Chevrolet", "color":"Gray", "category":"Sedan", "year":2025}`,
			createFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:           "make is missing",
			body:           `{"model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			createFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name: "repository error",
			body: `{"make": "Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			createFn: func(car *models.Car) error {
				return errors.New("repository error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedData:   nil,
		},
		{
			name: "car created successfully",
			body: `{"make": "Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			createFn: func(car *models.Car) error {
				car.ID = "A1"
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedData:   &models.Car{ID: "A1", Make: "Chevrolet", Model: "Onix", Color: "Gray", Category: "Sedan", Year: 2025},
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewCarController(
				services.NewCarService(
					&MockCarRepository{
						CreateFn: tc.createFn,
					},
				),
			)

			req, err := http.NewRequest("POST", "/cars", strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			http.HandlerFunc(controller.Create).ServeHTTP(resp, req)
			if resp.Code != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, resp.Code)
			}

			if tc.expectedData != nil {
				var car models.Car
				if err := json.NewDecoder(resp.Body).Decode(&car); err != nil {
					t.Fatal(err)
				}

				if car.ID != tc.expectedData.ID {
					t.Errorf("expected ID %v, got %v", tc.expectedData.ID, car.ID)
				}

				if car.Make != tc.expectedData.Make {
					t.Errorf("expected make %v, got %v", tc.expectedData.Make, car.Make)
				}

				if car.Model != tc.expectedData.Model {
					t.Errorf("expected model %v, got %v", tc.expectedData.Model, car.Model)
				}

				if car.Color != tc.expectedData.Color {
					t.Errorf("expected color %v, got %v", tc.expectedData.Color, car.Color)
				}

				if car.Category != tc.expectedData.Category {
					t.Errorf("expected category %v, got %v", tc.expectedData.Category, car.Category)
				}

				if car.Year != tc.expectedData.Year {
					t.Errorf("expected year %v, got %v", tc.expectedData.Year, car.Year)
				}
			}
		})
	}
}

func Test_Car_Update(t *testing.T) {
	tCases := []struct {
		name           string
		idParam        string
		body           string
		updateFn       func(car *models.Car) error
		expectedStatus int
		expectedData   *models.Car
	}{
		{
			name:           "missing param",
			idParam:        "",
			body:           ``,
			updateFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:           "invalid request body: missing comma",
			idParam:        "ABC123",
			body:           `{"make": "Chevrolet" "model": "Onix"}`,
			updateFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:           "year is missing",
			idParam:        "ABC123",
			body:           `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan"}`,
			updateFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:           "category is missing",
			idParam:        "ABC123",
			body:           `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "year":2025}`,
			updateFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:           "color is missing",
			idParam:        "ABC123",
			body:           `{"make":"Chevrolet", "model":"Onix", "category":"Sedan", "year":2025}`,
			updateFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:           "model is missing",
			idParam:        "ABC123",
			body:           `{"make":"Chevrolet", "color":"Gray", "category":"Sedan", "year":2025}`,
			updateFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:           "make is missing",
			idParam:        "ABC123",
			body:           `{"model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			updateFn:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedData:   nil,
		},
		{
			name:    "car not found",
			idParam: "ABC123",
			body:    `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			updateFn: func(car *models.Car) error {
				return e.ErrCarNotFound
			},
			expectedStatus: http.StatusNotFound,
			expectedData:   nil,
		},
		{
			name:    "repository error",
			idParam: "ABC123",
			body:    `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			updateFn: func(car *models.Car) error {
				return errors.New("repository error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedData:   nil,
		},
		{
			name:    "car updated successfully",
			idParam: "ABC123",
			body:    `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			updateFn: func(car *models.Car) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedData:   &models.Car{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Color: "Gray", Category: "Sedan", Year: 2025},
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewCarController(
				services.NewCarService(
					&MockCarRepository{
						UpdateFn: tc.updateFn,
					},
				),
			)

			req, err := http.NewRequest("PUT", "/cars/"+tc.idParam, strings.NewReader(tc.body))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			http.HandlerFunc(controller.Update).ServeHTTP(resp, req)
			if resp.Code != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, resp.Code)
			}

			if tc.expectedData != nil {
				var car models.Car
				if err := json.NewDecoder(resp.Body).Decode(&car); err != nil {
					t.Fatal(err)
				}

				if car.ID != tc.expectedData.ID {
					t.Errorf("expected ID %v, got %v", tc.expectedData.ID, car.ID)
				}

				if car.Make != tc.expectedData.Make {
					t.Errorf("expected make %v, got %v", tc.expectedData.Make, car.Make)
				}

				if car.Model != tc.expectedData.Model {
					t.Errorf("expected model %v, got %v", tc.expectedData.Model, car.Model)
				}

				if car.Color != tc.expectedData.Color {
					t.Errorf("expected color %v, got %v", tc.expectedData.Color, car.Color)
				}

				if car.Category != tc.expectedData.Category {
					t.Errorf("expected category %v, got %v", tc.expectedData.Category, car.Category)
				}

				if car.Year != tc.expectedData.Year {
					t.Errorf("expected year %v, got %v", tc.expectedData.Year, car.Year)
				}
			}
		})
	}
}
