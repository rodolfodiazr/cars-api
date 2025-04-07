package routes

import (
	"cars/controllers"
	"cars/models"
	e "cars/pkg/errors"
	"cars/pkg/middleware"
	"cars/repositories"
	"cars/services"
	"net/http"
)

var mux = http.NewServeMux()

var data = map[string]models.Car{
	"JHK290XJ": {
		ID:   "JHK290XJ",
		Make: "Ford", Model: "F10", Package: "Base",
		Color: "Silver", Year: 2010, Category: "Truck",
		Mileage: 120123, Price: 1999900,
	},
	"FWL37LA": {
		ID:   "FWL37LA",
		Make: "Toyota", Model: "Camry", Package: "SE",
		Color: "White", Year: 2019, Category: "Sedan",
		Mileage: 3999, Price: 2899000,
	},
	"1I3XJRLLC": {
		ID:   "1I3XJRLLC",
		Make: "Toyota", Model: "Rav4", Package: "XSE",
		Color: "Red", Year: 2018, Category: "SUV",
		Mileage: 24001, Price: 2275000,
	},
	"DKU43920S": {
		ID:   "DKU43920S",
		Make: "Ford", Model: "Bronco", Package: "Badlands",
		Color: "Bumt Orange", Year: 2022, Category: "SUV",
		Mileage: 1, Price: 4499000,
	},
}

func Register() *http.ServeMux {
	repo := repositories.NewCarRepository(data)
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
			http.Error(w, e.ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
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
			http.Error(w, e.ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		}
	}
}
