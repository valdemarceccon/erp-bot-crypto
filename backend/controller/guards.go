package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/store"
)

type Guards struct {
	roleStore store.Role
}

func NewGuards(roleStore store.Role) *Guards {
	return &Guards{
		roleStore: roleStore,
	}
}

func (g *Guards) WithPermission(permissonName model.Permission, handler func(c *fiber.Ctx) error) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		user := getCurrentUserFromContext(c)

		if user == nil {
			log.Printf("controller: guard: user_obj not set")
			return fiber.ErrUnauthorized
		}

		permissions, err := g.roleStore.FromUser(user.Id)

		if err != nil {
			log.Println("controller: guard:", err)
			return fiber.ErrForbidden
		}
		for _, role := range permissions {
			if role == model.Permission(permissonName) {
				return handler(c)
			}
		}
		log.Printf("controller: guard: user '%s' don't have '%s' permission", user.Username, permissonName)
		return fiber.ErrForbidden
	}
}
