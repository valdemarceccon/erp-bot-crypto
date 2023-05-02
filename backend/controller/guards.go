package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/valdemarceccon/crypto-bot-erp/backend/middleware/constants"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
)

func WithPermission(repo repository.Role, permissonName model.Permission, handler func(c *fiber.Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := c.Locals(constants.ContextKeyCurrentUser).(*model.User)

		if user == nil {
			log.Printf("controller: guard: user_obj not set")
			return fiber.ErrUnauthorized
		}

		permissions, err := repo.UserPermissions(user.Id)

		if err != nil {
			log.Println("controller: guard:", err)
			return c.SendStatus(fiber.StatusForbidden)
		}
		for _, role := range permissions {
			if role == model.Permission(permissonName) {
				return handler(c)
			}
		}
		log.Printf("controller: guard: user '%s' don't have '%s' permission", user.Username, permissonName)
		return c.SendStatus(fiber.StatusForbidden)
	}
}
