package dto

import "cars/models"

func ToModelCreate(req CreateCarRequest) *models.Car {
	return &models.Car{
		Make:  req.Make,
		Model: req.Model,
		Year:  req.Year,
	}
}

func ToModelUpdate(id string, req UpdateCarRequest) *models.Car {
	return &models.Car{
		ID:     id,
		BodyID: req.ID,
		Make:   req.Make,
		Model:  req.Model,
		Year:   req.Year,
	}
}

func ToResponse(car models.Car) CarResponse {
	return CarResponse{
		ID:    car.ID,
		Make:  car.Make,
		Model: car.Model,
		Year:  car.Year,
	}
}

func ToResponseList(cars models.Cars) []CarResponse {
	out := make([]CarResponse, 0, len(cars))
	for _, car := range cars {
		out = append(out, ToResponse(car))
	}
	return out
}
