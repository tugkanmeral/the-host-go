# config/environments (Environment files)

Contains `.env` files for different environments (development, staging, production).

## Contents
Example layout:
```
environments/
├── .env.development     # Local development
├── .env.staging         # Staging
└── .env.production      # Production
```

## Variables each file should define
```
# Server
PORT=8000

# Database
MONGO_URI=mongodb://localhost:27017
MONGO_DB=your_db_name

# JWT
JWT_SECRET=your_secret_key
JWT_EXPIRATION=3600

# Environment
APP_ENV=development
```

## Usage

### Development
```bash
export $(cat config/environments/.env.development | xargs)
go run ./cmd/main.go
```

### Production
```bash
export $(cat config/environments/.env.production | xargs)
go run ./cmd/main.go
```

## Notes
- Add `.env` files to `.gitignore` (they hold production-sensitive values)
- Use different values per environment
- Never commit secrets
