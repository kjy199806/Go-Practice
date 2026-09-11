package main

import "log"

func main() {
	// 1. Initialize SQLite database (creates table & seeds users.json if empty)
	db, err := InitDB("users.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close() // Ensures the database connection closes when main exits

	printData(db)
}
