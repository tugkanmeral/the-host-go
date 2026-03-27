# config (Configuration)

Loads settings from environment variables.

## Contents
Example layout:
```
config/
└── config.go        # Load all configuration
```

## Example config struct
```go
type Config struct {
    // Server
    Port     string
    
    // Database
    MongoURI string
    DBName   string
    
    // JWT
    JWTSecret string
    
    // Environment
    Env string // development, staging, production
}
```

## Usage
```go
func LoadConfig() *Config {
    // Read from environment variables
    // Apply defaults
    // Return configuration
}
```

## Environment variables
```
PORT=8000
MONGO_URI=mongodb://localhost:27017
MONGO_DB=your_db_name
JWT_SECRET=your_secret_key
APP_ENV=development
```

## Notes
- You can use a `.env` file (e.g. with godotenv)
- In production, rely on real environment variables
- Manage secrets via `.env` locally
