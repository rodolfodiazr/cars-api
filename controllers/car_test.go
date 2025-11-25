package controllers

import (
	"cars/models"
	e "cars/pkg/errors"
	"cars/pkg/httpx"
	"cars/services"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type MockCarRepository struct {
	FindFn   func(id string) (models.Car, error)
	ListFn   func(filters models.CarFilters) (models.Cars, error)
	CreateFn func(car *models.Car) error
	UpdateFn func(car *models.Car) error
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

func Test_Car_Get(t *testing.T) {
	tCases := []struct {
		name             string
		idParam          string
		findFn           func(id string) (models.Car, error)
		expectedStatus   int
		expectedResponse any
	}{
		{
			name:             "missing param",
			idParam:          "",
			findFn:           nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid path format: expected /cars/{id}"},
		},
		{
			name:    "repository error",
			idParam: "ABC123",
			findFn: func(id string) (models.Car, error) {
				return models.Car{}, errors.New("repository error")
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: httpx.ErrorResponse{Message: "Internal server error"},
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
			expectedResponse: httpx.SuccessResponse{
				Data: models.Car{
					ID:       "ABC123",
					Make:     "Chevrolet",
					Model:    "Onix",
					Color:    "Black",
					Category: "Sedan",
					Year:     2025,
				},
			},
		},
		{
			name:    "car not found",
			idParam: "DEF456",
			findFn: func(id string) (models.Car, error) {
				return models.Car{}, e.ErrCarNotFound
			},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: httpx.ErrorResponse{Message: "Car not found"},
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewCarController(
				services.NewCarService(
					&MockCarRepository{FindFn: tc.findFn},
				),
			)

			resp := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/cars/"+tc.idParam, nil)
			controller.Get(resp, req)

			// Check status code
			if resp.Code != tc.expectedStatus {
				t.Errorf("expected status %v, got %v", tc.expectedStatus, resp.Code)
				return
			}

			// Check error response
			if expected, ok := tc.expectedResponse.(httpx.ErrorResponse); ok {
				var got httpx.ErrorResponse
				if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
					t.Fatal(err)
					return
				}

				if got.Message != expected.Message {
					t.Errorf("expected error %v, got %v", expected.Message, got.Message)
				}
				return
			}

			// Check success response
			var got httpx.SuccessResponse
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatal(err)
				return
			}

			gotDataBytes, err := json.Marshal(got.Data)
			if err != nil {
				t.Fatal(err)
				return
			}

			gotCar := models.Car{}
			if err := json.Unmarshal(gotDataBytes, &gotCar); err != nil {
				t.Fatal(err)
				return
			}

			expectedCar, _ := tc.expectedResponse.(httpx.SuccessResponse).Data.(models.Car)
			if !reflect.DeepEqual(gotCar, expectedCar) {
				t.Errorf("expected car %+v, got %+v", expectedCar, gotCar)
			}
		})
	}
}

