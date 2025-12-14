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
		r.Post("/", cars.Create)

		r.Route("/{id:[A-Za-z0-9-]+}", func(r chi.Router) {
			r.Get("/", cars.Get)
			r.Put("/", cars.Update)
		})
	})
	return r
}
