package models

import "github.com/golang-jwt/jwt/v5"

// Credentials represents the expected JSON payload for user registration and login.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims represents the JWT payload carrying user authorization data.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
