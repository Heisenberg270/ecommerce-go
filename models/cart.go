package models

import "time"

// Cart represents a user’s shopping cart.
type Cart struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// CartItem is a line item in a cart.
type CartItem struct {
	CartID    int `db:"cart_id" json:"cart_id"`
	ProductID int `db:"product_id" json:"product_id"`
	Quantity  int `db:"quantity" json:"quantity"`
}
