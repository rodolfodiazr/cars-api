package services

import (
	"cars/models"
)

type FindFunc func(id string) (models.Car, error)
type ListFunc func(filters models.CarFilters) (models.Cars, error)
type CreateFunc func(car *models.Car) error
type UpdateFunc func(car *models.Car) error
type DeleteFunc func(id string) error

type MockCarRepository struct {
	FindFn   FindFunc
	ListFn   ListFunc
	CreateFn CreateFunc
	UpdateFn UpdateFunc
	DeleteFn DeleteFunc
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
