package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

// ContextUserID is the key we use to store the user ID in request contexts.
const ContextUserID = contextKey("userID")

// AuthMiddleware parses a Bearer JWT and stores the user ID in the context.
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "missing token", http.StatusUnauthorized)
				return
			}
			tokenString := strings.TrimPrefix(auth, "Bearer ")
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			claims := token.Claims.(jwt.MapClaims)
			sub, ok := claims["sub"].(float64)
			if !ok {
				http.Error(w, "invalid token claims", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserID, int(sub))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
