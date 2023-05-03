package repository

import (
	"errors"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
)

var (
	ErrUserNotFound     = errors.New("user: User not found")
	ErrUserOrEmailInUse = errors.New("user: Username or Email taken")

	ErrApiKeyNotFound                = errors.New("user: Api key not found")
	ErrCouldNotUpdateApikey          = errors.New("could not update api_key")
	ErrCouldNoteRetrieveAffectedRows = errors.New("could not retrieve rows affected")
)

type User interface {
	Create(user *model.User) error
	Get(id uint32) (*model.User, error)
	GetAll() ([]model.User, error)
	Update(user *model.User) error
	Delete(id uint32) error
	SearchByUsername(string) (*model.User, error)
	ListApiKeys() ([]model.ApiKey, error)
	ListUserApiKeys(uint32) ([]model.ApiKey, error)

	AddApiKey(*model.ApiKey) error
	GetApiKey(id, userId uint32) (*model.ApiKey, error)
	SaveApiKey(apiKey *model.ApiKey) error
	ListUsersPermission(userId uint32) ([]model.Permission, error)
}
