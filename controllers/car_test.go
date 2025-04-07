package controllers

import (
	"cars/models"
	e "cars/pkg/errors"
	"cars/pkg/utils"
	"cars/repositories"
	"cars/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockCarRepository struct {
	cars map[string]models.Car
}

func NewMockCarRepository(data map[string]models.Car) repositories.CarRepository {
	repo := &MockCarRepository{
		cars: data,
	}
	return repo
}

func (m *MockCarRepository) Find(id string) (models.Car, error) {
	car, exists := m.cars[id]
	if !exists {
		return models.Car{}, e.ErrCarNotFound
	}
	return car, nil
}

func (m *MockCarRepository) List() (models.Cars, error) {
	var list []models.Car
	for _, car := range m.cars {
		list = append(list, car)
	}
	return list, nil
}

func (m *MockCarRepository) Create(car *models.Car) error {
	id, err := utils.GenerateID()
	if err != nil {
		return err
	}

	car.ID = id
	m.cars[car.ID] = *car
	return nil
}

func (m *MockCarRepository) Update(car *models.Car) error {
	if _, exists := m.cars[car.ID]; !exists {
		return e.ErrCarNotFound
	}

	m.cars[car.ID] = *car
	return nil
}

func setupController() *CarController {
	data := map[string]models.Car{
		"ABC123": {
			ID:       "ABC123",
			Make:     "Chevrolet",
			Model:    "Onix",
			Package:  "ABC",
			Color:    "Black",
			Category: "Sedan",
			Year:     2025,

			Mileage: 150,
			Price:   23000000,
		},
		"DEF456": {
			ID:       "DEF456",
			Make:     "Toyota",
			Model:    "Yaris",
			Package:  "DEF",
			Color:    "Red",
			Category: "Sedan",
			Year:     2025,

			Mileage: 150,
			Price:   23000000,
		},
		"GHI789": {
			ID:       "GHI789",
			Make:     "Renault",
			Model:    "Arkana",
			Package:  "GHI",
			Color:    "White",
			Category: "Sedan",
			Year:     2025,

			Mileage: 150,
			Price:   23000000,
		},
	}

	return NewCarController(
		services.NewCarService(
			NewMockCarRepository(
				data,
			),
		),
	)
}

