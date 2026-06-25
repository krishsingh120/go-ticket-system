package models

import "time"

// User represents a registered user.
type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

// Ticket represents a support ticket.
type Ticket struct {
	ID          int       `json:"id"`
	UserID      int       `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // open, in_progress, closed
	CreatedAt   time.Time `json:"created_at"`
}

// RegisterRequest is the payload for user registration.
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginRequest is the payload for user login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateTicketRequest is the payload for creating a ticket.
type CreateTicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// UpdateTicketStatusRequest is the payload for updating ticket status.
type UpdateTicketStatusRequest struct {
	Status string `json:"status"`
}
