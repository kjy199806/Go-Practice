package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"practice/config"
	"practice/models"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

func RequireAuth(next http.HandlerFunc, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		// 1. Try reading token from cookie
		cookie, err := r.Cookie("auth_token")
		if err == nil {
			tokenString = cookie.Value
		} else {
			// 2. Fallback to Authorization header if cookie isn't present
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// Reject request if no token was supplied
		if tokenString == "" {
			respondWithError(w, http.StatusUnauthorized, "Missing auth token")
			return
		}

		// 3. Cryptographically parse and validate the token
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return cfg.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		// 4. Token is valid — inject UserID into request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

		// 5. Pass request forward to the protected handler
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