func Test_Car_Get(t *testing.T) {
	controller := setupController()

	tCases := []struct {
		index   string
		idParam string

		expectedStatus int
		expectedBody   models.Car
	}{
		{
			index:   "Case 1",
			idParam: "ABC123",

			expectedStatus: http.StatusOK,
			expectedBody: models.Car{
				ID:       "ABC123",
				Make:     "Chevrolet",
				Model:    "Onix",
				Color:    "Black",
				Category: "Sedan",
				Year:     2025,
			},
		},
		{
			index:   "Case 2",
			idParam: "DEF456",

			expectedStatus: http.StatusOK,
			expectedBody: models.Car{
				ID:       "DEF456",
				Make:     "Toyota",
				Model:    "Yaris",
				Package:  "DEF",
				Color:    "Red",
				Category: "Sedan",
				Year:     2025,
			},
		},
		{
			index:   "Case 3",
			idParam: "XYZ000",

			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tCase := range tCases {
		req, err := http.NewRequest("GET", "/cars/"+tCase.idParam, nil)
		if err != nil {
			t.Fatalf("[%s] Failed to create request: %v", tCase.index, err)
		}

		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.Get)
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != tCase.expectedStatus {
			t.Errorf("[%s] Expected status %v, got %v", tCase.index, tCase.expectedStatus, status)
		}

		if tCase.expectedBody.ID != "" {
			var car models.Car
			if err := json.Unmarshal(resp.Body.Bytes(), &car); err != nil {
				t.Fatalf("[%s]: %v", tCase.index, err)
			}

			if car.ID != tCase.expectedBody.ID {
				t.Errorf("[%s] Expected Car ID %v, got %v", tCase.index, tCase.expectedBody.ID, car.ID)
			}
		}
	}
}

func Test_Car_List(t *testing.T) {
	controller := setupController()

	req, err := http.NewRequest("GET", "/cars", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(controller.List)
	handler.ServeHTTP(resp, req)
	if status := resp.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", status)
	}

	var cars models.Cars
	if err := json.Unmarshal(resp.Body.Bytes(), &cars); err != nil {
		t.Fatal(err)
	}

	if len(cars) != 3 {
		t.Errorf("Expected 3 cars, got %v", len(cars))
	}
}

func Test_Car_Create(t *testing.T) {
	controller := setupController()

	tCases := []struct {
		index string
		body  string

		expectedStatus int
		expectedData   models.Car
	}{
		{
			index: "Case 1",
			body:  `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan"}`,

			expectedStatus: http.StatusBadRequest,
			expectedData:   models.Car{},
		},
		{
			index: "Case 2",
			body:  `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "year":2025}`,

			expectedStatus: http.StatusBadRequest,
			expectedData:   models.Car{},
		},
		{
			index: "Case 3",
			body:  `{"make":"Chevrolet", "model":"Onix", "category":"Sedan", "year":2025}`,

			expectedStatus: http.StatusBadRequest,
			expectedData:   models.Car{},
		},
		{
			index: "Case 4",
			body:  `{"make":"Chevrolet", "color":"Gray", "category":"Sedan", "year":2025}`,

			expectedStatus: http.StatusBadRequest,
			expectedData:   models.Car{},
		},
		{
			index: "Case 5",
			body:  `{"make":"Chevrolet", "model":"Onix", "color":"Gray", "category":"Sedan", "year":2025}`,

			expectedStatus: http.StatusCreated,
			expectedData: models.Car{
				Make:     "Chevrolet",
				Model:    "Onix",
				Color:    "Gray",
				Category: "Sedan",
				Year:     2025,
			},
		},
	}

	for _, tCase := range tCases {
		req, err := http.NewRequest("POST", "/cars", strings.NewReader(tCase.body))
		if err != nil {
			t.Fatalf("[%s] Failed to create request: %v", tCase.index, err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.Create)
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != tCase.expectedStatus {
			t.Errorf("[%s] Expected status %v, got %v", tCase.index, tCase.expectedStatus, status)
		}

		if tCase.expectedData.Make != "" {
			var car models.Car
			if err := json.Unmarshal(resp.Body.Bytes(), &car); err != nil {
				t.Fatalf("[%s]: %v", tCase.index, err)
			}

			if car.Make != tCase.expectedData.Make {
				t.Errorf("[%s] Expected Car Make %v, got %v", tCase.index, tCase.expectedData.Make, car.Make)
			}
		}
	}
}

func Test_Car_Update(t *testing.T) {
	controller := setupController()

	tCases := []struct {
		index   string
		idParam string
		body    string

		expectedStatus int
		expectedData   models.Car
	}{
		{
			index:   "Case 1",
			idParam: "ABC123",
			body:    `{"make":"Chevrolet", "model":"Onix", "color":"Black", "category":"Sedan"}`,

			expectedStatus: http.StatusBadRequest,
		},
		{
			index:   "Case 2",
			idParam: "ABC123",
			body:    `{"make":"Chevrolet", "model":"Onix", "color":"Black", "year":2025}`,

			expectedStatus: http.StatusBadRequest,
		},
		{
			index:   "Case 3",
			idParam: "ABC123",
			body:    `{"make":"Chevrolet", "model":"Onix", "category":"Sedan", "year":2025}`,

			expectedStatus: http.StatusBadRequest,
		},
		{
			index:   "Case 4",
			idParam: "ABC123",
			body:    `{"make":"Chevrolet", "color":"Black", "category":"Sedan", "year":2025}`,

			expectedStatus: http.StatusBadRequest,
		},
		{
			index:   "Case 5",
			idParam: "ABC123",
			body:    `{"make":"Ford", "model":"F-150", "color":"Red", "category":"Pickup", "year":2024}`,

			expectedStatus: http.StatusOK,
			expectedData: models.Car{
				ID:       "ABC123",
				Make:     "Ford",
				Model:    "F-150",
				Color:    "Red",
				Category: "Pickup",
				Year:     2024,
			},
		},
	}

	for _, tCase := range tCases {
		req, err := http.NewRequest("PUT", "/cars/"+tCase.idParam, strings.NewReader(tCase.body))
		if err != nil {
			t.Fatalf("[%s] Failed to create request: %v", tCase.index, err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(controller.Update)
		handler.ServeHTTP(resp, req)

		if status := resp.Code; status != tCase.expectedStatus {
			t.Errorf("[%s] Expected status %v, got %v", tCase.index, tCase.expectedStatus, status)
		}

		if tCase.expectedData.ID != "" {
			var car models.Car
			if err := json.Unmarshal(resp.Body.Bytes(), &car); err != nil {
				t.Fatalf("[%s]: %v", tCase.index, err)
			}

			if car.ID != tCase.expectedData.ID {
				t.Errorf("[%s] Expected Car ID %v, got %v", tCase.index, tCase.expectedData.ID, car.ID)
			}

			if car.Make != tCase.expectedData.Make {
				t.Errorf("[%s] Expected Car Make %v, got %v", tCase.index, tCase.expectedData.Make, car.Make)
			}
		}

	}
}
