package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller"
	"github.com/valdemarceccon/crypto-bot-erp/backend/middleware/constants"
	"github.com/valdemarceccon/crypto-bot-erp/backend/store"
)

type Auth struct {
	userStore store.User
	roleStore store.Role
	jwtKey    string
}

func NewAuthMiddleware(userStore store.User, roleStore store.Role, jwtKey string) *Auth {
	return &Auth{
		userStore: userStore,
		jwtKey:    jwtKey,
		roleStore: roleStore,
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

	user, err := a.userStore.Get(claims.UserId)

	if err != nil {
		log.Printf("auth: %s", err)
		if err == store.ErrUserNotFound {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if user == nil {
		log.Printf("auth: user is nil")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Locals(constants.ContextKeyCurrentUser, user)

	return c.Next()
}
