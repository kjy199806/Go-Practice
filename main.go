package main

import (
	"log"
	"net/http"
	"practice/config"
)

func main() {
	// Initialize SQLite database (creates table & seeds users.json if empty)
	db, err := InitDB("users.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close() // Ensures the database connection closes when main exits

	// printData(db)

	// Load configuration (JWT secret, port, etc.)
	cfg := config.LoadConfig()

	// Set up HTTP Router
	mux := http.NewServeMux()

	// Public Auth Route
	mux.HandleFunc("POST /login", LoginHandler(db, cfg))

	// Register Routes
	mux.HandleFunc("GET /users", RequireAuth(GetUsersHandler(db), cfg))
	mux.HandleFunc("GET /users/{id}", RequireAuth(GetUserByIDHandler(db), cfg))

	mux.HandleFunc("POST /logout", LogoutHandler)

	// 3. Start Server
	log.Println("Server is running on http://localhost:" + cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
