package routes

import (
	"cars/controllers"
	"cars/pkg/middleware"
	"cars/services"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// Register initializes and configures the application's HTTP routes.
//
// It creates the car controller using the provided service,
// applies global middlewares, and registers all car endpoints.
//
// Returns a configured *chi.Mux router.
func Register(service services.CarService) *chi.Mux {
	cars := controllers.NewCarController(service)

	r := chi.NewRouter()

	r.Use(chimw.CleanPath)
	r.Use(chimw.Recoverer)

	r.Use(middleware.Logging)

	r.Route("/cars", func(r chi.Router) {
		// GET /cars
		// Retrieves a list of cars.
		// Supports optional query parameters for filtering.
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
			// All fields must be provided.
			r.Put("/", cars.Update)

			// PATCH /cars/{id}
			// Partially updates the car resource.
			// Only the provided fields are updated; all other fields remain unchanged.
			r.Patch("/", cars.PartialUpdate)

			// DELETE /cars/{id}
			// Deletes a car by its ID.
			r.Delete("/", cars.Delete)
		})
	})
	return r
}
