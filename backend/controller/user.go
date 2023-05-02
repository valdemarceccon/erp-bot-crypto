package controller

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/valdemarceccon/crypto-bot-erp/backend/controller/schema"
	"github.com/valdemarceccon/crypto-bot-erp/backend/middleware/constants"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
)

type UserController struct {
	userRepository repository.User
	roleRepository repository.Role
	validate       *validator.Validate
}

func getCurrentUserFromContext(c *fiber.Ctx) *model.User {
	return c.Locals(constants.ContextKeyCurrentUser).(*model.User)
}

func NewUserController(ur repository.User, role repository.Role) *UserController {
	return &UserController{
		userRepository: ur,
		roleRepository: role,
		validate:       validator.New(),
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

func (uc *UserController) AddApiKey(c *fiber.Ctx) error {
	user := getCurrentUserFromContext(c)
	var requestBody schema.ApiKeyRequest
	err := c.BodyParser(&requestBody)

	if err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	// TODO: validate request body

	err = uc.userRepository.AddApiKey(&model.ApiKey{
		UserId:     user.Id,
		ApiKeyName: requestBody.ApiKeyName,
		Exchange:   requestBody.Exchange,
		ApiKey:     requestBody.ApiKey,
		ApiSecret:  requestBody.ApiKeySecret,
		Status:     model.ApiKeyStatusInactive,
	})

	if err != nil {
		log.Println(err)
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (uc *UserController) ClientToggleApiKey(c *fiber.Ctx) error {
	apiKeyId, err := c.ParamsInt("apiKeyId")
	userId := getCurrentUserFromContext(c)

	if err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	apiKey, err := uc.userRepository.GetApiKey(uint32(apiKeyId), userId.Id)

	if err != nil {
		log.Println(err)
		if err == repository.ErrApiKeyNotFound {
			return fiber.ErrNotFound
		}

		return fiber.ErrInternalServerError
	}

	switch apiKey.Status {
	case model.ApiKeyStatusActive:
		apiKey.Status = model.ApiKeyStatusWaitingDeactivation
	case model.ApiKeyStatusInactive:
		apiKey.Status = model.ApiKeyStatusWaitingActivation
	case model.ApiKeyStatusWaitingActivation:
		apiKey.Status = model.ApiKeyStatusInactive
	case model.ApiKeyStatusWaitingDeactivation:
		apiKey.Status = model.ApiKeyStatusActive
	}

	err = uc.userRepository.SaveApiKey(apiKey)

	if err != nil {
		fmt.Println(err)
		if err == repository.ErrApiKeyNotFound {
			return fiber.ErrNotFound
		}

		return fiber.ErrInternalServerError
	}

	return c.JSON(schema.FromApiKeyModel(apiKey))
}
