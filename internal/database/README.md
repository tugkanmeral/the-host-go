# database (MongoDB connection and access)

MongoDB connection, collections, and query helpers live here.

## Contents
Example layout:
```
database/
├── mongodb.go       # Connect and initialize MongoDB
├── user_repo.go     # User collection queries
└── post_repo.go     # Post collection queries
```

## Responsibilities
1. **MongoDB connection**
   - `ConnectDB()` — connect to the database
   - `DisconnectDB()` — close the connection
   - `GetDB()` — return the database handle

2. **Repository pattern** (CRUD)
   - `CreateUser()` — create a user
   - `GetUser()` — fetch a user
   - `UpdateUser()` — update a user
   - `DeleteUser()` — delete a user

## Example connection
```go
// mongodb.go
var mongoClient *mongo.Client

func ConnectDB() error {
    // Connect to MongoDB
    // Read connection string from config
}

func GetDB() *mongo.Database {
    return mongoClient.Database(dbName)
}
```

## Notes
- Read the connection string from environment variables
- Use a connection pool
- Handle errors explicitly
