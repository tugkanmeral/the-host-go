package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	model "github.com/tugkanmeral/the-host-go/internal/models/api"
	"github.com/tugkanmeral/the-host-go/internal/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
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

	token, err := h.auth.Login(c.UserContext(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			return c.Status(fiber.StatusUnauthorized).JSON(model.ErrorResponse{
				Message: "Kullanıcı adı veya şifre hatalı",
			})
		case errors.Is(err, service.ErrTokenGeneration):
			return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
				Message: "Token oluşturulurken bir hata oluştu",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
				Message: "Kullanıcı aranırken bir hata oluştu",
			})
		}
	}

	return c.JSON(model.LoginResponse{
		Token: token,
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
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

	err := h.auth.Register(c.UserContext(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUsernameTaken):
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Message: "Username conflict",
			})
		case errors.Is(err, service.ErrRegisterFailed):
			return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
				Message: "User cannot be registered!",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
				Message: "User cannot be registered!",
			})
		}
	}

	return nil
}
