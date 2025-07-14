package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/Heisenberg270/ecommerce-go/models"
)

// AuthHandler holds DB and JWT config
type AuthHandler struct {
	DB        *sqlx.DB
	JWTSecret string
}

// NewAuthHandler constructs an AuthHandler
func NewAuthHandler(db *sqlx.DB, secret string) *AuthHandler {
	return &AuthHandler{
		DB:        db,
		JWTSecret: secret,
	}
}

// Signup handles POST /users/signup
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var inp struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(inp.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}
	// insert user
	var user models.User
	insert := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, created_at`
	if err := h.DB.QueryRowx(insert, inp.Email, string(hash)).StructScan(&user); err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login handles POST /users/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var inp struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&inp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// fetch user
	var user models.User
	if err := h.DB.Get(&user, "SELECT * FROM users WHERE email=$1", inp.Email); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(inp.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	// create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})
	signed, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		http.Error(w, "failed to sign token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": signed})
}

// Logout or other auth endpoints can be added later
