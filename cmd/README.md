# cmd

Application entry points: **HTTP API** and **interactive CLI**.

## Layout

```
cmd/
├── api/
│   └── main.go       # Fiber server: config, MongoDB, routes, listen
└── cli/
    └── main.go       # Bubble Tea TUI: config, MongoDB, cli.Run()
```

## api

- Loads environment and opens MongoDB.
- Registers auth and note routes (see `cmd/api/main.go` and `internal/handlers`).

```bash
go run ./cmd/api
```

## cli

- Same configuration and database wiring as the API.
- Starts the terminal UI (`internal/cli.Run`), which uses `internal/cli/page` for screens and `internal/cli/appsvc` for authenticated service calls.

```bash
go run ./cmd/cli
```

## Build binaries

```bash
go build -o the-host-api ./cmd/api
go build -o the-host-cli ./cmd/cli
```
