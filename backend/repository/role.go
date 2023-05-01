package repository

import "github.com/valdemarceccon/crypto-bot-erp/backend/model"

type Transactions string

const (
	ListUsers Transactions = "ListUsers"
)

type Role interface {
	Create(user *model.Role) error
	Get(id uint32) (*model.Role, error)
	GetAll() ([]model.Role, error)
	Update(user *model.Role) error
	Delete(id uint32) error
	SearchByName(string) (*model.Role, error)
	UserPermissions(userId uint32) ([]model.Permission, error)
}
