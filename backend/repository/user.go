package repository

import (
	"errors"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
)

var (
	ErrUserNotFound     = errors.New("user: User not found")
	ErrUserOrEmailInUse = errors.New("user: Username or Email taken")
)

type User interface {
	Create(user *model.User) error
	Get(id uint32) (*model.User, error)
	GetAll() ([]model.User, error)
	Update(user *model.User) error
	Delete(id uint32) error
	SearchByUsername(string) (*model.User, error)
	ListApiKeys() ([]model.ApiKey, error)
	AddApiKey(*model.ApiKey) error
}
