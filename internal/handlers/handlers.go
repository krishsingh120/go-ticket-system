package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"ticket-system/internal/middleware"
	"ticket-system/internal/models"
	"ticket-system/internal/repository"
	"ticket-system/internal/service"
	"ticket-system/internal/utils"
)

// HealthCheck responds with a status ok.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

// Register handles user registration.
func Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := service.RegisterUser(req)
	if err != nil {
		if errors.Is(err, repository.ErrConflict) {
			utils.JSONError(w, http.StatusConflict, "Username already exists")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

// Login handles user login and returns a JWT.
func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, err := service.LoginUser(req)
	if err != nil {
		utils.JSONError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"token": token})
}

// CreateTicket creates a new ticket for the logged-in user.
func CreateTicket(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	var req models.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	ticket, err := service.CreateTicket(userID, req)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to create ticket")
		return
	}

	utils.JSONResponse(w, http.StatusCreated, ticket)
}

// ListTickets returns all tickets for the logged-in user.
func ListTickets(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	tickets, err := service.GetUserTickets(userID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, "Failed to retrieve tickets")
		return
	}

	utils.JSONResponse(w, http.StatusOK, tickets)
}

// GetTicket returns a specific ticket by ID for the logged-in user.
func GetTicket(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	
	// Get ID from URL path, standard library 1.22 support (e.g., /tickets/{id})
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	ticket, err := service.GetUserTicketByID(id, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Ticket not found")
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, "Failed to retrieve ticket")
		return
	}

	utils.JSONResponse(w, http.StatusOK, ticket)
}

// UpdateTicketStatus updates the status of a specific ticket.
func UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid ticket ID")
		return
	}

	var req models.UpdateTicketStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = service.UpdateTicketStatus(id, userID, req.Status)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			utils.JSONError(w, http.StatusNotFound, "Ticket not found")
			return
		}
		if errors.Is(err, service.ErrInvalidStatusTransition) || errors.Is(err, service.ErrInvalidStatus) {
			utils.JSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.JSONError(w, http.StatusInternalServerError, "Failed to update ticket")
		return
	}

	utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "Ticket status updated"})
}
