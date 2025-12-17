package controllers

import (
	"cars/models"
	"cars/pkg/httpx"
	"cars/pkg/logger"
	"cars/services"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
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

	id := chi.URLParam(r, "id")

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

	filters, err := parseCarFilters(r.URL.Query())
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

	id := chi.URLParam(r, "id")

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

// parseCarFilters converts URL query parameters into a CarFilters struct.
func parseCarFilters(q url.Values) (models.CarFilters, error) {
	f := models.CarFilters{
		Make:  q.Get("make"),
		Model: q.Get("model"),
	}

	if yearStr := q.Get("year"); yearStr != "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return f, fmt.Errorf("invalid year: %q", yearStr)
		}
		f.Year = year
	}

	return f, nil
}
