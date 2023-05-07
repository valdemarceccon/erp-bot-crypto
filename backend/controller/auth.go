package controller

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller/schema"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/store"
	"golang.org/x/crypto/bcrypt"
)

type AuthController interface {
	LoginHandler(*fiber.Ctx) error
	RegisterHandler(*fiber.Ctx) error
}

type JwtAuthController struct {
	UserStore store.User
	JwtSecret string
}

type JwtAuthControllerOptions func(*JwtAuthController)

func WithHS256Secret(secret string) JwtAuthControllerOptions {
	return func(jac *JwtAuthController) {
		jac.JwtSecret = secret
	}
}

func NewJwtAuthController(ur store.User, options ...JwtAuthControllerOptions) AuthController {
	var controller JwtAuthController

	for _, opt := range options {
		opt(&controller)
	}

	controller.UserStore = ur

	return &controller
}

func (jac *JwtAuthController) RegisterHandler(c *fiber.Ctx) error {
	var registerRequest schema.RegisterRequest
	err := c.BodyParser(&registerRequest)

	if err != nil {
		log.Println(err)
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Println(err)
		return err
	}

	newUser := &model.User{
		Fullname: registerRequest.Fullname,
		Username: registerRequest.Username,
		Password: string(hashedPassword),
		Email:    registerRequest.Email,
		Telegram: registerRequest.Telegram,
	}

	err = jac.UserStore.New(newUser)

	if err != nil {
		log.Println(err)
		if err == store.ErrUserOrEmailInUse {

			return fiber.NewError(fiber.StatusForbidden, "Username or email already in use")
		}

		return fiber.ErrInternalServerError
	}

	token, err := generateUserToken(jac.JwtSecret, newUser)

	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(schema.LoginResponse{Token: token})
}

func (jac *JwtAuthController) LoginHandler(c *fiber.Ctx) error {
	var loginVal schema.LoginRequest
	err := c.BodyParser(&loginVal)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(schema.MessageResponse{
			Message: "invalid username or password",
		})
	}

	user, err := jac.validateUser(loginVal.Username, loginVal.Password)

	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusForbidden).JSON(schema.MessageResponse{
			Message: "invalid username or password",
		})
	}

	signedToken, err := generateUserToken(jac.JwtSecret, user)

	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.JSON(schema.LoginResponse{
		Token: signedToken,
	})
}

func (jac *JwtAuthController) validateUser(username, password string) (*model.User, error) {
	dbUser, err := jac.UserStore.ByUsername(username)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)); err != nil {
		log.Println(err)
		return nil, err
	}
	return dbUser, nil
}

type UserClaim struct {
	jwt.RegisteredClaims
	Username string
	UserId   uint32
}

func generateUserToken(secret string, user *model.User) (string, error) {
	claims := UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		},
		Username: user.Username,
		UserId:   user.Id,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))

}