// func Test_Car_List(t *testing.T) {
// 	tCases := []struct {
// 		name             string
// 		listFn           func() (models.Cars, error)
// 		expectedStatus   int
// 		expectedResponse httpx.Response
// 	}{
// 		{
// 			name: "list with 3 cars",
// 			listFn: func() (models.Cars, error) {
// 				return models.Cars{
// 					{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Package: "ABC", Color: "Black", Category: "Sedan", Year: 2025},
// 					{ID: "DEF456", Make: "Toyota", Model: "Yaris", Package: "DEF", Color: "Red", Category: "Sedan", Year: 2025},
// 					{ID: "GHI789", Make: "Renault", Model: "Arkana", Package: "GHI", Color: "White", Category: "Sedan", Year: 2025},
// 				}, nil
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedResponse: httpx.Response{
// 				Data: models.Cars{
// 					{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Package: "ABC", Color: "Black", Category: "Sedan", Year: 2025},
// 					{ID: "DEF456", Make: "Toyota", Model: "Yaris", Package: "DEF", Color: "Red", Category: "Sedan", Year: 2025},
// 					{ID: "GHI789", Make: "Renault", Model: "Arkana", Package: "GHI", Color: "White", Category: "Sedan", Year: 2025},
// 				},
// 			},
// 		},
// 		{
// 			name: "repository error",
// 			listFn: func() (models.Cars, error) {
// 				return nil, errors.New("repository error")
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedResponse: httpx.Response{
// 				Error: "repository error",
// 			},
// 		},
// 	}

// 	for _, tc := range tCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			controller := NewCarController(
// 				services.NewCarService(
// 					&MockCarRepository{
// 						ListFn: tc.listFn,
// 					},
// 				),
// 			)

// 			req, _ := http.NewRequest("GET", "/cars", nil)
// 			resp := httptest.NewRecorder()
// 			http.HandlerFunc(controller.List).ServeHTTP(resp, req)

// 			if resp.Code != tc.expectedStatus {
// 				t.Errorf("expected status %v, got %v", tc.expectedStatus, resp.Code)
// 			}

// 			var respBody httpx.Response
// 			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
// 				t.Fatal(err)
// 			}

// 			if tc.expectedResponse.Error != "" {
// 				if respBody.Error != tc.expectedResponse.Error {
// 					t.Errorf("expected error %v, got %v", tc.expectedResponse.Error, respBody.Error)
// 				}
// 				return
// 			}

// 			dataBytes, err := json.Marshal(respBody.Data)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			var cars models.Cars
// 			if err := json.Unmarshal(dataBytes, &cars); err != nil {
// 				t.Fatal(err)
// 			}

// 			expectedCars := tc.expectedResponse.Data.(models.Cars)
// 			if len(cars) != len(expectedCars) {
// 				t.Errorf("expected %d cars, got %d", len(expectedCars), len(cars))
// 			}

// 			actualCars := make(map[string]models.Car)
// 			for _, c := range cars {
// 				actualCars[c.ID] = c
// 			}

// 			for _, expected := range expectedCars {
// 				got, ok := actualCars[expected.ID]
// 				if !ok {
// 					t.Errorf("expected car ID %s not found", expected.ID)
// 					continue
// 				}

// 				if got.Make != expected.Make {
// 					t.Errorf("expected make %v, got %v", expected.Make, got.Make)
// 				}

// 				if got.Model != expected.Model {
// 					t.Errorf("expected model %v, got %v", expected.Model, got.Model)
// 				}

// 				if got.Color != expected.Color {
// 					t.Errorf("expected color %v, got %v", expected.Color, got.Color)
// 				}

// 				if got.Category != expected.Category {
// 					t.Errorf("expected category %v, got %v", expected.Category, got.Category)
// 				}

// 				if got.Year != expected.Year {
// 					t.Errorf("expected year %v, got %v", expected.Year, got.Year)
// 				}
// 			}
// 		})
// 	}
// }

// func Test_Car_Create(t *testing.T) {
// 	tCases := []struct {
// 		name             string
// 		body             string
// 		createFn         func(car *models.Car) error
// 		expectedStatus   int
// 		expectedResponse httpx.Response
// 	}{
// 		{
// 			name:           "invalid request body: missing comma",
// 			body:           `{"make": "Chevrolet" "model": "Onix"}`,
// 			createFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "invalid request body",
// 			},
// 		},
// 		{
// 			name:           "year is missing",
// 			body:           `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan"}`,
// 			createFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "year is not valid",
// 			},
// 		},
// 		{
// 			name:           "category is missing",
// 			body:           `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "year":2025}`,
// 			createFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "category is required",
// 			},
// 		},
// 		{
// 			name:           "color is missing",
// 			body:           `{"make":"Chevrolet", "model":"Onix", "category":"Sedan", "year":2025}`,
// 			createFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "color is required",
// 			},
// 		},
// 		{
// 			name:           "model is missing",
// 			body:           `{"make":"Chevrolet", "color":"Gray", "category":"Sedan", "year":2025}`,
// 			createFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "model is required",
// 			},
// 		},
// 		{
// 			name:           "make is missing",
// 			body:           `{"model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
// 			createFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "make is required",
// 			},
// 		},
// 		{
// 			name: "repository error",
// 			body: `{"make": "Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
// 			createFn: func(car *models.Car) error {
// 				return errors.New("repository error")
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedResponse: httpx.Response{
// 				Error: "repository error",
// 			},
// 		},
// 		{
// 			name: "car created successfully",
// 			body: `{"make": "Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
// 			createFn: func(car *models.Car) error {
// 				car.ID = "A1"
// 				return nil
// 			},
// 			expectedStatus: http.StatusCreated,
// 			expectedResponse: httpx.Response{
// 				Data: models.Car{ID: "A1", Make: "Chevrolet", Model: "Onix", Color: "Gray", Category: "Sedan", Year: 2025},
// 			},
// 		},
// 	}

// 	for _, tc := range tCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			controller := NewCarController(
// 				services.NewCarService(
// 					&MockCarRepository{
// 						CreateFn: tc.createFn,
// 					},
// 				),
// 			)

// 			req, err := http.NewRequest("POST", "/cars", strings.NewReader(tc.body))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			resp := httptest.NewRecorder()
// 			http.HandlerFunc(controller.Create).ServeHTTP(resp, req)
// 			if resp.Code != tc.expectedStatus {
// 				t.Errorf("expected status %v, got %v", tc.expectedStatus, resp.Code)
// 			}

// 			var respBody httpx.Response
// 			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
// 				t.Fatal(err)
// 			}

// 			if tc.expectedResponse.Error != "" {
// 				if respBody.Error != tc.expectedResponse.Error {
// 					t.Errorf("expected error %v, got %v", tc.expectedResponse.Error, respBody.Error)
// 				}
// 				return
// 			}

// 			dataBytes, err := json.Marshal(respBody.Data)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			var car models.Car
// 			if err := json.Unmarshal(dataBytes, &car); err != nil {
// 				t.Fatal(err)
// 			}

// 			expectedCar := tc.expectedResponse.Data.(models.Car)
// 			if car.ID != expectedCar.ID {
// 				t.Errorf("expected ID %v, got %v", expectedCar.ID, car.ID)
// 			}

// 			if car.Make != expectedCar.Make {
// 				t.Errorf("expected make %v, got %v", expectedCar.Make, car.Make)
// 			}

// 			if car.Model != expectedCar.Model {
// 				t.Errorf("expected model %v, got %v", expectedCar.Model, car.Model)
// 			}

// 			if car.Color != expectedCar.Color {
// 				t.Errorf("expected color %v, got %v", expectedCar.Color, car.Color)
// 			}

// 			if car.Category != expectedCar.Category {
// 				t.Errorf("expected category %v, got %v", expectedCar.Category, car.Category)
// 			}

// 			if car.Year != expectedCar.Year {
// 				t.Errorf("expected year %v, got %v", expectedCar.Year, car.Year)
// 			}
// 		})
// 	}
// }

// func Test_Car_Update(t *testing.T) {
// 	tCases := []struct {
// 		name             string
// 		idParam          string
// 		body             string
// 		updateFn         func(car *models.Car) error
// 		expectedStatus   int
// 		expectedResponse httpx.Response
// 	}{
// 		{
// 			name:           "missing param",
// 			idParam:        "",
// 			body:           ``,
// 			updateFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "id is required",
// 			},
// 		},
// 		{
// 			name:           "invalid request body: missing comma",
// 			idParam:        "ABC123",
// 			body:           `{"make": "Chevrolet" "model": "Onix"}`,
// 			updateFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "invalid request body",
// 			},
// 		},
// 		{
// 			name:           "year is missing",
// 			idParam:        "ABC123",
// 			body:           `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan"}`,
// 			updateFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "year is not valid",
// 			},
// 		},
// 		{
// 			name:           "category is missing",
// 			idParam:        "ABC123",
// 			body:           `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "year":2025}`,
// 			updateFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "category is required",
// 			},
// 		},
// 		{
// 			name:           "color is missing",
// 			idParam:        "ABC123",
// 			body:           `{"make":"Chevrolet", "model":"Onix", "category":"Sedan", "year":2025}`,
// 			updateFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "color is required",
// 			},
// 		},
// 		{
// 			name:           "model is missing",
// 			idParam:        "ABC123",
// 			body:           `{"make":"Chevrolet", "color":"Gray", "category":"Sedan", "year":2025}`,
// 			updateFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "model is required",
// 			},
// 		},
// 		{
// 			name:           "make is missing",
// 			idParam:        "ABC123",
// 			body:           `{"model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
// 			updateFn:       nil,
// 			expectedStatus: http.StatusBadRequest,
// 			expectedResponse: httpx.Response{
// 				Error: "make is required",
// 			},
// 		},
// 		{
// 			name:    "car not found",
// 			idParam: "ABC123",
// 			body:    `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
// 			updateFn: func(car *models.Car) error {
// 				return e.ErrCarNotFound
// 			},
// 			expectedStatus: http.StatusNotFound,
// 			expectedResponse: httpx.Response{
// 				Error: "car not found",
// 			},
// 		},
// 		{
// 			name:    "repository error",
// 			idParam: "ABC123",
// 			body:    `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
// 			updateFn: func(car *models.Car) error {
// 				return errors.New("repository error")
// 			},
// 			expectedStatus: http.StatusInternalServerError,
// 			expectedResponse: httpx.Response{
// 				Error: "repository error",
// 			},
// 		},
// 		{
// 			name:    "car updated successfully",
// 			idParam: "ABC123",
// 			body:    `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
// 			updateFn: func(car *models.Car) error {
// 				return nil
// 			},
// 			expectedStatus: http.StatusOK,
// 			expectedResponse: httpx.Response{
// 				Data: models.Car{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Color: "Gray", Category: "Sedan", Year: 2025},
// 			},
// 		},
// 	}

// 	for _, tc := range tCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			controller := NewCarController(
// 				services.NewCarService(
// 					&MockCarRepository{
// 						UpdateFn: tc.updateFn,
// 					},
// 				),
// 			)

// 			req, err := http.NewRequest("PUT", "/cars/"+tc.idParam, strings.NewReader(tc.body))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			req.Header.Set("Content-Type", "application/json")

// 			resp := httptest.NewRecorder()
// 			http.HandlerFunc(controller.Update).ServeHTTP(resp, req)
// 			if resp.Code != tc.expectedStatus {
// 				t.Errorf("expected status %v, got %v", tc.expectedStatus, resp.Code)
// 			}

// 			var respBody httpx.Response
// 			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
// 				t.Fatal(err)
// 			}

// 			if tc.expectedResponse.Error != "" {
// 				if respBody.Error != tc.expectedResponse.Error {
// 					t.Errorf("expected error %v, got %v", tc.expectedResponse.Error, respBody.Error)
// 				}
// 				return
// 			}

// 			dataBytes, err := json.Marshal(respBody.Data)
// 			if err != nil {
// 				t.Fatal(err)
// 			}

// 			var car models.Car
// 			if err := json.Unmarshal(dataBytes, &car); err != nil {
// 				t.Fatal(err)
// 			}

// 			expectedCar := tc.expectedResponse.Data.(models.Car)
// 			if car.ID != expectedCar.ID {
// 				t.Errorf("expected ID %v, got %v", expectedCar.ID, car.ID)
// 			}

// 			if car.Make != expectedCar.Make {
// 				t.Errorf("expected make %v, got %v", expectedCar.Make, car.Make)
// 			}

// 			if car.Model != expectedCar.Model {
// 				t.Errorf("expected model %v, got %v", expectedCar.Model, car.Model)
// 			}

// 			if car.Color != expectedCar.Color {
// 				t.Errorf("expected color %v, got %v", expectedCar.Color, car.Color)
// 			}

// 			if car.Category != expectedCar.Category {
// 				t.Errorf("expected category %v, got %v", expectedCar.Category, car.Category)
// 			}

// 			if car.Year != expectedCar.Year {
// 				t.Errorf("expected year %v, got %v", expectedCar.Year, car.Year)
// 			}
// 		})
// 	}
// }
