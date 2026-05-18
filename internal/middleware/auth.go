package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zonafirmann/next-go-pos-backend/internal/models"
)

// RequireAuth is a middleware that intercepts incoming HTTP requests
// and strictly validates the JSON Web Token (JWT) before allowing access.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract the Authorization header from the incoming request
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Unauthorized: Missing Authorization header"}`, http.StatusUnauthorized)
			return
		}

		// 2. Enforce the global standard "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": "Unauthorized: Invalid token format"}`, http.StatusUnauthorized)
			return
		}

		// Extract the actual JWT string
		tokenString := parts[1]

		// 3. Initialize the claims struct and fetch the secret key
		claims := &models.Claims{}
		jwtSecret := os.Getenv("JWT_SECRET")

		// 4. Parse the token, verify the cryptographic signature, and decode the claims
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		// Reject if the signature is invalid, forged, or the token is expired
		if err != nil || !token.Valid {
			http.Error(w, `{"error": "Unauthorized: Invalid, forged, or expired token"}`, http.StatusUnauthorized)
			return
		}

		// 5. Token is valid! Inject the verified username into the request Context
		// This allows downstream handlers to securely identify who is making the request
		ctx := context.WithValue(r.Context(), "username", claims.Username)

		// 6. Pass the baton to the next handler (allow entry)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
