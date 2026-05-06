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

// Register initializes and configures the application's HTTP routes.
//
// It sets up the dependency chain (repository → service → controller),
// applies global middlewares, and registers all endpoints under the "/cars" path.
//
// Routes:
//
//	GET    /cars          - List all cars (supports optional filtering via query params)
//	POST   /cars          - Create a new car
//	GET    /cars/{id}     - Retrieve a car by ID
//	PUT    /cars/{id}     - Replace an existing car (full update)
//	DELETE /cars/{id}     - Delete a car by ID
//
// Middleware applied:
//
//   - CleanPath: normalizes URL paths
//   - Recoverer: recovers from panics and returns HTTP 500
//   - Logging: custom request logging middleware
//
// Returns:
//
//	A configured *chi.Mux router ready to be used by an HTTP server.
func Register() *chi.Mux {
	repo := repositories.NewCarRepository(data.Cars())
	service := services.NewCarService(repo)
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
			// All fields must be provided; partial updates are not supported.
			r.Put("/", cars.Update)

			// DELETE /cars/{id}
			// Deletes a car by its ID.
			r.Delete("/", cars.Delete)
		})
	})
	return r
}
