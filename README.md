# The Host Go

Go project with a **Fiber HTTP API** and a **terminal (TUI) CLI** that share the same MongoDB layer and **`internal/service`** business logic (auth, notes, and extensions).

## Project structure

```
the-host-go/
├── cmd/
│   ├── api/                      # HTTP API entry point (Fiber)
│   │   └── main.go
│   └── cli/                      # Interactive CLI entry point (Bubble Tea)
│       └── main.go
│
├── internal/
│   ├── cli/                      # TUI screens, session, app loop
│   ├── handlers/                 # HTTP handlers
│   ├── service/                  # Business logic (Auth, Note, …)
│   ├── models/                   # Entity + API DTOs
│   ├── database/                 # MongoDB connection & access
│   ├── middleware/               # HTTP middleware (e.g. JWT)
│   ├── auth/                     # JWT helpers, password hashing
│   └── config/                   # Configuration loading
│
└── config/
    └── environments/             # .env files (.development, .staging, .production)
```

## Folder reference

| Folder | Purpose |
|--------|---------|
| `cmd/api` | Starts the Fiber server; registers `/api/auth` and `/api/note` routes. |
| `cmd/cli` | Connects to MongoDB, loads config, runs the `internal/cli` TUI. |
| `internal/cli` | Terminal UI (session, screens). |
| `internal/handlers` | REST endpoint implementations. |
| `internal/service` | Shared domain logic for API and CLI. |
| `internal/models` | Database entities and API models. |
| `internal/database` | MongoDB connection and data access. |
| `internal/middleware` | Cross-cutting HTTP concerns (JWT, CORS, …). |
| `internal/auth` | JWT issuance/validation, password helpers. |
| `internal/config` | Environment / `.env` loading. |
| `config/environments` | Per-environment `.env` files. |

## Setup

### 1. Clone and dependencies

Module: `github.com/tugkanmeral/the-host-go`

```bash
git clone <repository-url> the-host-go
cd the-host-go
go mod download
```

### 2. Environment file

Create `config/environments/.env.development` and set values:

```bash
PORT=8000
MONGO_URI=mongodb://localhost:27017
MONGO_DB=your_database_name
JWT_SECRET=your_super_secret_key_here
JWT_EXPIRATION=3600
APP_ENV=development
```

### 3. MongoDB

```bash
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

### 4. Run

**HTTP API** (listen address comes from config; default example `:8000`):

```bash
go run ./cmd/api
```

**CLI** (same `.env` and MongoDB as the API):

```bash
go run ./cmd/cli
```

## Development flow

1. **Model** (`internal/models/`) — data structures
2. **Database** (`internal/database/`) — CRUD / queries
3. **Service** (`internal/service/`) — shared logic for API and CLI
4. **Handlers** (`internal/handlers/`) — HTTP endpoints
5. **CLI** (`internal/cli/`) — terminal flows and screens
6. **Routes** — Fiber groups in `cmd/api/main.go`
7. **Verify** — via API client or the CLI

## API overview

Route groups (summary):

- `POST /api/auth/register`, `POST /api/auth/login`
- `POST /api/note/`, `GET /api/note/`, `GET /api/note/:id`, `PUT /api/note/:id`, `DELETE /api/note/:id`

## API response format

```json
{
    "success": true,
    "message": "Operation successful",
    "data": { }
}
```

## Environment management

### Development

```bash
source config/environments/.env.development
go run ./cmd/api
# or
go run ./cmd/cli
```

### Production (API)

```bash
source config/environments/.env.production
go build -o the-host-api ./cmd/api
./the-host-api
```

### Production (CLI)

```bash
go build -o the-host-cli ./cmd/cli
./the-host-cli
```

## Important notes

- Keep secrets in `.env` files only; do not commit them.
- Ensure `.env` files are listed in `.gitignore`.
- Use a MongoDB connection pool in production workloads.
- Keep JWT expiration reasonable; hash passwords with bcrypt.
- Handle errors consistently and add logging where it helps operations.

## Libraries

- **Fiber** — HTTP framework
- **MongoDB Driver v2** — database
- **JWT** — authentication
- **golang.org/x/crypto** — password hashing (bcrypt)
- **godotenv** — environment variables
- **Bubble Tea / Lipgloss / Bubbles** — TUI CLI

## Help

See package-level `README.md` files and inline comments in the code.
