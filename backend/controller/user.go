package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller/schema"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
)

type UserController interface {
	ListUsers(c *fiber.Ctx) error
}

type UserControllerImpl struct {
	userRepository repository.UserRepository
}

func NewUserController(ur repository.UserRepository) UserController {
	return &UserControllerImpl{
		userRepository: ur,
	}
}

func (uc *UserControllerImpl) ListUsers(c *fiber.Ctx) error {

	users, err := uc.userRepository.GetAll()

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	ret := make([]*schema.UserResponse, 0)
	for _, u := range users {
		ret = append(ret, schema.FromUserModel(&u))
	}

	return c.JSON(ret)
}
