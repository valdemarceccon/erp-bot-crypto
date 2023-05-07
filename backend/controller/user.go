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
	"github.com/valdemarceccon/crypto-bot-erp/backend/store"
)

type UserController struct {
	userStore store.User
	roleStore store.Role
	apiStore  store.ApiKey
	validate  *validator.Validate
}

func getCurrentUserFromContext(c *fiber.Ctx) *model.User {
	return c.Locals(constants.ContextKeyCurrentUser).(*model.User)
}

func NewUserController(ur store.User, role store.Role, apiKey store.ApiKey) *UserController {
	return &UserController{
		userStore: ur,
		roleStore: role,
		apiStore:  apiKey,
		validate:  validator.New(),
	}
}

func (uc *UserController) ListUsers(c *fiber.Ctx) error {

	users, err := uc.userStore.List()

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
	permitions, err := uc.roleStore.FromUser(user.Id)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	respPermissions := make([]schema.PermissionResponse, len(permitions))

	for i, v := range permitions {
		respPermissions[i] = schema.PermissionResponse{
			Name: string(v),
		}
	}

	return c.JSON(schema.UserMeResponse{
		UserResponse: *schema.FromUserModel(user),
		Permissions:  respPermissions,
	})
}

func (uc *UserController) ListApiKeys(c *fiber.Ctx) error {
	user := getCurrentUserFromContext(c)

	resp, err := uc.apiStore.FromUser(user.Id)

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
	resp, err := uc.apiStore.List()

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

	err = uc.apiStore.New(&model.ApiKey{
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

	return c.Status(fiber.StatusCreated).JSON(nil)
}

func (uc *UserController) ClientToggleApiKey(c *fiber.Ctx) error {
	apiKeyId, err := c.ParamsInt("apiKeyId")
	userId := getCurrentUserFromContext(c)

	if err != nil {
		log.Println(err)
		return fiber.ErrBadRequest
	}

	apiKey, err := uc.apiStore.Get(uint32(apiKeyId), userId.Id)

	if err != nil {
		log.Println(err)
		if err == store.ErrApiKeyNotFound {
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

	err = uc.apiStore.Save(apiKey)

	if err != nil {
		fmt.Println(err)
		if err == store.ErrApiKeyNotFound {
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

	balanceETH := balanceList[0].Coin[0].WalletBalance
	f := new(big.Float)
	parsedBalance, _, err := f.Parse(balanceETH, 10)
	if err != nil {
		return nil, err
	}
	return parsedBalance, nil
}

func (uc *UserController) getApiKeyForUser(c *fiber.Ctx) (*model.ApiKey, error) {
	apiKeyId, err := c.ParamsInt("apiKeyId")

	if err != nil {
		log.Println(err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "missing api key id")
	}

	userId, err := c.ParamsInt("userId")

	if err != nil {
		log.Println(err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "missing user id")
	}

	apiKey, err := uc.apiStore.Get(uint32(apiKeyId), uint32(userId))

	if err != nil {
		log.Println(err)
		if err == store.ErrApiKeyNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "api key not found for the user")
		}

		return nil, fiber.ErrInternalServerError
	}

	return apiKey, nil
}

func (uc *UserController) updateBotStatus(apiKey *model.ApiKey, balance *big.Float, newStatus model.ApiKeyStatus) error {
	if newStatus == model.ApiKeyStatusActive {
		err := uc.userStore.StartBot(apiKey, balance)

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("was not able to register the bot start. %v", err))

		}
	} else if newStatus == model.ApiKeyStatusInactive {
		err := uc.userStore.StopBot(apiKey, balance)

		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("was not able to register the bot stop. %v", err))
		}
	}

	apiKey.Status = newStatus

	return uc.apiStore.Save(apiKey)
}

func (uc *UserController) AdminToggleApiKey(c *fiber.Ctx) error {

	apiKey, err := uc.getApiKeyForUser(c)

	if err != nil {
		return err
	}

	balance, err := getWalletBalanceETH(apiKey)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("was not able to fetch wallet ballance. %v", err))
	}

	newStatus := getNewApiKeyStatus(apiKey.Status)

	err = uc.updateBotStatus(apiKey, balance, newStatus)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(schema.FromApiKeyModel(apiKey))
}

func getNewApiKeyStatus(status model.ApiKeyStatus) model.ApiKeyStatus {
	var newStatus model.ApiKeyStatus
	switch status {
	case model.ApiKeyStatusActive:
		newStatus = model.ApiKeyStatusInactive
	case model.ApiKeyStatusInactive:
		newStatus = model.ApiKeyStatusActive
	case model.ApiKeyStatusWaitingActivation:
		newStatus = model.ApiKeyStatusActive
	case model.ApiKeyStatusWaitingDeactivation:
		newStatus = model.ApiKeyStatusInactive
	}

	return newStatus
}

func (uc *UserController) CalculateComission(c fiber.Ctx) error {
	return fiber.NewError(fiber.StatusInternalServerError, store.ErrNotImplemented.Error())
}
