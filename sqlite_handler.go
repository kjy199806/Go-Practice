package main

import (
	"database/sql"
	"fmt"
	"practice/models"
)

// GetAllUsers retrieves all users from SQLite
func GetAllUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query(`
		SELECT ID, FirstName, LastName, Email, 
			   Street, City, State, PostalCode, Country 
		FROM users
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID, &u.FirstName, &u.LastName, &u.Email,
			&u.Address.Street, &u.Address.City, &u.Address.State, &u.Address.PostalCode, &u.Address.Country,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

// GetUserByID retrieves a single user by ID from SQLite
func GetUserByID(db *sql.DB, id int) (*models.User, error) {
	var u models.User
	err := db.QueryRow(`
		SELECT ID, FirstName, LastName, Email, Street, City, State, PostalCode, Country 
		FROM users WHERE ID = ?
	`, id).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Email,
		&u.Address.Street, &u.Address.City, &u.Address.State, &u.Address.PostalCode, &u.Address.Country,
	)

	if err != nil {
		return nil, err // Returns sql.ErrNoRows if not found
	}

	return &u, nil
}

// GetUserByEmail retrieves a single user by matched email from SQLite
func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {

	// Fetch user by email from database
	var u models.User
	err := db.QueryRow(`
		SELECT ID, FirstName, LastName, Email, 
				Street, City, State, PostalCode, Country 
		FROM users WHERE Email = ?
	`, email).Scan(
		&u.ID, &u.FirstName, &u.LastName, &u.Email,
		&u.Address.Street, &u.Address.City, &u.Address.State, &u.Address.PostalCode, &u.Address.Country,
	)

	if err != nil {
		return nil, err // Returns sql.ErrNoRows if not found
	}

	return &u, nil
}
