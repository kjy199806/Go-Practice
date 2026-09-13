package main

import (
	"log/slog"
	"net/http"
	"os"
	"practice/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Initialize SQLite database (creates table & seeds users.json if empty)
	db, err := InitDB("users.db")
	if err != nil {
		slog.Error("failed to initialize database", "err", err)
		os.Exit(1)
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
	slog.Info("server is running", "address", "http://localhost:"+cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		slog.Error("server failed to start", "err", err)
		os.Exit(1)
	}
}
