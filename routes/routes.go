package routes

import (
	"cars/controllers"
	"cars/data"
	"cars/pkg/middleware"
	"cars/repositories"
	"cars/services"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func Register() *chi.Mux {
	repo := repositories.NewCarRepository(data.Cars())
	service := services.NewCarService(repo)
	cars := controllers.NewCarController(service)

	r := chi.NewRouter()

	r.Use(chimw.CleanPath)
	r.Use(chimw.Recoverer)

	r.Use(middleware.Logging)

	r.Route("/cars", func(r chi.Router) {
		r.Get("/", cars.List)

		// POST /cars
		// Creates a new car.
		// All required fields must be provided in the request body.
		r.Post("/", cars.Create)

		r.Route("/{id:[A-Za-z0-9-]+}", func(r chi.Router) {
			// GET /cars/{id}
			// Retrieves a car by its ID.
			r.Get("/", cars.Get)

			// PUT /cars/{id}
			// Performs a full replacement of the car resource.
			// All fields must be provided; partial updates are not supported.
			r.Put("/", cars.Update)
		})
	})
	return r
}
