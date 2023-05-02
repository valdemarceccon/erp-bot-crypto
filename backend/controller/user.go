package controller

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller/schema"
	"github.com/valdemarceccon/crypto-bot-erp/backend/middleware/constants"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
)

type UserController struct {
	userRepository repository.User
	roleRepository repository.Role
}

func NewUserController(ur repository.User, role repository.Role) *UserController {
	return &UserController{
		userRepository: ur,
		roleRepository: role,
	}
}

func (uc *UserController) ListUsers(c *fiber.Ctx) error {

	users, err := uc.userRepository.GetAll()

	if err != nil {
		log.Println("user: ", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	ret := make([]*schema.UserResponse, 0)
	for _, u := range users {
		ret = append(ret, schema.FromUserModel(&u))
	}

	return c.JSON(ret)
}

func (uc *UserController) Me(c *fiber.Ctx) error {
	user := c.Locals(constants.ContextKeyCurrentUser).(*model.User)

	return c.JSON(schema.FromUserModel(user))
}

func (uc *UserController) ListApiKeys(c *fiber.Ctx) error {
	resp, err := uc.userRepository.ListApiKeys()

	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	result := make([]schema.ApiKeyResponse, 0)

	for _, v := range resp {
		result = append(result, *schema.FromApiKeyModel(&v))
	}

	return c.JSON(result)
}
