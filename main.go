package main

import (
	"cars/data"
	"cars/repositories"
	"cars/routes"
	"cars/services"
	"log"
	"net/http"
)

func main() {
	repo := repositories.NewCarRepository(data.Cars())
	service := services.NewCarService(repo)

	r := routes.Register(service)

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
