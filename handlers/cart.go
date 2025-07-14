package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/Heisenberg270/ecommerce-go/models"
)

// CartHandler holds the DB reference.
type CartHandler struct {
	DB *sqlx.DB
}

// NewCartHandler constructs a CartHandler.
func NewCartHandler(db *sqlx.DB) *CartHandler {
	return &CartHandler{DB: db}
}

// CreateCart creates a new cart for the authenticated user.
func (h *CartHandler) CreateCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ContextUserID).(int)
	var cart models.Cart
	err := h.DB.Get(&cart,
		`INSERT INTO carts (user_id) VALUES ($1) RETURNING id, user_id, created_at`,
		userID)
	if err != nil {
		http.Error(w, "failed to create cart", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(cart)
}

// AddItem adds or updates an item in the cart.
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	cartID, _ := strconv.Atoi(chi.URLParam(r, "cartID"))
	var in struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	_, err := h.DB.Exec(`
    INSERT INTO cart_items (cart_id, product_id, quantity)
    VALUES ($1, $2, $3)
    ON CONFLICT (cart_id, product_id) DO UPDATE
      SET quantity = cart_items.quantity + EXCLUDED.quantity
  `, cartID, in.ProductID, in.Quantity)
	if err != nil {
		http.Error(w, "failed to add item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetCart returns the cart and its items.
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	cartID, _ := strconv.Atoi(chi.URLParam(r, "cartID"))
	var cart models.Cart
	if err := h.DB.Get(&cart,
		`SELECT id, user_id, created_at FROM carts WHERE id=$1`, cartID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "cart not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to fetch cart", http.StatusInternalServerError)
		}
		return
	}

	var items []struct {
		models.CartItem
		ProductName string  `db:"name" json:"product_name"`
		Price       float64 `db:"price" json:"unit_price"`
	}
	query := `
    SELECT ci.cart_id, ci.product_id, ci.quantity,
           p.name, p.price
    FROM cart_items ci
    JOIN products p ON p.id=ci.product_id
    WHERE ci.cart_id=$1`
	if err := h.DB.Select(&items, query, cartID); err != nil {
		http.Error(w, "failed to fetch items", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Cart  models.Cart `json:"cart"`
		Items interface{} `json:"items"`
	}{Cart: cart, Items: items}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RemoveItem deletes an item from the cart.
func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	cartID, _ := strconv.Atoi(chi.URLParam(r, "cartID"))
	productID, _ := strconv.Atoi(chi.URLParam(r, "productID"))
	if _, err := h.DB.Exec(
		`DELETE FROM cart_items WHERE cart_id=$1 AND product_id=$2`,
		cartID, productID,
	); err != nil {
		http.Error(w, "failed to remove item", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
