package dto

import (
	"cars/models"
)

func ToModelCreate(req CreateCarRequest) *models.Car {
	return toModel(req)
}

func ToModelUpdate(id string, req UpdateCarRequest) *models.Car {
	car := toModel(req)
	car.ID = id
	return car
}

func toModel(req CarUpsertRequest) *models.Car {
	return &models.Car{
		Make:     req.Make,
		Model:    req.Model,
		Package:  req.Package,
		Color:    req.Color,
		Category: req.Category,
		Year:     req.Year,
		Mileage:  req.Mileage,
		Price:    req.Price,
	}
}

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

func ToResponseList(cars models.Cars) []CarResponse {
	out := make([]CarResponse, 0, len(cars))
	for i := range cars {
		out = append(out, ToResponse(&cars[i]))
	}
	return out
}
