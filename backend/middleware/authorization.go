package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
)

type Auth struct {
	userRepo repository.User
	roleRepo repository.Role
	jwtKey   string
}

func NewAuthMiddleware(userRepo repository.User, roleRepo repository.Role, jwtKey string) *Auth {
	return &Auth{
		userRepo: userRepo,
		jwtKey:   jwtKey,
		roleRepo: roleRepo,
	}
}

func (a *Auth) UserExists(c *fiber.Ctx) error {
	token := c.Locals("user").(*jwt.Token)

	if token == nil {
		log.Println("auth: token not set")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var claims controller.UserClaim
	jwt.ParseWithClaims(token.Raw, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(a.jwtKey), nil
	})

	user, err := a.userRepo.Get(claims.UserId)

	if err != nil {
		log.Printf("auth: %s", err)
		if err == repository.ErrUserNotFound {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if user == nil {
		log.Printf("auth: user is nil")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Locals("user_obj", user)

	return c.Next()
}
