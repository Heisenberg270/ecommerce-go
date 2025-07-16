package models

import "time"

// Order represents a completed (or pending) purchase.
type Order struct {
	ID          int       `db:"id" json:"id"`
	UserID      int       `db:"user_id" json:"user_id"`
	TotalAmount float64   `db:"total_amount" json:"total_amount"`
	Status      string    `db:"status" json:"status"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// OrderItem is a line item within an order.
type OrderItem struct {
	OrderID   int     `db:"order_id" json:"order_id"`
	ProductID int     `db:"product_id" json:"product_id"`
	Quantity  int     `db:"quantity" json:"quantity"`
	UnitPrice float64 `db:"unit_price" json:"unit_price"`
}
