# The Host Go - Fiber Web API

Go Web API project with Fiber, MongoDB, JWT authentication, and multiple environments.

## Project Structure

```
the-host-go/
├── cmd/                          # Application entry point
│   └── main.go                   # Server startup
│
├── internal/                     # Application logic (not importable by external packages)
│   ├── handlers/                 # HTTP request handlers
│   ├── models/                   # Data models (structs)
│   ├── database/                 # MongoDB connection and repository
│   ├── middleware/               # Middleware functions
│   ├── auth/                     # JWT and authentication
│   ├── config/                   # Configuration management
│   └── utils/                    # Helper functions
│
└── config/
    └── environments/             # .env files (.development, .staging, .production)
```

## Folder Reference

| Folder | Purpose |
|--------|---------|
| `cmd` | Where the app starts. Server bootstrap and route registration. |
| `handlers` | Functions that implement API endpoints. |
| `models` | Data structures for the database and API. |
| `database` | MongoDB connection, collections, and queries. |
| `middleware` | Cross-cutting concerns such as JWT, logging, and CORS. |
| `auth` | JWT issuance, validation, and password handling. |
| `config` | Load configuration values. |
| `utils` | Reusable helper functions. |
| `environments` | `.env` files for different environments. |

## Setup

### 1. Initialize the Go module
```bash
go mod init github.com/yourusername/the-host-go
```

### 2. Install dependencies
```bash
# Fiber Framework
go get github.com/gofiber/fiber/v2

# MongoDB Driver
go get go.mongodb.org/mongo-driver

# JWT
go get github.com/golang-jwt/jwt/v5

# Password hashing
go get golang.org/x/crypto

# Environment variables
go get github.com/joho/godotenv
```

### 3. Create the environment file
Create `config/environments/.env.development` and fill in the values:

```bash
PORT=8000
MONGO_URI=mongodb://localhost:27017
MONGO_DB=your_database_name
JWT_SECRET=your_super_secret_key_here
JWT_EXPIRATION=3600
APP_ENV=development
```

### 4. Start MongoDB
```bash
# With Docker
docker run -d -p 27017:27017 --name mongodb mongo:latest

# Or start a locally installed MongoDB instance
```

### 5. Run the application
```bash
go run ./cmd/main.go
```

## Development Flow

1. **Define a model** (`internal/models/`) — create the data structure
2. **Database repository** (`internal/database/`) — implement CRUD
3. **Write a handler** (`internal/handlers/`) — endpoint logic
4. **Add middleware** (`internal/middleware/`) — wire middleware as needed
5. **Register routes** (`cmd/main.go`) — connect handlers to routes
6. **Test** — exercise your API

## Example: user creation flow

1. Add the `User` struct in `internal/models/user.go`
2. Implement `CreateUser` in `internal/database/user_repo.go`
3. Add `CreateUserHandler` in `internal/handlers/user_handler.go`
4. Enforce JWT in `internal/middleware/auth_middleware.go` where required
5. Register `POST /users` in `cmd/main.go`

## API response format

All responses follow a standard shape:

```json
{
    "success": true,
    "message": "Operation successful",
    "data": { ... }
}
```

## Environment management

### Development
```bash
source config/environments/.env.development
go run ./cmd/main.go
```

### Production
```bash
source config/environments/.env.production
go build -o app ./cmd/main.go
./app
```

## Important notes

- Keep secrets in `.env` files
- Add `.env` files to `.gitignore`
- Use a MongoDB connection pool
- Keep JWT expiration reasonably short
- Hash passwords with bcrypt
- Handle errors consistently
- Add logging (request details as appropriate)

## Libraries

- **Fiber** — Web framework
- **MongoDB Driver** — Database
- **JWT** — Authentication
- **bcrypt** — Password hashing
- **godotenv** — Environment variables

## Help

For issues, read the README files in each package or check inline comments in the code.

Happy coding! 🚀
