package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller/schema"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
)

type AuthController interface {
	LoginHandler(*fiber.Ctx) error
}

type JwtAuthController struct {
	UserRepository repository.UserRepository
	JwtSecret      string
}

type JwtAuthControllerOptions func(*JwtAuthController)

func WithHS256Secret(secret string) JwtAuthControllerOptions {
	return func(jac *JwtAuthController) {
		jac.JwtSecret = secret
	}
}

func NewJwtAuthController(ur repository.UserRepository, options ...JwtAuthControllerOptions) AuthController {
	var controller JwtAuthController

	for _, opt := range options {
		opt(&controller)
	}

	controller.UserRepository = ur

	return &controller
}

func (jac *JwtAuthController) LoginHandler(c *fiber.Ctx) error {
	var loginVal schema.LoginRequest
	err := c.BodyParser(&loginVal)

	if err != nil {
		return fiber.ErrUnauthorized
	}

	user, err := jac.UserRepository.ValidateUser(loginVal.Username, loginVal.Password)

	if err != nil {
		return fiber.ErrForbidden
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"username": user.Username,
		"sub":      user.Id,
		"exp":      time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jac.JwtSecret))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(schema.LoginResponse{
		Token: signedToken,
	})
}
