package dto

import (
	"cars/models"
)

// ToModelCreate maps a CreateCarRequest to a Car model.
//
// It is intended for create operations, where the ID is generated
// by the system and must not be set in the request.
func ToModelCreate(req CreateCarRequest) *models.Car {
	return toModel(req)
}

// ToModelUpdate maps an UpdateCarRequest to a Car model,
// assigning the provided ID to the resulting entity.
//
// It is intended for update operations, where the ID identifies
// the existing resource being modified.
func ToModelUpdate(id string, req UpdateCarRequest) *models.Car {
	car := toModel(req)
	car.ID = id
	return car
}

// toModel maps a CarUpsertRequest to a Car model.
//
// It is an internal helper shared by create and update mappings.
func toModel(req CarUpsertRequest) *models.Car {
	return &models.Car{
		Make:     req.Make,
		Model:    req.Model,
		Color:    req.Color,
		Category: req.Category,
		Year:     req.Year,
		Package:  req.Package,
		Mileage:  req.Mileage,
		Price:    req.Price,
	}
}

// ToResponse maps a Car model to a CarResponse.
//
// If the provided car is nil, it returns an empty response.
func ToResponse(car *models.Car) CarResponse {
	if car == nil {
		return CarResponse{}
	}

	return CarResponse{
		ID:       car.ID,
		Make:     car.Make,
		Model:    car.Model,
		Package:  car.Package,
		Color:    car.Color,
		Category: car.Category,
		Year:     car.Year,
		Mileage:  car.Mileage,
		Price:    car.Price,
	}
}

// ToResponseList maps a collection of Car models to a slice of CarResponse.
func ToResponseList(cars models.Cars) []CarResponse {
	out := make([]CarResponse, len(cars))
	for i := range cars {
		out[i] = ToResponse(&cars[i])
	}
	return out
}
