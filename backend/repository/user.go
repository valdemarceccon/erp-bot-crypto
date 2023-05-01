package repository

import (
	"errors"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
)

var (
	ErrUserNotFound = errors.New("user: User not found")
)

type User struct {
	model.User
	Timestamps
}

func toDomainUser(dbUser *User) *model.User {
	return &model.User{
		Id:       dbUser.Id,
		Name:     dbUser.Name,
		Username: dbUser.Username,
		Password: dbUser.Password,
		Email:    dbUser.Email,
		Telegram: dbUser.Telegram,
	}
}

func toDBModel(user *model.User) *User {
	return &User{
		User: model.User{
			Id:       user.Id,
			Name:     user.Name,
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
			Telegram: user.Telegram,
		},
	}

}

type UserRepository interface {
	Create(user *model.User) error
	Get(id uint32) (*model.User, error)
	GetAll() ([]model.User, error)
	Update(user *model.User) error
	Delete(id uint32) error
	SearchByUsername(string) (*model.User, error)
}
