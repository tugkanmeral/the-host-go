# middleware (Cross-cutting handlers)

Middleware that runs on incoming requests lives here.

## Contents
Example layout:
```
middleware/
├── auth_middleware.go       # JWT verification
├── logger_middleware.go     # Request logging
└── cors_middleware.go       # CORS settings
```

## Common middleware

### JWT middleware
```go
func JWTMiddleware(c *fiber.Ctx) error {
    // Check Authorization header
    // Verify JWT
    // Attach user to context
}
```

### Logger middleware
```go
func LoggerMiddleware(c *fiber.Ctx) error {
    // Log request (method, URL, IP)
    // Call next handler
}
```

## Middleware order
1. CORS
2. Logger
3. JWT (for protected routes)
4. Route handler

## Notes
- Call `c.Next()` to continue the chain
- You can store data on the context: `c.Locals("key", value)`
- Protect routes with JWT where required
