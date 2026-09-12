package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

// GET /users
func GetUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := GetAllUsers(db)
		if err != nil {
			log.Printf("Error getting users: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve users")
			return
		}

		respondWithJSON(w, http.StatusOK, users)
	}
}

// GET /users/{id}
func GetUserByIDHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			respondWithError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		user, err := GetUserByID(db, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusNotFound, "User not found")
				return
			}
			log.Printf("Error getting user %d: %v", id, err)
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user")
			return
		}

		respondWithJSON(w, http.StatusOK, user)
	}
}
