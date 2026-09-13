package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"practice/config"
	"practice/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// POST /login
func LoginHandler(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Fetch user by email from database (sqlite handler)
		u, err := GetUserByEmail(db, req.Email)

		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
				return
			}
			respondWithError(w, http.StatusInternalServerError, "Internal Server error")
			return
		}

		// Generate secure JWT token
		token, err := GenerateJWT(u, cfg.JWTSecret)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Token creation failed")
			return
		}

		// Set JWT as an HttpOnly, Secure cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "auth_token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		// Set XSRF token in cookie
		respondWithJSON(w, http.StatusOK, u)
	}
}

// GenerateJWT creates a signed token string containing the user's ID and Email
func GenerateJWT(user *models.User, jwtSecret []byte) (string, error) {
	// Set token to expire in 1 hour
	expirationTime := time.Now().Add(60 * time.Minute)

	claims := &models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Declare token with HMAC-SHA256 signing algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create complete signed JWT string using secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully logged out"})
}
