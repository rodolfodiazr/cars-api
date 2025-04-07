package main

import (
	"cars/routes"
	"log"
	"net/http"
)

func main() {
	r := routes.Register()

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
