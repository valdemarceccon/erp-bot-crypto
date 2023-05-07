package store

import (
	"errors"
	"math/big"

	"github.com/hirokisan/bybit/v2"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
)

var (
	ErrUserNotFound     = errors.New("user: user not found")
	ErrUserOrEmailInUse = errors.New("user: username or Email taken")

	ErrApiKeyNotFound                = errors.New("user: api key not found")
	ErrCouldNotUpdateApikey          = errors.New("could not update api_key")
	ErrCouldNoteRetrieveAffectedRows = errors.New("could not retrieve rows affected")

	ErrNotImplemented = errors.New("method not implemented")
)

type Transactions string

const (
	ListUsers Transactions = "ListUsers"
)

type User interface {
	New(user *model.User) error
	Get(id uint32) (*model.User, error)
	List() ([]model.User, error)
	Save(user *model.User) error
	Delete(id uint32) error
	ByUsername(string) (*model.User, error)

	SaveClosedPnL(userId, apiKeyId uint32, pnlResult []bybit.V5GetClosedPnLItem) error
	StartBot(apikey *model.ApiKey, balance *big.Float) error
	StopBot(apikey *model.ApiKey, balance *big.Float) error
}

type ApiKey interface {
	New(*model.ApiKey) error
	Get(id, userId uint32) (*model.ApiKey, error)
	List() ([]model.ApiKey, error)
	Save(apiKey *model.ApiKey) error
	FromUser(uint32) ([]model.ApiKey, error)
	ListActive(uint32) ([]model.ApiKey, error)
}

type Role interface {
	New(user *model.Role) error
	Get(id uint32) (*model.Role, error)
	List() ([]model.Role, error)
	Save(user *model.Role) error
	Delete(id uint32) error
	ByName(string) (*model.Role, error)
	FromUser(userId uint32) ([]model.Permission, error)
	// ListUsersPermission(userId uint32) ([]model.Permission, error)
}
