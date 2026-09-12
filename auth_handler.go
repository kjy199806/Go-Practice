package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

// POST /login
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
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

		// Return data
		respondWithJSON(w, http.StatusOK, LoginResponse{
			Token: "xsrf-token",
			User:  *u,
		})
	}
}
