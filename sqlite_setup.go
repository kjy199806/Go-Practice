package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"practice/models"
	"time"

	_ "modernc.org/sqlite"
)

// InitDB initializes the SQLite connection, applies table schemas, and seeds initial data.
func InitDB(dbPath string) (*sql.DB, error) {

	// 1. Open the SQLite database file (creates it if it doesn't exist)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database file: %w", err)
	}

	// 2. Configure connection pooling for SQLite
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// 3. Verify connection works
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	// 4. Create the users table schema
	if err := createSchemas(db); err != nil {
		return nil, fmt.Errorf("schema creation failed: %w", err)
	}

	// 5. Seed initial users data if database is empty
	if err := seedUsersIfEmpty(db); err != nil {
		log.Printf("Warning: Seeding initial JSON data failed: %v", err)
	}

	return db, nil
}

// createSchemas handles creation of database tables
func createSchemas(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		FirstName TEXT NOT NULL,
		LastName TEXT NOT NULL,
		Email TEXT NOT NULL UNIQUE,
		Street TEXT,
		City TEXT,
		State TEXT,
		PostalCode TEXT,
		Country TEXT
	);`

	_, err := db.Exec(query)
	return err
}

// seedUsersIfEmpty checks if table has records; if not, loads from data/users.json
func seedUsersIfEmpty(db *sql.DB) error {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check record count: %w", err)
	}

	// Skip seeding if data is already present in the database
	if count > 0 {
		return nil
	}

	log.Println("Database is empty. Seeding initial records from data/users.json...")

	// Read JSON seed file
	data, err := os.ReadFile("data/users.json")
	if err != nil {
		return fmt.Errorf("unable to read seed JSON file: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return fmt.Errorf("unable to parse seed JSON file: %w", err)
	}

	// Use a SQL Transaction for bulk insertion performance
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO users (ID, FirstName, LastName, Email, Street, City, State, PostalCode, Country)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	for _, c := range users {
		_, err := stmt.Exec(
			c.ID, c.FirstName, c.LastName, c.Email,
			c.Address.Street, c.Address.City, c.Address.State, c.Address.PostalCode, c.Address.Country,
		)
		if err != nil {
			return fmt.Errorf("failed to insert user ID %d: %w", c.ID, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Println("Database successfully seeded!")
	return nil
}

func printData(db *sql.DB) {
	// --- QUICK DATABASE TEST ---
	// Fetch the total count of users
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to query user count: %v", err)
	}
	fmt.Printf("Successfully connected! Total users in DB: %d\n", count)

	// Fetch and print the all user in for loop
	for i := 0; i < count; i++ {
		var u models.User

		// Pass (i + 1) to match SQLite's 1-based IDs (1, 2, 3...)
		targetID := i + 1

		err = db.QueryRow(`
			SELECT ID, FirstName, LastName, Email, 
				   Street, City, State, PostalCode, Country 
			FROM users WHERE ID = ?
		`, targetID).Scan(
			&u.ID, &u.FirstName, &u.LastName, &u.Email,
			&u.Address.Street, &u.Address.City, &u.Address.State, &u.Address.PostalCode, &u.Address.Country,
		)
		if err != nil {
			log.Fatalf("Failed to query user ID %d: %v", targetID, err)
		}
		fmt.Printf("%d User in DB: ID=%d, Name=%s %s, Email=%s, City=%s\n",
			i, u.ID, u.FirstName, u.LastName, u.Email, u.Address.City)
	}

	log.Println("Database test passed!")
}
