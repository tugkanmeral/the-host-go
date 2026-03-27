# models (Data models)

Defines data structures (structs) for the database and API.

## Contents
Structs that map to MongoDB collections.

Example layout:
```
models/
├── user.go          # User model
├── post.go          # Post model
└── response.go      # API response models
```

## Model example
```go
type User struct {
    ID       primitive.ObjectID `bson:"_id,omitempty"`
    Email    string             `bson:"email"`
    Password string             `bson:"password"`
    Name     string             `bson:"name"`
}
```

## Notes
- `bson` tags must match MongoDB field names
- `json` tags shape API responses
- `ObjectID` is MongoDB’s primary key type
