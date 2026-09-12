package main

import (
	"log"
	"net/http"
)

func main() {
	// 1. Initialize SQLite database (creates table & seeds users.json if empty)
	db, err := InitDB("users.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close() // Ensures the database connection closes when main exits

	// printData(db)

	// 2. Set up HTTP Router
	mux := http.NewServeMux()

	// Public Auth Route
	mux.HandleFunc("POST /login", LoginHandler(db))

	// Register Routes
	mux.HandleFunc("GET /users", GetUsersHandler(db))
	mux.HandleFunc("GET /users/{id}", GetUserByIDHandler(db))

	// 3. Start Server
	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
