# auth (JWT authentication)

JWT issuance, validation, and authentication logic live here.

## Contents
Example layout:
```
auth/
├── jwt.go           # Create and verify JWTs
├── password.go      # Password hashing and verification
└── claims.go        # JWT claim definitions
```

## Responsibilities

### JWT token management
```go
func GenerateToken(userID string) (string, error) {
    // Issue token
    // Use secret key
    // Set expiration
}

func VerifyToken(tokenString string) (*jwt.Claims, error) {
    // Verify token
    // Extract claims
    // Check validity
}
```

### Password operations
```go
func HashPassword(password string) string {
    // Hash password (bcrypt)
}

func CheckPassword(hashed, plaintext string) bool {
    // Compare hash and plaintext
}
```

### JWT claims
```go
type Claims struct {
    UserID string
    Email  string
    jwt.StandardClaims
}
```

## Flow
1. User logs in → verify password
2. On success → issue JWT
3. Return token to the client
4. Each request is validated in JWT middleware

## Notes
- Read the secret from an environment variable
- Keep token lifetime short (15–30 minutes)
- Never store passwords in plain text
