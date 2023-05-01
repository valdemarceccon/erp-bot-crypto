package controller

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller/schema"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthController interface {
	LoginHandler(*fiber.Ctx) error
	RegisterHandler(*fiber.Ctx) error
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

func (jac *JwtAuthController) RegisterHandler(c *fiber.Ctx) error {
	var registerRequest schema.RegisterRequest
	err := c.BodyParser(&registerRequest)

	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	newUser := &model.User{
		Name:     registerRequest.Fullname,
		Username: registerRequest.Username,
		Password: string(hashedPassword),
		Email:    registerRequest.Email,
		Telegram: registerRequest.Telegram,
	}

	err = jac.UserRepository.Create(newUser)

	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(schema.RegisterReponse{
		Message: "ok",
	})
}

func (jac *JwtAuthController) LoginHandler(c *fiber.Ctx) error {
	var loginVal schema.LoginRequest
	err := c.BodyParser(&loginVal)

	if err != nil {
		log.Println(err)
		return fiber.ErrUnauthorized
	}

	user, err := jac.validateUser(loginVal.Username, loginVal.Password)

	if err != nil {
		log.Println(err)
		return fiber.ErrForbidden
	}

	signedToken, err := generateUserToken(jac.JwtSecret, user)

	if err != nil {
		log.Println(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(schema.LoginResponse{
		Token: signedToken,
	})
}

func (jac *JwtAuthController) validateUser(username, password string) (*model.User, error) {
	dbUser, err := jac.UserRepository.SearchByUsername(username)

	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)); err != nil {
		return nil, err
	}
	return dbUser, nil
}

func generateUserToken(secret string, user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"username": user.Username,
		"sub":      user.Id,
		"exp":      time.Now().Add(time.Hour * 12).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}
