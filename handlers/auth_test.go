package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// setupMockDB gives you a sqlx.DB backed by sqlmock
func setupMockDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	return sqlx.NewDb(db, "sqlmock"), mock
}

func TestSignup(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		mockSetup  func(sqlmock.Sqlmock)
		wantStatus int
	}{
		{
			name:       "invalid JSON",
			body:       `{"email":`,
			mockSetup:  func(m sqlmock.Sqlmock) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate email",
			body: `{"email":"dupe@example.com","password":"pass"}`,
			mockSetup: func(m sqlmock.Sqlmock) {
				// INSERT ... RETURNING returns a unique violation
				m.ExpectQuery(`INSERT INTO users`).
					WithArgs("dupe@example.com", sqlmock.AnyArg()).
					WillReturnError(&pq.Error{Code: "23505"})
			},
			wantStatus: http.StatusInternalServerError, // handler maps all errors to 500
		},
		{
			name: "successful signup",
			body: `{"email":"new@example.com","password":"pass"}`,
			mockSetup: func(m sqlmock.Sqlmock) {
				now := time.Now()
				m.ExpectQuery(`INSERT INTO users .*RETURNING id, email, created_at`).
					WithArgs("new@example.com", sqlmock.AnyArg()).
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "email", "created_at"}).
						AddRow(1, "new@example.com", now),
					)
			},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			tt.mockSetup(mock)

			ah := NewAuthHandler(db, "secret")
			req := httptest.NewRequest("POST", "/users/signup", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			ah.Signup(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("signup status = %d; want %d", w.Code, tt.wantStatus)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	pw := "mypassword"
	hash, _ := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

	tests := []struct {
		name       string
		body       string
		mockSetup  func(sqlmock.Sqlmock)
		wantStatus int
		wantToken  bool
	}{
		{
			name:       "invalid JSON",
			body:       `{"email":`,
			mockSetup:  func(m sqlmock.Sqlmock) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "wrong credentials",
			body: `{"email":"noone@example.com","password":"bad"}`,
			mockSetup: func(m sqlmock.Sqlmock) {
				m.ExpectQuery(`SELECT \* FROM users WHERE email=\$1`).
					WithArgs("noone@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "successful login",
			body: `{"email":"user@example.com","password":"` + pw + `"}`,
			mockSetup: func(m sqlmock.Sqlmock) {
				now := time.Now()
				m.ExpectQuery(`SELECT \* FROM users WHERE email=\$1`).
					WithArgs("user@example.com").
					WillReturnRows(sqlmock.
						NewRows([]string{"id", "email", "password_hash", "created_at"}).
						AddRow(42, "user@example.com", string(hash), now),
					)
			},
			wantStatus: http.StatusOK,
			wantToken:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			tt.mockSetup(mock)

			ah := NewAuthHandler(db, "secret")
			req := httptest.NewRequest("POST", "/users/login", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			ah.Login(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("login status = %d; want %d", w.Code, tt.wantStatus)
			}
			if tt.wantToken {
				var resp struct {
					Token string `json:"token"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to parse JSON: %v", err)
				}
				if resp.Token == "" {
					t.Errorf("expected non-empty token")
				}
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet expectations: %v", err)
			}
		})
	}
}
