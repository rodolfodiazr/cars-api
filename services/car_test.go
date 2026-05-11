package services

import "cars/models"

type MockCarRepository struct {
	FindFunc func(id string) (models.Car, error)
}

func (m *MockCarRepository) Find(id string) (models.Car, error) {
	return m.FindFunc(id)
}
