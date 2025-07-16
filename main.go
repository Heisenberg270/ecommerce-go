package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Heisenberg270/ecommerce-go/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	initDB()

	r := chi.NewRouter()
	// CORS — allow your frontend dev server to talk to us
	r.Use(cors.Handler(cors.Options{
		// put your actual domains here in production
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // 5 minutes
	}))

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
	// Protected routes: Carts & Orders
	r.Group(func(r chi.Router) {
		r.Use(handlers.AuthMiddleware(jwtSecret))

		// Cart
		ch := handlers.NewCartHandler(db)
		r.Post("/carts", ch.CreateCart)
		r.Route("/carts/{cartID}", func(r chi.Router) {
			r.Post("/items", ch.AddItem)
			r.Get("/", ch.GetCart)
			r.Delete("/items/{productID}", ch.RemoveItem)
		})
		// Orders
		oh := handlers.NewOrderHandler(db)
		r.Post("/orders", oh.CreateOrder)
		r.Get("/orders", oh.ListOrders)
		r.Get("/orders/{orderID}", oh.GetOrder)
	})

	log.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
