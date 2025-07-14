package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Heisenberg270/ecommerce-go/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	initDB()

	r := chi.NewRouter()

	// Auth routes
	jwtSecret := os.Getenv("JWT_SECRET")
	ah := handlers.NewAuthHandler(db, jwtSecret)
	r.Post("/users/signup", ah.Signup)
	r.Post("/users/login", ah.Login)

	// Health check
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
	})

	// Product routes
	ph := handlers.NewProductHandler(db)
	r.Route("/products", func(r chi.Router) {
		r.Post("/", ph.Create)
		r.Get("/", ph.List)
		r.Get("/{id}", ph.Get)
		r.Put("/{id}", ph.Update)
		r.Delete("/{id}", ph.Delete)
	})

	// Cart routes (protected)
	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware(jwtSecret))

		ch := handlers.NewCartHandler(db)
		r.Post("/carts", ch.CreateCart)
		r.Route("/carts/{cartID}", func(r chi.Router) {
			r.Post("/items", ch.AddItem)
			r.Get("/", ch.GetCart)
			r.Delete("/items/{productID}", ch.RemoveItem)
		})
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
