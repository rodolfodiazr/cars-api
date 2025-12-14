package controllers

import (
	"cars/models"
	e "cars/pkg/errors"
	"cars/pkg/httpx"
	"cars/services"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
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
			name:    "car not found",
			idParam: "DEF456",
			findFn: func(id string) (models.Car, error) {
				return models.Car{}, e.ErrCarNotFound
			},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: httpx.ErrorResponse{Message: "Car not found"},
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
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewCarController(
				services.NewCarService(
					&MockCarRepository{FindFn: tc.findFn},
				),
			)

			router := chi.NewRouter()
			router.Route("/cars", func(r chi.Router) {
				r.Get("/{id:[A-Za-z0-9-]+}", controller.Get)
			})

			resp := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/cars/"+tc.idParam, nil)

			router.ServeHTTP(resp, req)

			// Check status code
			if resp.Code != tc.expectedStatus {
				t.Fatalf("expected status %v, got %v", tc.expectedStatus, resp.Code)
			}

			body := resp.Body.Bytes()

			switch expected := tc.expectedResponse.(type) {
			case httpx.ErrorResponse:
				var got httpx.ErrorResponse
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}

				if got.Message != expected.Message {
					t.Fatalf("expected error %v, got %v", expected.Message, got.Message)
				}
			case httpx.SuccessResponse:
				var got httpx.SuccessResponse
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}

				gotDataBytes, err := json.Marshal(got.Data)
				if err != nil {
					t.Fatal(err)
				}

				gotCar := models.Car{}
				if err := json.Unmarshal(gotDataBytes, &gotCar); err != nil {
					t.Fatal(err)
				}

				expectedCar := tc.expectedResponse.(httpx.SuccessResponse).Data.(models.Car)
				if !reflect.DeepEqual(gotCar, expectedCar) {
					t.Fatalf("expected car %+v, got %+v", expectedCar, gotCar)
				}
			}
		})
	}
}

func Test_Car_List(t *testing.T) {
	tCases := []struct {
		name             string
		queryParams      string
		listFn           func(f models.CarFilters) (models.Cars, error)
		expectedStatus   int
		expectedResponse any
	}{
		{
			name: "list with 3 cars",
			listFn: func(f models.CarFilters) (models.Cars, error) {
				return models.Cars{
					{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Package: "ABC", Color: "Black", Category: "Sedan", Year: 2025},
					{ID: "DEF456", Make: "Toyota", Model: "Yaris", Package: "DEF", Color: "Red", Category: "Sedan", Year: 2025},
					{ID: "GHI789", Make: "Renault", Model: "Arkana", Package: "GHI", Color: "White", Category: "Sedan", Year: 2025},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: httpx.SuccessResponse{
				Data: models.Cars{
					{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Package: "ABC", Color: "Black", Category: "Sedan", Year: 2025},
					{ID: "DEF456", Make: "Toyota", Model: "Yaris", Package: "DEF", Color: "Red", Category: "Sedan", Year: 2025},
					{ID: "GHI789", Make: "Renault", Model: "Arkana", Package: "GHI", Color: "White", Category: "Sedan", Year: 2025},
				},
			},
		},
		{
			name:             "invalid params",
			queryParams:      "?year=MMXXV",
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: httpx.ErrorResponse{Message: "Internal server error"},
		},
		{
			name: "repository error",
			listFn: func(f models.CarFilters) (models.Cars, error) {
				return nil, errors.New("repository error")
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: httpx.ErrorResponse{Message: "Internal server error"},
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewCarController(
				services.NewCarService(
					&MockCarRepository{ListFn: tc.listFn},
				),
			)

			router := chi.NewRouter()
			router.Route("/cars", func(r chi.Router) {
				r.Get("/", controller.List)
			})

			resp := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("%s%s", "/cars", tc.queryParams), nil)

			router.ServeHTTP(resp, req)

			// Check status code
			if resp.Code != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d", tc.expectedStatus, resp.Code)
			}

			body := resp.Body.Bytes()

			switch expected := tc.expectedResponse.(type) {
			case httpx.ErrorResponse:
				var got httpx.ErrorResponse
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}

				if got.Message != expected.Message {
					t.Fatalf("expected error %v, got %v", expected.Message, got.Message)
				}
			case httpx.SuccessResponse:
				var got httpx.SuccessResponse
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}

				gotDataBytes, err := json.Marshal(got.Data)
				if err != nil {
					t.Fatal(err)
				}

				var gotCars models.Cars
				if err := json.Unmarshal(gotDataBytes, &gotCars); err != nil {
					t.Fatal(err)
				}

				expectedCars := tc.expectedResponse.(httpx.SuccessResponse).Data.(models.Cars)
				if len(gotCars) != len(expectedCars) {
					t.Fatalf("expected %d cars, got %d", len(expectedCars), len(gotCars))
				}

				if !reflect.DeepEqual(gotCars, expectedCars) {
					t.Fatalf("expected cars %+v, got %+v", expectedCars, gotCars)
				}
			}
		})
	}
}

