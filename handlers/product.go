package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/Heisenberg270/ecommerce-go/models"
)

// ProductHandler holds DB reference
type ProductHandler struct {
	DB *sqlx.DB
}

// NewProductHandler returns a handler with DB injected
func NewProductHandler(db *sqlx.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

// Create handles POST /products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// If table is empty, restart the ID sequence
	var cnt int
	if err := h.DB.Get(&cnt, "SELECT COUNT(*) FROM products"); err != nil {
		http.Error(w, "Failed to check products count", http.StatusInternalServerError)
		return
	}
	if cnt == 0 {
		if _, err := h.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1"); err != nil {
			http.Error(w, "Failed to reset product ID sequence", http.StatusInternalServerError)
			return
		}
	}

	query := `INSERT INTO products (name, description, price)
            VALUES ($1, $2, $3)
            RETURNING id, created_at, updated_at`
	if err := h.DB.QueryRowx(query, p.Name, p.Description, p.Price).
		StructScan(&p); err != nil {
		http.Error(w, "Failed to insert product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// List handles GET /products
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products := []models.Product{}
	if err := h.DB.Select(&products, "SELECT * FROM products ORDER BY id"); err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// Get handles GET /products/{id}
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	var p models.Product
	if err := h.DB.Get(&p, "SELECT * FROM products WHERE id=$1", id); err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// Update handles PUT /products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = h.DB.Exec(
		`UPDATE products SET name=$1, description=$2, price=$3, updated_at=now() WHERE id=$4`,
		p.Name, p.Description, p.Price, id,
	)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Delete handles DELETE /products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	if _, err := h.DB.Exec("DELETE FROM products WHERE id=$1", id); err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
