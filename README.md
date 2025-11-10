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


### Code Comments & Naming

- All exported functions, structs, and interfaces are commented using GoDoc style for clarity.
- Variable and function names are descriptive and use `camelCase` for variables and `PascalCase` for types and exported functions.
- Example: `userRepository` for a variable, `UserRepository` for an interface, `passwordHasher` for a dependency, etc.

### Application Workflow (Architecture Diagram)

Below is a high-level workflow of the application. You can view and edit the diagram in [draw.io](https://app.diagrams.net/):

![Application Workflow Diagram](./docs/go-auth-sessions-workflow.png)

**Diagram Description:**

1. **Client** sends HTTP requests (register, login, protected resource) to the server.
2. **Server** routes requests to the appropriate handler.
3. **Handlers** use injected dependencies:
	- `UserRepository` for database operations
	- `PasswordHasher` for password hashing/verification
	- `JWTManager` for token generation/validation
4. **Database** stores user data (ClickHouse).
5. **JWT** is used for stateless authentication.
6. **Middleware** handles logging, error recovery, and authentication.

> You can edit or export the diagram using draw.io. Place the PNG or XML file in the `docs/` folder.

---

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