package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tugkanmeral/the-host-go/internal/auth"
	"github.com/tugkanmeral/the-host-go/internal/database"
	model "github.com/tugkanmeral/the-host-go/internal/models/api"
	"github.com/tugkanmeral/the-host-go/internal/models/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Geçersiz istek gövdesi",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Email ve şifre zorunludur",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := database.GetDB()
	var user entity.User

	err := db.Collection("Users").FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{
				Message: "Kullanıcı adı veya şifre hatalı",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Message: "Kullanıcı aranırken bir hata oluştu",
		})
	}

	if !auth.CheckPassword(user.Password, req.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{
			Message: "Kullanıcı adı veya şifre hatalı",
		})
	}

	token, err := auth.GenerateToken(user.ID.Hex(), user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Message: "Token oluşturulurken bir hata oluştu",
		})
	}

	return c.JSON(model.LoginResponse{
		Token: token,
	})
}

func Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Invalid body",
		})
	}

	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Username and password fields are required",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db := database.GetDB()
	var existingUser entity.User
	err := db.Collection(database.UserCollectionName).FindOne(ctx, bson.M{"username": req.Username}).Decode(&existingUser)

	if existingUser.ID != bson.NilObjectID {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "Username conflict",
		})
	}

	hashedPass, _ := auth.HashPassword(req.Password)
	var newUser entity.User
	newUser.Username = req.Username
	newUser.Password = hashedPass

	_, err = db.Collection(database.UserCollectionName).InsertOne(ctx, newUser)
	if err != nil {
		fmt.Println("Failed to insert user:", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Message: "User cannot be registered!",
		})
	}

	return nil
}
