package routes

import (
	"cars/controllers"
	"cars/data"
	e "cars/pkg/errors"
	"cars/pkg/httpx"
	"cars/pkg/middleware"
	"cars/repositories"
	"cars/services"
	"net/http"
)

var mux = http.NewServeMux()

func Register() *http.ServeMux {
	repo := repositories.NewCarRepository(data.Cars())
	service := services.NewCarService(repo)
	cars := controllers.NewCarController(service)

	mux.Handle("/cars",
		middleware.Logging(
			handleCars(cars),
		),
	)
	mux.Handle("/cars/",
		middleware.Logging(
			handleCarByID(cars),
		),
	)
	return mux
}

func handleCars(cars *controllers.CarController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cars.List(w, r)
		case http.MethodPost:
			cars.Create(w, r)
		default:
			httpx.HandleServiceError(w, e.NewMethodNotAllowedError())
		}
	}
}

func handleCarByID(cars *controllers.CarController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cars.Get(w, r)
		case http.MethodPut:
			cars.Update(w, r)
		default:
			httpx.HandleServiceError(w, e.NewMethodNotAllowedError())
		}
	}
}
