# utils (Helpers)

Shared helper functions used across the project.

## Contents
Example layout:
```
utils/
├── response.go      # Standard API response shape
├── validators.go    # Validation helpers
└── helpers.go       # Other utilities
```

## Example functions

### Response formatter
```go
type ApiResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

func SuccessResponse(data interface{}) ApiResponse {
    return ApiResponse{Success: true, Data: data}
}

func ErrorResponse(msg string) ApiResponse {
    return ApiResponse{Success: false, Message: msg}
}
```

### Validators
```go
func IsValidEmail(email string) bool {
    // Validate email
}

func IsValidPassword(password string) bool {
    // Enforce password rules
}
```

## Notes
- Put reusable code here
- Keep functions generic so other packages can import them
