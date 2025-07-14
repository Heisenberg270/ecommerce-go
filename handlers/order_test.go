package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/Heisenberg270/ecommerce-go/models"
)

func setupOrderMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func TestCreateOrder(t *testing.T) {
	db, mock := setupOrderMock(t)
	// 1) Cart items fetch
	rows := sqlmock.NewRows([]string{"cart_id", "product_id", "quantity", "price"}).
		AddRow(1, 10, 2, 5.00)
	mock.ExpectQuery(`FROM cart_items`).
		WithArgs(1).
		WillReturnRows(rows)

	// 2) Insert order
	now := time.Now()
	mock.ExpectQuery(`INSERT INTO orders`).
		WithArgs(42, 10.00, "pending").
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "user_id", "total_amount", "status", "created_at"},
		).AddRow(100, 42, 10.00, "pending", now))

	// 3) Begin transaction
	mock.ExpectBegin()

	// 4) Insert order_items
	mock.ExpectExec(`INSERT INTO order_items`).
		WithArgs(100, 10, 2, 5.00).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 5) Clear cart_items
	mock.ExpectExec(`DELETE FROM cart_items`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 6) Commit transaction
	mock.ExpectCommit()

	// Call handler
	oh := NewOrderHandler(db)
	payload := []byte(`{"cart_id":1}`)
	req := httptest.NewRequest("POST", "/orders", bytes.NewReader(payload))
	req = req.WithContext(context.WithValue(req.Context(), ContextUserID, 42))
	w := httptest.NewRecorder()

	oh.CreateOrder(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("CreateOrder status = %d; want %d", w.Code, http.StatusCreated)
	}
	var resp models.Order
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp.ID != 100 || resp.UserID != 42 {
		t.Errorf("got %+v; want ID=100, UserID=42", resp)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// ... TestCreateOrder_EmptyCart, and other tests ...
