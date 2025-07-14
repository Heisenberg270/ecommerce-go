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

// OrderHandler manages orders
type OrderHandler struct {
	DB *sqlx.DB
}

// NewOrderHandler constructs an OrderHandler
func NewOrderHandler(db *sqlx.DB) *OrderHandler {
	return &OrderHandler{DB: db}
}

// CreateOrder handles POST /orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	// 1) Get authenticated user
	userID := r.Context().Value(ContextUserID).(int)

	// 2) Parse cart_id from JSON
	var in struct {
		CartID int `json:"cart_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// 3) Fetch cart items with prices
	var items []struct {
		models.CartItem
		UnitPrice float64 `db:"price"`
	}
	itemQ := `
    SELECT ci.cart_id, ci.product_id, ci.quantity, p.price
    FROM cart_items ci
    JOIN products p ON p.id = ci.product_id
    WHERE ci.cart_id = $1`
	if err := h.DB.Select(&items, itemQ, in.CartID); err != nil {
		http.Error(w, "failed to fetch cart items", http.StatusInternalServerError)
		return
	}
	if len(items) == 0 {
		http.Error(w, "cart is empty", http.StatusBadRequest)
		return
	}

	// 4) Compute total
	total := 0.0
	for _, it := range items {
		total += float64(it.Quantity) * it.UnitPrice
	}

	// 5) Insert into orders
	var order models.Order
	ordQ := `
    INSERT INTO orders (user_id, total_amount, status)
    VALUES ($1, $2, $3)
    RETURNING id, user_id, total_amount, status, created_at`
	if err := h.DB.Get(&order, ordQ, userID, total, "pending"); err != nil {
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	// 6) Insert order_items and clear cart in a tx
	tx, err := h.DB.Beginx()
	if err != nil {
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}
	for _, it := range items {
		if _, err := tx.Exec(
			`INSERT INTO order_items (order_id, product_id, quantity, unit_price)
       VALUES ($1, $2, $3, $4)`,
			order.ID, it.ProductID, it.Quantity, it.UnitPrice,
		); err != nil {
			tx.Rollback()
			http.Error(w, "failed to insert order items", http.StatusInternalServerError)
			return
		}
	}
	// clear the cart
	if _, err := tx.Exec(`DELETE FROM cart_items WHERE cart_id=$1`, in.CartID); err != nil {
		tx.Rollback()
		http.Error(w, "failed to clear cart", http.StatusInternalServerError)
		return
	}
	tx.Commit()

	// 7) Return the created order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// ListOrders handles GET /orders
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(ContextUserID).(int)
	var orders []models.Order
	if err := h.DB.Select(
		&orders,
		`SELECT id, user_id, total_amount, status, created_at
     FROM orders WHERE user_id=$1 ORDER BY id`,
		userID,
	); err != nil {
		http.Error(w, "failed to fetch orders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetOrder handles GET /orders/{orderID}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "orderID")
	orderID, _ := strconv.Atoi(idParam)

	// Fetch order
	var order models.Order
	if err := h.DB.Get(
		&order,
		`SELECT id, user_id, total_amount, status, created_at
     FROM orders WHERE id=$1`, orderID,
	); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "order not found", http.StatusNotFound)
		} else {
			http.Error(w, "failed to fetch order", http.StatusInternalServerError)
		}
		return
	}

	// Fetch items
	var items []struct {
		models.OrderItem
		ProductName string `db:"name" json:"product_name"`
	}
	itemQ := `
    SELECT oi.order_id, oi.product_id, oi.quantity, oi.unit_price, p.name
    FROM order_items oi
    JOIN products p ON p.id = oi.product_id
    WHERE oi.order_id = $1`
	if err := h.DB.Select(&items, itemQ, orderID); err != nil {
		http.Error(w, "failed to fetch order items", http.StatusInternalServerError)
		return
	}

	// Assemble response
	resp := struct {
		Order models.Order `json:"order"`
		Items interface{}  `json:"items"`
	}{Order: order, Items: items}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
