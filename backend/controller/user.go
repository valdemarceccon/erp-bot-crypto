package controller

import (
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/hirokisan/bybit/v2"
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
	permitions, err := uc.userRepository.ListUsersPermission(user.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	respPermissions := make([]schema.PermissionResponse, len(permitions))

	for i, v := range permitions {
		respPermissions[i] = schema.PermissionResponse{
			Name: string(v),
		}
	}

	return c.JSON(schema.UserMeReponse{
		UserResponse: *schema.FromUserModel(user),
		Permissions:  respPermissions,
	})
}

func (uc *UserController) ListApiKeys(c *fiber.Ctx) error {
	user := getCurrentUserFromContext(c)

	resp, err := uc.userRepository.ListUserApiKeys(user.Id)

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

func (uc *UserController) ListAllApiKeys(c *fiber.Ctx) error {
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

func getWalletBalanceETH(apiKey *model.ApiKey) (*big.Float, error) {
	balance, err := bybit.
		NewClient().
		WithAuth(apiKey.ApiKey, apiKey.ApiSecret).
		V5().
		Account().
		GetWalletBalance(
			bybit.AccountType(bybit.AccountTypeNormal), []bybit.Coin{bybit.CoinETH},
		)

	if err != nil {
		return nil, err
	}

	balanceList := balance.Result.List

	if len(balanceList) == 0 || len(balanceList[0].Coin) == 0 || balanceList[0].Coin[0].WalletBalance == "" {
		return nil, errors.New("no balance info available")
	}

	for _, v := range balanceList {
		for _, c := range v.Coin {
			fmt.Println(c.WalletBalance)
		}
	}
	balanceETH := balanceList[0].Coin[0].WalletBalance
	parsedBalance, _, err := big.ParseFloat(balanceETH, 10, 30, big.ToNearestAway)
	return parsedBalance, nil
}

func (uc *UserController) AdminToggleApiKey(c *fiber.Ctx) error {
	apiKeyId, err := c.ParamsInt("apiKeyId")

	if err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	userId, err := c.ParamsInt("userId")

	if err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	apiKey, err := uc.userRepository.GetApiKey(uint32(apiKeyId), uint32(userId))

	if err != nil {
		log.Println(err)
		if err == repository.ErrApiKeyNotFound {
			return fiber.ErrNotFound
		}

		return fiber.ErrInternalServerError
	}

	var newStatus model.ApiKeyStatus
	switch apiKey.Status {
	case model.ApiKeyStatusActive:
		newStatus = model.ApiKeyStatusInactive
	case model.ApiKeyStatusInactive:
		newStatus = model.ApiKeyStatusActive
	case model.ApiKeyStatusWaitingActivation:
		newStatus = model.ApiKeyStatusActive
	case model.ApiKeyStatusWaitingDeactivation:
		newStatus = model.ApiKeyStatusInactive
	}

	balance, err := getWalletBalanceETH(apiKey)
	if err != nil {
		log.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("was not able to fetch wallet ballance. %v", err))
	}

	if newStatus == model.ApiKeyStatusActive {
		err = uc.userRepository.StartBot(apiKey, balance)

		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("was not able to register the bot start. %v", err))

		}
	} else if newStatus == model.ApiKeyStatusInactive {
		err = uc.userRepository.StopBot(apiKey, balance)

		if err != nil {
			log.Println(err)
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("was not able to register the bot stop. %v", err))

		}
	}

	apiKey.Status = newStatus

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
