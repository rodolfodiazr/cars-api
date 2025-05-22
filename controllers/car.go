package controllers

import (
	"cars/models"
	e "cars/pkg/errors"
	"cars/pkg/httpx"
	"cars/pkg/logger"
	"cars/services"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// CarController manages HTTP requests related to cars.
type CarController struct {
	service services.CarService
}

// NewCarController creates a new instance of CarController.
func NewCarController(service services.CarService) *CarController {
	return &CarController{service: service}
}

// Get handles retrieving a car by its ID.
func (c *CarController) Get(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	id := c.getIDFromURL(r)
	if id == "" {
		log.Println("[ERROR] Car ID is required")
		httpx.Error(w, http.StatusBadRequest, e.ErrIDIsRequired.Error())
		return
	}

	car, err := c.service.Find(id)
	if err != nil {
		if errors.Is(err, e.ErrCarNotFound) {
			log.Printf("[ERROR] Car not found: Car ID: %v", id)
			httpx.Error(w, http.StatusNotFound, err.Error())
			return
		}

		log.Printf("[ERROR] Failed to retrieve car: %v", err)
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := httpx.JSON(w, http.StatusOK, car); err != nil {
		log.Printf("[ERROR] Failed to encode car response: %v", err)
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[INFO] Car with ID %s retrieved successfully", id)
}

// List handles retrieving all available cars.
func (c *CarController) List(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	cars, err := c.service.List()
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve car list: %v", err)
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := httpx.JSON(w, http.StatusOK, cars); err != nil {
		log.Printf("[ERROR] Failed to encode car list response: %v", err)
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[INFO] Successfully retrieved %d cars", len(cars))
}

// Create handles creating a new car.
func (c *CarController) Create(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	var car models.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		log.Printf("[ERROR] Invalid request body: %v", err)
		httpx.Error(w, http.StatusBadRequest, e.ErrInvalidRequestBody.Error())
		return
	}

	if err := car.Validate(); err != nil {
		log.Printf("[ERROR] Validation error: %v", err)
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.service.Create(&car); err != nil {
		log.Printf("[ERROR] Failed to create car: %v", err)
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := httpx.JSON(w, http.StatusCreated, car); err != nil {
		log.Printf("[ERROR] Failed to encode created car response: %v", err)
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[INFO] Car with ID %s has been created successfully", car.ID)
}

// Update handles updating an existing car.
func (c *CarController) Update(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	id := c.getIDFromURL(r)
	if id == "" {
		log.Println("[ERROR] Car ID is required")
		httpx.Error(w, http.StatusBadRequest, e.ErrIDIsRequired.Error())
		return
	}

	var car models.Car
	if err := json.NewDecoder(r.Body).Decode(&car); err != nil {
		log.Printf("[ERROR] Invalid request body: %v", err)
		httpx.Error(w, http.StatusBadRequest, e.ErrInvalidRequestBody.Error())
		return
	}

	if err := car.Validate(); err != nil {
		log.Printf("[ERROR] Validation error: %v", err)
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	car.ID = id // Put Car ID in the struct
	if err := c.service.Update(&car); err != nil {
		if errors.Is(err, e.ErrCarNotFound) {
			log.Printf("[ERROR] Car not found: Car ID: %v", id)
			httpx.Error(w, http.StatusNotFound, err.Error())
			return
		}

		log.Printf("[ERROR] Failed to update car: %v", err)
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := httpx.JSON(w, http.StatusOK, car); err != nil {
		log.Printf("[ERROR] Failed to encode updated car response: %v", err)
		httpx.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[INFO] Car with ID %s has been updated successfully", id)
}

// getIDFromURL extracts the car ID from the request URL.
func (c *CarController) getIDFromURL(r *http.Request) string {
	id := strings.TrimPrefix(r.URL.Path, "/cars/")
	if strings.Contains(id, "/") || id == "" {
		fmt.Println("ID: ", id)
		return ""
	}
	return id
}
