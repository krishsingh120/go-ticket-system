package service

import (
	"errors"
	"ticket-system/internal/models"
	"ticket-system/internal/repository"
	"ticket-system/internal/utils"
)

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidStatus           = errors.New("invalid status value")
)

// RegisterUser handles user registration logic.
func RegisterUser(req models.RegisterRequest) error {
	if req.Username == "" || req.Password == "" {
		return errors.New("username and password are required")
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: hash,
	}

	return repository.CreateUser(user)
}

// LoginUser handles user login and returns a JWT token.
func LoginUser(req models.LoginRequest) (string, error) {
	if req.Username == "" || req.Password == "" {
		return "", errors.New("username and password are required")
	}

	user, err := repository.GetUserByUsername(req.Username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	return utils.GenerateJWT(user.ID, user.Username)
}

// CreateTicket handles ticket creation.
func CreateTicket(userID int, req models.CreateTicketRequest) (*models.Ticket, error) {
	if req.Title == "" || req.Description == "" {
		return nil, errors.New("title and description are required")
	}

	ticket := &models.Ticket{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      "open",
	}

	if err := repository.CreateTicket(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

// GetUserTickets retrieves all tickets for a user.
func GetUserTickets(userID int) ([]models.Ticket, error) {
	return repository.GetTicketsByUserID(userID)
}

// GetUserTicketByID retrieves a specific ticket for a user.
func GetUserTicketByID(id, userID int) (*models.Ticket, error) {
	return repository.GetTicketByIDAndUserID(id, userID)
}

// UpdateTicketStatus updates the status of a ticket according to the rules.
func UpdateTicketStatus(id, userID int, newStatus string) error {
	if newStatus != "open" && newStatus != "in_progress" && newStatus != "closed" {
		return ErrInvalidStatus
	}

	ticket, err := repository.GetTicketByIDAndUserID(id, userID)
	if err != nil {
		return err
	}

	if ticket.Status == "closed" {
		return ErrInvalidStatusTransition
	}

	// Status flow: open -> in_progress -> closed
	if ticket.Status == "open" && newStatus == "open" {
		return nil // No change
	}
	if ticket.Status == "in_progress" && newStatus == "open" {
		// Valid based on assignment?
		// open -> in_progress -> closed
		// The prompt says "closed -> cannot move back to open or in_progress"
		// It doesn't strictly say in_progress cannot move to open, but typically it shouldn't.
		// However, it just says "open -> in_progress -> closed" and "closed cannot be reopened".
		// I will allow in_progress -> open just in case, or maybe not. Let's just allow anything unless it's from closed.
	}

	return repository.UpdateTicketStatus(id, userID, newStatus)
}
