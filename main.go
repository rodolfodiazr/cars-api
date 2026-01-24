package main

import (
	"cars/routes"
	"log"
	"net/http"
)

func main() {
	r := routes.Register()

	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
