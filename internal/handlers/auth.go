package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/zonafirmann/next-go-pos-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler injects the database connection into the authentication HTTP handlers.
type AuthHandler struct {
	DB *pgx.Conn
}

// Register handles the creation of a new cashier account.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// 1. Hash the password using Bcrypt for secure storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
		return
	}

	// 2. Insert the new user record into the PostgreSQL database
	_, err = h.DB.Exec(context.Background(),
		"INSERT INTO users (username, password_hash, role) VALUES ($1, $2, 'cashier')",
		creds.Username, string(hashedPassword))

	if err != nil {
		http.Error(w, `{"error": "Username already exists or database transaction failed"}`, http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "User registered successfully. You may now log in."}`))
}

// Login authenticates a user and issues a JWT session token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	// 1. Retrieve the user record from the database
	var expectedPasswordHash, role string
	err := h.DB.QueryRow(context.Background(),
		"SELECT password_hash, role FROM users WHERE username=$1",
		creds.Username).Scan(&expectedPasswordHash, &role)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, `{"error": "User not found"}`, http.StatusUnauthorized)
			return
		}
		http.Error(w, `{"error": "Internal database error"}`, http.StatusInternalServerError)
		return
	}

	// 2. Verify the provided password against the stored Bcrypt hash
	if err = bcrypt.CompareHashAndPassword([]byte(expectedPasswordHash), []byte(creds.Password)); err != nil {
		http.Error(w, `{"error": "Invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	// 3. If credentials are valid, generate a JWT token (Valid for 24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Username: creds.Username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		http.Error(w, `{"error": "Failed to generate security token"}`, http.StatusInternalServerError)
		return
	}

	// 4. Send the JWT token back to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":   tokenString,
		"message": "Authentication successful",
	})
}
