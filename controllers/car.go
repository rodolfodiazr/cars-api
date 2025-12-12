package controllers

import (
	"cars/models"
	e "cars/pkg/errors"
	"cars/pkg/httpx"
	"cars/pkg/logger"
	"cars/services"
	"fmt"
	"net/http"
	"strconv"
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

	id, err := c.getIDFromURL(r)
	if err != nil {
		log.Print(err)
		httpx.HandleServiceError(w, err)
		return
	}

	car, err := c.service.Find(id)
	if err != nil {
		log.Printf("error retrieving car id=%s: %v", id, err)
		httpx.HandleServiceError(w, err)
		return
	}

	if err := httpx.JSON(w, http.StatusOK, car); err != nil {
		log.Printf("error encoding car response: %v", err)
		httpx.HandleServiceError(w, err)
		return
	}

	log.Printf("car retrieved id=%s", id)
}

// List handles retrieving all available cars.
func (c *CarController) List(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	filters, err := c.parseCarFilters(r)
	if err != nil {
		log.Printf("error parsing car filters: %v", err)
		httpx.HandleServiceError(w, err)
		return
	}

	cars, err := c.service.List(filters)
	if err != nil {
		log.Printf("error retrieving cars: %v", err)
		httpx.HandleServiceError(w, err)
		return
	}

	if err := httpx.JSON(w, http.StatusOK, cars); err != nil {
		log.Printf("error encoding cars response: %v", err)
		httpx.HandleServiceError(w, err)
		return
	}

	log.Printf("%d cars retrieved", len(cars))
}

// Create handles creating a new car.
func (c *CarController) Create(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	car, err := httpx.DecodeAndValidate[models.Car](r)
	if err != nil {
		log.Printf("error decoding car payload: %v", err)
		httpx.HandleServiceError(w, err)
		return
	}

	if err := c.service.Create(car); err != nil {
		log.Printf("error creating car: %v", err)
		httpx.HandleServiceError(w, err)
		return
	}

	if err := httpx.JSON(w, http.StatusCreated, car); err != nil {
		log.Printf("error encoding created car response: %v", err)
		httpx.HandleServiceError(w, err)
		return
	}

	log.Printf("car created id=%s", car.ID)
}

// Update handles updating an existing car.
func (c *CarController) Update(w http.ResponseWriter, r *http.Request) {
	log := logger.FromContext(r.Context())

	id, err := c.getIDFromURL(r)
	if err != nil {
		log.Print(err)
		httpx.HandleServiceError(w, err)
		return
	}

	car, err := httpx.DecodeAndValidate[models.Car](r)
	if err != nil {
		log.Printf("error decoding car payload: %v", err)
		httpx.HandleServiceError(w, err)
		return
	}

	car.BodyID = car.ID // Preserve original
	car.ID = id         // Override with URL ID

	if err := c.service.Update(car); err != nil {
		log.Printf("error updating car id=%s: %v", car.ID, err)
		httpx.HandleServiceError(w, err)
		return
	}

	if err := httpx.JSON(w, http.StatusOK, car); err != nil {
		log.Printf("error encoding updated car response: %v", err)
		httpx.HandleServiceError(w, err)
		return
	}

	log.Printf("car updated id=%s", id)
}

// getIDFromURL extracts the car ID from the request URL.
func (c *CarController) getIDFromURL(r *http.Request) (string, error) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[0] != "cars" {
		return "", e.NewInvalidCarPathError()
	}

	return parts[1], nil
}

// parseCarFilters converts URL query parameters into a CarFilters struct.
func (c *CarController) parseCarFilters(r *http.Request) (models.CarFilters, error) {
	var filters models.CarFilters

	filters.Make = r.URL.Query().Get("make")
	filters.Model = r.URL.Query().Get("model")

	if yearStr := r.URL.Query().Get("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return filters, fmt.Errorf("invalid year: %q", yearStr)
		}
		filters.Year = year
	}
	return filters, nil
}
