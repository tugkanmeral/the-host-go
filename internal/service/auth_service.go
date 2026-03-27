package service

import (
	"context"
	"errors"
	"time"

	"github.com/tugkanmeral/the-host-go/internal/auth"
	"github.com/tugkanmeral/the-host-go/internal/database"
	"github.com/tugkanmeral/the-host-go/internal/models/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthService struct {
	db *mongo.Database
}

func NewAuthService(db *mongo.Database) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user entity.User
	err := s.db.Collection(database.UserCollectionName).FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", ErrInvalidCredentials
		}
		return "", ErrInternal
	}

	if !auth.CheckPassword(user.Password, password) {
		return "", ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(user.ID.Hex(), user.Username)
	if err != nil {
		return "", ErrTokenGeneration
	}

	return token, nil
}

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var existing entity.User
	err := s.db.Collection(database.UserCollectionName).FindOne(ctx, bson.M{"username": username}).Decode(&existing)
	if err == nil {
		return ErrUsernameTaken
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return ErrInternal
	}

	hashedPass, err := auth.HashPassword(password)
	if err != nil {
		return ErrInternal
	}

	newUser := entity.User{
		Username: username,
		Password: hashedPass,
	}

	_, err = s.db.Collection(database.UserCollectionName).InsertOne(ctx, newUser)
	if err != nil {
		return ErrRegisterFailed
	}

	return nil
}
