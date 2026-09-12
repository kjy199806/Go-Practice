package main

import (
	"database/sql"
	"fmt"
)

// GetAllUsers retrieves all users from SQLite
func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`
		SELECT ID, FirstName, LastName, Email, 
			   Street, City, State, PostalCode, Country 
		FROM users
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
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
func GetUserByID(db *sql.DB, id int) (*User, error) {
	var u User
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
