package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB initializes the database connection and runs migrations.
func InitDB(dbPath string) error {
	var err error
	// Use sqlite3 driver
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	// Verify connection
	if err = DB.Ping(); err != nil {
		return err
	}

	return runMigrations()
}

func runMigrations() error {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL
	);`

	ticketsTable := `
	CREATE TABLE IF NOT EXISTS tickets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	if _, err := DB.Exec(usersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	if _, err := DB.Exec(ticketsTable); err != nil {
		return fmt.Errorf("failed to create tickets table: %w", err)
	}

	return nil
}
