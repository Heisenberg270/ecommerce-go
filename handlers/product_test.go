package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/Heisenberg270/ecommerce-go/models"
)

// setupMock creates a sqlx.DB hooked to sqlmock
func setupMock(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockExpect func(sqlmock.Sqlmock)
		wantStatus int
	}{
		{
			name:       "invalid JSON",
			body:       `{"name":`,
			mockExpect: func(_ sqlmock.Sqlmock) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "valid request",
			body: `{"name":"Test","description":"Desc","price":1.23}`,
			mockExpect: func(m sqlmock.Sqlmock) {
				// Expect COUNT(*) query for sequence reset
				m.ExpectQuery(`SELECT COUNT\(\*\) FROM products`).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
				m.ExpectExec(`ALTER SEQUENCE products_id_seq RESTART WITH 1`).
					WillReturnResult(sqlmock.NewResult(0, 0))

				// Expect the INSERT ... RETURNING query
				m.ExpectQuery(`INSERT INTO products`).
					WithArgs("Test", "Desc", 1.23).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(1, time.Now(), time.Now()))
			},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMock(t)
			tt.mockExpect(mock)

			handler := NewProductHandler(db)
			req := httptest.NewRequest("POST", "/products", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Create(w, req)
			if w.Code != tt.wantStatus {
				t.Fatalf("status = %d; want %d", w.Code, tt.wantStatus)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestListProducts(t *testing.T) {
	// prepare some fake products
	now := time.Now()
	sample := []models.Product{
		{ID: 1, Name: "A", Description: "Alpha", Price: 9.99, CreatedAt: now, UpdatedAt: now},
	}

	db, mock := setupMock(t)
	rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "created_at", "updated_at"})
	for _, p := range sample {
		rows.AddRow(p.ID, p.Name, p.Description, p.Price, p.CreatedAt, p.UpdatedAt)
	}
	mock.ExpectQuery(`SELECT \* FROM products ORDER BY id`).WillReturnRows(rows)

	handler := NewProductHandler(db)
	req := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()
	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("List status = %d; want %d", w.Code, http.StatusOK)
	}
	var got []models.Product
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if len(got) != len(sample) {
		t.Fatalf("got %d products; want %d", len(got), len(sample))
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
