package mappers

import (
	"cars/api/dto"
	c "cars/car"
)

// ToModelCreate maps a CreateCarRequest to a Car model.
//
// It is intended for create operations, where the ID is generated
// by the system and must not be set in the request.
func ToModelCreate(req dto.CreateCarRequest) *c.Car {
	return toModel(req)
}

// ToModelUpdate maps an UpdateCarRequest to a Car model,
// assigning the provided ID to the resulting entity.
//
// It is intended for update operations, where the ID identifies
// the existing resource being modified.
func ToModelUpdate(id string, req dto.UpdateCarRequest) *c.Car {
	car := toModel(req)
	car.ID = id
	return car
}

// ToModelPartialUpdate maps an PartialUpdateCarRequest to a CarUpdateParams model,
// assigning the provided ID to the resulting entity.
//
// It is intended for partial update operations, where the ID identifies
// the existing resource being modified.
func ToModelPartialUpdate(id string, req dto.PartialUpdateCarRequest) c.CarUpdateParams {
	return c.CarUpdateParams{
		ID:       id,
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

// toModel maps a CarUpsertRequest to a Car model.
//
// It is an internal helper shared by create and update mappings.
func toModel(req dto.CarUpsertRequest) *c.Car {
	return &c.Car{
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
func ToResponse(car *c.Car) dto.CarResponse {
	if car == nil {
		return dto.CarResponse{}
	}

	return dto.CarResponse{
		ID:       car.ID,
		Make:     car.Make,
		Model:    car.Model,
		Color:    car.Color,
		Category: car.Category,
		Year:     car.Year,
		Package:  car.Package,
		Mileage:  car.Mileage,
		Price:    car.Price,
	}
}

// ToResponseList maps a collection of Car models to a slice of CarResponse.
func ToResponseList(cars c.Cars) []dto.CarResponse {
	out := make([]dto.CarResponse, len(cars))
	for i := range cars {
		out[i] = ToResponse(&cars[i])
	}
	return out
}
