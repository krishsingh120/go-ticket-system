package main

import (
	"log"
	"net/http"
	"os"

	"ticket-system/internal/database"
	"ticket-system/internal/handlers"
	"ticket-system/internal/middleware"
)

func main() {
	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "tickets.db" // default fallback
	}

	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully.")

	// Create a new ServeMux (router)
	router := http.NewServeMux()

	// Public routes
	router.HandleFunc("GET /health", handlers.HealthCheck)
	router.HandleFunc("POST /auth/register", handlers.Register)
	router.HandleFunc("POST /auth/login", handlers.Login)

	// Protected routes (require JWT)
	protected := http.NewServeMux()
	protected.HandleFunc("POST /tickets", handlers.CreateTicket)
	protected.HandleFunc("GET /tickets", handlers.ListTickets)
	protected.HandleFunc("GET /tickets/{id}", handlers.GetTicket)
	protected.HandleFunc("PATCH /tickets/{id}/status", handlers.UpdateTicketStatus)

	// Apply JWT middleware to protected routes
	// Note: We mount the protected multiplexer to the main router with a slight trick
	// since 1.22 paths matching in main router needs to route to the protected router.
	// Since standard lib mux matches prefixes if they end with /, we can just register them directly.
	// But it's easier to just wrap the handler functions directly in the main router.
	
	router.Handle("POST /tickets", middleware.JWTMiddleware(http.HandlerFunc(handlers.CreateTicket)))
	router.Handle("GET /tickets", middleware.JWTMiddleware(http.HandlerFunc(handlers.ListTickets)))
	router.Handle("GET /tickets/{id}", middleware.JWTMiddleware(http.HandlerFunc(handlers.GetTicket)))
	router.Handle("PATCH /tickets/{id}/status", middleware.JWTMiddleware(http.HandlerFunc(handlers.UpdateTicketStatus)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddr := ":" + port
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
