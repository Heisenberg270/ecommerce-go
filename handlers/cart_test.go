package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/Heisenberg270/ecommerce-go/models"
)

func setupCartMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func TestCreateCart(t *testing.T) {
	db, mock := setupCartMock(t)
	// Expect INSERT ... RETURNING
	now := time.Now()
	mock.ExpectQuery(`INSERT INTO carts .*RETURNING id, user_id, created_at`).
		WithArgs(99).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "created_at"}).
			AddRow(123, 99, now),
		)

	ch := NewCartHandler(db)
	req := httptest.NewRequest("POST", "/carts", nil)
	// inject authenticated user ID
	req = req.WithContext(context.WithValue(req.Context(), ContextUserID, 99))
	w := httptest.NewRecorder()

	ch.CreateCart(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("CreateCart status = %d; want %d", w.Code, http.StatusCreated)
	}
	var cart models.Cart
	if err := json.Unmarshal(w.Body.Bytes(), &cart); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if cart.ID != 123 || cart.UserID != 99 {
		t.Errorf("got cart %+v; want ID=123, UserID=99", cart)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestAddItem_BadJSON(t *testing.T) {
	db, _ := setupCartMock(t)
	ch := NewCartHandler(db)
	req := httptest.NewRequest("POST", "/carts/1/items", bytes.NewBufferString(`{"product_id":`))
	w := httptest.NewRecorder()

	ch.AddItem(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("AddItem bad JSON status = %d; want %d", w.Code, http.StatusBadRequest)
	}
}

func TestAddItem_Success(t *testing.T) {
	db, mock := setupCartMock(t)
	mock.ExpectExec(`INSERT INTO cart_items`).
		WithArgs(5, 10, 3).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ch := NewCartHandler(db)
	req := httptest.NewRequest("POST", "/carts/5/items",
		bytes.NewBufferString(`{"product_id":10,"quantity":3}`))
	// inject cartID=5 into chi route context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cartID", "5")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	ch.AddItem(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("AddItem status = %d; want %d", w.Code, http.StatusNoContent)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestGetCart_NotFound(t *testing.T) {
	db, mock := setupCartMock(t)
	mock.ExpectQuery(`SELECT id, user_id, created_at FROM carts`).
		WithArgs(7).
		WillReturnError(sql.ErrNoRows)

	ch := NewCartHandler(db)
	req := httptest.NewRequest("GET", "/carts/7", nil)
	// inject cartID=7 into chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cartID", "7")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	ch.GetCart(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("GetCart not found status = %d; want %d", w.Code, http.StatusNotFound)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestRemoveItem(t *testing.T) {
	db, mock := setupCartMock(t)
	mock.ExpectExec(`DELETE FROM cart_items`).
		WithArgs(8, 20).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ch := NewCartHandler(db)
	req := httptest.NewRequest("DELETE", "/carts/8/items/20", nil)
	// inject cartID=8 & productID=20 into chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cartID", "8")
	rctx.URLParams.Add("productID", "20")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	w := httptest.NewRecorder()
	ch.RemoveItem(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("RemoveItem status = %d; want %d", w.Code, http.StatusNoContent)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
