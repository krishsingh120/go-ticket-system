package repository

import (
	"database/sql"
	"errors"
	"ticket-system/internal/database"
	"ticket-system/internal/models"
)

var ErrNotFound = errors.New("record not found")
var ErrConflict = errors.New("record already exists")

// CreateUser inserts a new user into the database.
func CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, password_hash) VALUES (?, ?)`
	result, err := database.DB.Exec(query, user.Username, user.PasswordHash)
	if err != nil {
		// Basic check for unique constraint failure
		if err.Error() == "UNIQUE constraint failed: users.username" {
			return ErrConflict
		}
		return err
	}
	id, _ := result.LastInsertId()
	user.ID = int(id)
	return nil
}

// GetUserByUsername retrieves a user by their username.
func GetUserByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, password_hash FROM users WHERE username = ?`
	row := database.DB.QueryRow(query, username)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return user, nil
}

// CreateTicket inserts a new ticket into the database.
func CreateTicket(ticket *models.Ticket) error {
	query := `INSERT INTO tickets (user_id, title, description, status) VALUES (?, ?, ?, ?)`
	result, err := database.DB.Exec(query, ticket.UserID, ticket.Title, ticket.Description, ticket.Status)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	ticket.ID = int(id)

	// Fetch created_at
	row := database.DB.QueryRow(`SELECT created_at FROM tickets WHERE id = ?`, ticket.ID)
	_ = row.Scan(&ticket.CreatedAt)

	return nil
}

// GetTicketsByUserID retrieves all tickets for a specific user.
func GetTicketsByUserID(userID int) ([]models.Ticket, error) {
	query := `SELECT id, user_id, title, description, status, created_at FROM tickets WHERE user_id = ?`
	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []models.Ticket
	for rows.Next() {
		var t models.Ticket
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt); err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Return empty slice instead of nil for JSON serialization
	if tickets == nil {
		tickets = []models.Ticket{}
	}

	return tickets, nil
}

// GetTicketByIDAndUserID retrieves a specific ticket, ensuring it belongs to the user.
func GetTicketByIDAndUserID(id, userID int) (*models.Ticket, error) {
	query := `SELECT id, user_id, title, description, status, created_at FROM tickets WHERE id = ? AND user_id = ?`
	row := database.DB.QueryRow(query, id, userID)

	t := &models.Ticket{}
	err := row.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return t, nil
}

// UpdateTicketStatus updates the status of a specific ticket, ensuring it belongs to the user.
func UpdateTicketStatus(id, userID int, status string) error {
	query := `UPDATE tickets SET status = ? WHERE id = ? AND user_id = ?`
	result, err := database.DB.Exec(query, status, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