func Test_Car_Create(t *testing.T) {
	tCases := []struct {
		name             string
		body             string
		createFn         func(car *models.Car) error
		expectedStatus   int
		expectedResponse any
	}{
		{
			name:             "invalid request body: missing comma",
			body:             `{"make": "Chevrolet" "model": "Onix"}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "year is missing",
			body:             `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan"}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "category is missing",
			body:             `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "color is missing",
			body:             `{"make":"Chevrolet", "model":"Onix", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "model is missing",
			body:             `{"make":"Chevrolet", "color":"Gray", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "make is missing",
			body:             `{"model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "id not allowed on create",
			body:             `{"id": "ABC123", "make": "Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "ID must not be provided when creating a new record"},
		},
		{
			name: "repository error",
			body: `{"make": "Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			createFn: func(car *models.Car) error {
				return errors.New("repository error")
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: httpx.ErrorResponse{Message: "Internal server error"},
		},
		{
			name: "car created successfully",
			body: `{"make": "Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			createFn: func(car *models.Car) error {
				car.ID = "A1"
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedResponse: httpx.SuccessResponse{
				Data: models.Car{
					ID:       "A1",
					Make:     "Chevrolet",
					Model:    "Onix",
					Color:    "Gray",
					Category: "Sedan",
					Year:     2025,
				},
			},
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewCarController(
				services.NewCarService(
					&MockCarRepository{CreateFn: tc.createFn},
				),
			)

			router := chi.NewRouter()
			router.Route("/cars", func(r chi.Router) {
				r.Post("/", controller.Create)
			})

			resp := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/cars", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(resp, req)

			// Check status code
			if resp.Code != tc.expectedStatus {
				t.Fatalf("expected status %v, got %v", tc.expectedStatus, resp.Code)
			}

			body := resp.Body.Bytes()

			switch expected := tc.expectedResponse.(type) {
			case httpx.ErrorResponse:
				var got httpx.ErrorResponse
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}

				if got.Message != expected.Message {
					t.Fatalf("expected error %v, got %v", expected.Message, got.Message)
				}
			case httpx.SuccessResponse:
				var got httpx.SuccessResponse
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}

				gotDataBytes, err := json.Marshal(got.Data)
				if err != nil {
					t.Fatal(err)
				}

				var gotCar models.Car
				if err := json.Unmarshal(gotDataBytes, &gotCar); err != nil {
					t.Fatal(err)
				}

				expectedCar := tc.expectedResponse.(httpx.SuccessResponse).Data.(models.Car)
				if !reflect.DeepEqual(gotCar, expectedCar) {
					t.Fatalf("expected car %+v, got %+v", expectedCar, gotCar)
				}
			}
		})
	}
}

func Test_Car_Update(t *testing.T) {
	tCases := []struct {
		name             string
		idParam          string
		body             string
		updateFn         func(car *models.Car) error
		expectedStatus   int
		expectedResponse any
	}{
		{
			name:             "invalid request body: missing comma",
			idParam:          "ABC123",
			body:             `{"make": "Chevrolet" "model": "Onix"}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "year is missing",
			idParam:          "ABC123",
			body:             `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan"}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "category is missing",
			idParam:          "ABC123",
			body:             `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "color is missing",
			idParam:          "ABC123",
			body:             `{"make":"Chevrolet", "model":"Onix", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "model is missing",
			idParam:          "ABC123",
			body:             `{"make":"Chevrolet", "color":"Gray", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "make is missing",
			idParam:          "ABC123",
			body:             `{"model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Invalid request body"},
		},
		{
			name:             "body id mismatch",
			idParam:          "ABC123",
			body:             `{"id": "DEF456", "make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "ID in body does not match URL ID"},
		},
		{
			name:    "car not found",
			idParam: "ABC123",
			body:    `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			updateFn: func(car *models.Car) error {
				return e.ErrCarNotFound
			},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: httpx.ErrorResponse{Message: "Car not found"},
		},
		{
			name:    "repository error",
			idParam: "ABC123",
			body:    `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			updateFn: func(car *models.Car) error {
				return errors.New("repository error")
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: httpx.ErrorResponse{Message: "Internal server error"},
		},
		{
			name:    "car updated successfully",
			idParam: "ABC123",
			body:    `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			updateFn: func(car *models.Car) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: httpx.SuccessResponse{
				Data: models.Car{
					ID:       "ABC123",
					Make:     "Chevrolet",
					Model:    "Onix",
					Color:    "Gray",
					Category: "Sedan",
					Year:     2025},
			},
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := NewCarController(
				services.NewCarService(
					&MockCarRepository{UpdateFn: tc.updateFn},
				),
			)

			router := chi.NewRouter()
			router.Route("/cars", func(r chi.Router) {
				r.Put("/{id:[A-Za-z0-9-]+}", controller.Update)
			})

			resp := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/cars/"+tc.idParam, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(resp, req)

			// Check status code
			if resp.Code != tc.expectedStatus {
				t.Fatalf("expected status %v, got %v", tc.expectedStatus, resp.Code)
			}

			body := resp.Body.Bytes()

			switch expected := tc.expectedResponse.(type) {
			case httpx.ErrorResponse:
				var got httpx.ErrorResponse
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}

				if got.Message != expected.Message {
					t.Fatalf("expected error %v, got %v", expected.Message, got.Message)
				}
			case httpx.SuccessResponse:
				var got httpx.SuccessResponse
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}

				gotDataBytes, err := json.Marshal(got.Data)
				if err != nil {
					t.Fatal(err)
				}

				var gotCar models.Car
				if err := json.Unmarshal(gotDataBytes, &gotCar); err != nil {
					t.Fatal(err)
				}

				expectedCar := tc.expectedResponse.(httpx.SuccessResponse).Data.(models.Car)
				if !reflect.DeepEqual(gotCar, expectedCar) {
					t.Fatalf("expected car %+v, got %+v", expectedCar, gotCar)
				}
			}
		})
	}
}
