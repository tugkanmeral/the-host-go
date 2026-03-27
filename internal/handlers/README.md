# handlers (Request handlers)

Handler functions for your API endpoints live here.

## Contents
Each handler processes an HTTP request and returns a response.

Example layout:
```
handlers/
├── user_handler.go          # User endpoints
├── auth_handler.go          # Login / logout endpoints
└── health_handler.go        # Health check
```

## Handler example
```go
func GetUser(c *fiber.Ctx) error {
    // Handle request
    // Load from database
    // Return response
    return c.JSON(user)
}
```

## Notes
- Handlers take `*fiber.Ctx`
- Prefer JSON responses
- They run after applicable middleware
