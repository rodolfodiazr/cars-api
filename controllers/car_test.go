package controllers

import (
	"cars/api/dto"
	"cars/models"
	e "cars/pkg/errors"
	"cars/pkg/httpx"
	u "cars/pkg/utils"
	"cars/services"
	"encoding/json"
	"errors"
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
	DeleteFn func(id string) error
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
			expectedResponse: models.Car{
				ID:       "ABC123",
				Make:     "Chevrolet",
				Model:    "Onix",
				Color:    "Black",
				Category: "Sedan",
				Year:     2025,
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
			case models.Car:
				var gotCar models.Car
				if err := json.Unmarshal(body, &gotCar); err != nil {
					t.Fatal(err)
				}

				expectedCar := tc.expectedResponse
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
			name: "list all cars",
			listFn: func(f models.CarFilters) (models.Cars, error) {
				return models.Cars{
					{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Package: u.Ptr("ABC"), Color: "Black", Category: "Sedan", Year: 2025},
					{ID: "DEF456", Make: "Toyota", Model: "Yaris", Package: u.Ptr("DEF"), Color: "Red", Category: "Sedan", Year: 2025},
					{ID: "GHI789", Make: "Renault", Model: "Arkana", Package: u.Ptr("GHI"), Color: "White", Category: "Sedan", Year: 2025},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []dto.CarResponse{
				{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Package: u.Ptr("ABC"), Color: "Black", Category: "Sedan", Year: 2025},
				{ID: "DEF456", Make: "Toyota", Model: "Yaris", Package: u.Ptr("DEF"), Color: "Red", Category: "Sedan", Year: 2025},
				{ID: "GHI789", Make: "Renault", Model: "Arkana", Package: u.Ptr("GHI"), Color: "White", Category: "Sedan", Year: 2025},
			},
		},
		{
			name:        "list filtered by make",
			queryParams: "?make=Toyota",
			listFn: func(f models.CarFilters) (models.Cars, error) {
				if f.Make != "Toyota" {
					t.Errorf("expected make %q, got %q", "Toyota", f.Make)
				}
				return models.Cars{
					{ID: "DEF456", Make: "Toyota", Model: "Yaris", Package: u.Ptr("DEF"), Color: "Red", Category: "Sedan", Year: 2025},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []dto.CarResponse{
				{ID: "DEF456", Make: "Toyota", Model: "Yaris", Package: u.Ptr("DEF"), Color: "Red", Category: "Sedan", Year: 2025},
			},
		},
		{
			name:        "list filtered by model",
			queryParams: "?model=Arkana",
			listFn: func(f models.CarFilters) (models.Cars, error) {
				if f.Model != "Arkana" {
					t.Errorf("expected model %q, got %q", "Arkana", f.Model)
				}
				return models.Cars{
					{ID: "GHI789", Make: "Renault", Model: "Arkana", Package: u.Ptr("GHI"), Color: "White", Category: "Sedan", Year: 2025},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []dto.CarResponse{
				{ID: "GHI789", Make: "Renault", Model: "Arkana", Package: u.Ptr("GHI"), Color: "White", Category: "Sedan", Year: 2025},
			},
		},
		{
			name:        "list filtered by year",
			queryParams: "?year=2025",
			listFn: func(f models.CarFilters) (models.Cars, error) {
				if f.Year == nil {
					t.Fatal("expected year filter")
				}
				if *f.Year != 2025 {
					t.Errorf("expected year %d, got %d", 2025, *f.Year)
				}
				return models.Cars{
					{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Color: "Black", Category: "Sedan", Year: 2025},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []dto.CarResponse{
				{ID: "ABC123", Make: "Chevrolet", Model: "Onix", Color: "Black", Category: "Sedan", Year: 2025},
			},
		},
		{
			name:             "invalid year",
			queryParams:      "?year=MMXXV",
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Validation failed"},
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

			req := httptest.NewRequest(http.MethodGet, "/cars"+tc.queryParams, nil)

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			// Check status code
			if resp.Code != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d", tc.expectedStatus, resp.Code)
			}

			switch expected := tc.expectedResponse.(type) {
			case httpx.ErrorResponse:
				var got httpx.ErrorResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}

				if got.Message != expected.Message {
					t.Errorf("expected message %q, got %q", expected.Message, got.Message)
				}
			case []dto.CarResponse:
				var got []dto.CarResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if !reflect.DeepEqual(got, expected) {
					t.Errorf("expected response %+v, got %+v", expected, got)
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
		expectedLocation string
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
			expectedResponse: httpx.ErrorResponse{Message: "Validation failed"},
		},
		{
			name:             "category is missing",
			body:             `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Validation failed"},
		},
		{
			name:             "color is missing",
			body:             `{"make":"Chevrolet", "model":"Onix", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Validation failed"},
		},
		{
			name:             "model is missing",
			body:             `{"make":"Chevrolet", "color":"Gray", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Validation failed"},
		},
		{
			name:             "make is missing",
			body:             `{"model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Validation failed"},
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
			body: `{"make":"Chevrolet","model":"Onix","color":"Gray","category":"Sedan","year":2025}`,
			createFn: func(car *models.Car) error {
				car.ID = "A1"
				return nil
			},
			expectedStatus:   http.StatusCreated,
			expectedLocation: "/cars/A1",
			expectedResponse: dto.CarResponse{
				ID:       "A1",
				Make:     "Chevrolet",
				Model:    "Onix",
				Color:    "Gray",
				Category: "Sedan",
				Year:     2025,
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

			if tc.expectedLocation != "" {
				gotLocation := resp.Header().Get("Location")
				if gotLocation != tc.expectedLocation {
					t.Fatalf("expected Location %q, got %q", tc.expectedLocation, gotLocation)
				}
			}

			switch expected := tc.expectedResponse.(type) {
			case httpx.ErrorResponse:
				var got httpx.ErrorResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}

				if got.Message != expected.Message {
					t.Fatalf("expected error %v, got %v", expected.Message, got.Message)
				}
			case dto.CarResponse:
				var got dto.CarResponse
				if err := json.Unmarshal(resp.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}

				if !reflect.DeepEqual(got, expected) {
					t.Fatalf("expected car %+v, got %+v", expected, got)
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
			expectedResponse: httpx.ErrorResponse{Message: "Validation failed"},
		},
		{
			name:             "category is missing",
			idParam:          "ABC123",
			body:             `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Validation failed"},
		},
		{
			name:             "color is missing",
			idParam:          "ABC123",
			body:             `{"make":"Chevrolet", "model":"Onix", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Validation failed"},
		},
		{
			name:             "model is missing",
			idParam:          "ABC123",
			body:             `{"make":"Chevrolet", "color":"Gray", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Validation failed"},
		},
		{
			name:             "make is missing",
			idParam:          "ABC123",
			body:             `{"model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: httpx.ErrorResponse{Message: "Validation failed"},
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
			expectedResponse: models.Car{
				ID:       "ABC123",
				Make:     "Chevrolet",
				Model:    "Onix",
				Color:    "Gray",
				Category: "Sedan",
				Year:     2025},
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
			case models.Car:
				var gotCar models.Car
				if err := json.Unmarshal(body, &gotCar); err != nil {
					t.Fatal(err)
				}

				expectedCar := tc.expectedResponse.(models.Car)
				if !reflect.DeepEqual(gotCar, expectedCar) {
					t.Fatalf("expected car %+v, got %+v", expectedCar, gotCar)
				}
			}
		})
	}
}
