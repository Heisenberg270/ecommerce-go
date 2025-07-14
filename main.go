package main

import (
	"log"
	"net/http"

	"github.com/Heisenberg270/ecommerce-go/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	initDB()

	r := chi.NewRouter()
	// health check
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
	})

	// product routes
	ph := handlers.NewProductHandler(db)
	r.Route("/products", func(r chi.Router) {
		r.Post("/", ph.Create)
		r.Get("/", ph.List)
		r.Get("/{id}", ph.Get)
		r.Put("/{id}", ph.Update)
		r.Delete("/{id}", ph.Delete)
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
