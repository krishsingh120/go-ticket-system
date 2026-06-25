# Backend Intern Assignment: Ticket System in Golang

This is a backend service for a ticket system built in Golang. It provides REST APIs for user registration, login (JWT based), ticket creation, and ticket management.

## Architecture

The project uses a standard Go layout with modular packages:
- `cmd/api`: Contains the application entry point.
- `internal/database`: Handles SQLite connection and migrations.
- `internal/models`: Defines data structures (User, Ticket).
- `internal/repository`: Contains SQL queries and database interactions.
- `internal/service`: Enforces business logic (status flow rules).
- `internal/handlers`: Maps HTTP endpoints to service logic.
- `internal/middleware`: Validates JWT tokens and secures endpoints.

## Assumptions
- For simplicity, the JWT secret key is hardcoded in `internal/utils/utils.go`. In a production environment, it should be loaded from the environment variables.
- SQLite is used as the database as it provides a simple persistent store without external dependencies.
- A user can update the status from `in_progress` back to `open`, but neither can be updated from `closed`. The assignment requires `open -> in_progress -> closed`, and specifies `closed -> cannot move back to open or in_progress`.

## Deployment
- **Live URL:** `[Paste your Render/Railway deployment URL here]`
- **Health Check URL:** `[Paste your Render/Railway deployment URL here]/health`

## Local Run Instructions

To build and run the application locally using Docker, execute the following commands in the terminal:

1. **Build the Docker Image:**
   ```bash
   docker build -t ticket-system .
   ```

2. **Run the Docker Container:**
   ```bash
   docker run -p 8080:8080 ticket-system
   ```

3. **Verify the Application:**
   ```bash
   curl http://localhost:8080/health
   ```
   You should receive the response: `{"status": "ok"}`

## Endpoints

| Method | Endpoint | Purpose | Protected |
|---|---|---|---|
| GET | `/health` | Health check | No |
| POST | `/auth/register` | Register user | No |
| POST | `/auth/login` | Login and return JWT | No |
| POST | `/tickets` | Create ticket | Yes |
| GET | `/tickets` | List logged-in user tickets | Yes |
| GET | `/tickets/{id}` | Get own ticket by ID | Yes |
| PATCH | `/tickets/{id}/status` | Update own ticket status | Yes |

*Protected endpoints require an `Authorization: Bearer <token>` header.*
