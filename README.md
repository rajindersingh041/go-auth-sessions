# go-auth-sessions

## Project Structure & File Responsibilities

This project is organized into multiple Go files, each with a clear responsibility to improve maintainability, readability, and testability:

- **main.go**: Application entry point. Starts the server and initializes dependencies.
- **server.go**: Defines the `Server` struct, sets up HTTP routes, and applies middleware.
- **handlers.go**: Contains HTTP handler functions for each endpoint (e.g., health check, register, login, protected resource).
- **middleware.go**: Contains middleware functions for logging, recovery, and authentication.
- **response.go**: Provides helper functions for sending JSON responses and errors.
- **db.go**: Handles database initialization and connection setup (ClickHouse in this case).
- **models.go**: Defines data models, such as the `User` struct.
- **repository.go**: Implements the `UserRepository` struct and all database operations related to users (CRUD, existence checks, etc.).
- **jwt_simple.go**: (If present) Handles JWT creation and validation logic.

### Why Modularize?

Splitting the codebase into focused files helps:
- Make each file easier to understand and maintain.
- Encourage separation of concerns (business logic, data access, HTTP, etc.).
- Simplify testing and future extension.
- Reduce merge conflicts in team environments.

### How to Use Each File

- **main.go**: Run this file to start the application. It wires up the server and database.
- **server.go**: Used by `main.go` to create and configure the HTTP server.
- **handlers.go**: Called by the server to process incoming HTTP requests.
- **middleware.go**: Used by the server to wrap handlers with cross-cutting concerns (logging, error recovery, authentication).
- **response.go**: Used by handlers and middleware to send consistent JSON responses.
- **db.go**: Called by `main.go` to initialize the database connection.
- **models.go**: Imported wherever user data structures are needed.
- **repository.go**: Used by handlers to interact with the database via the repository pattern.
- **jwt_simple.go**: Used by authentication handlers and middleware for token management.