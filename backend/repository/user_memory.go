package repository

import (
	"errors"
	"log"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryInMemory struct {
	data   map[uint32]*model.User
	nextId uint32
}

func NewUserRepositoryInMemory() User {
	return &UserRepositoryInMemory{
		data:   make(map[uint32]*model.User),
		nextId: 1,
	}
}

func (repo *UserRepositoryInMemory) Create(user *model.User) error {
	nextId := repo.nextId
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Println(err)
		return err
	}

	user.Id = nextId

	user.Password = string(hashedPassword)
	repo.data[nextId] = user

	repo.nextId += 1
	return nil
}

func (repo *UserRepositoryInMemory) Get(id uint32) (*model.User, error) {
	if user, ok := repo.data[id]; ok {
		return user, nil
	}
	return nil, ErrUserNotFound
}

func (repo *UserRepositoryInMemory) GetAll() ([]model.User, error) {
	allUsers := make([]model.User, 0)
	for _, user := range repo.data {
		allUsers = append(allUsers, model.User{
			Id:       user.Id,
			Fullname: user.Fullname,
			Username: user.Username,
			Password: user.Password,
			Email:    user.Email,
			Telegram: user.Telegram,
		})
	}
	return allUsers, nil
}

func (repo *UserRepositoryInMemory) SearchByUsername(username string) (*model.User, error) {
	for _, v := range repo.data {
		if v.Username == username {
			return v, nil
		}
	}
	return nil, ErrUserNotFound
}

func (repo *UserRepositoryInMemory) ValidateUser(username, password string) (*model.User, error) {
	user, err := repo.SearchByUsername(username)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func (repo *UserRepositoryInMemory) Update(user *model.User) error {
	return errors.New("not implemented")
}
func (repo *UserRepositoryInMemory) Delete(id uint32) error {
	return errors.New("not implemented")
}

func (r *UserRepositoryInMemory) ListApiKeys() ([]model.ApiKey, error) {
	return nil, ErrNotImplemented
}

func (r *UserRepositoryInMemory) AddApiKey(*model.ApiKey) error {
	return ErrNotImplemented
}

func (r *UserRepositoryInMemory) GetApiKey(id, userId uint32) (*model.ApiKey, error) {
	return nil, ErrNotImplemented
}

func (r *UserRepositoryInMemory) SaveApiKey(apiKey *model.ApiKey) error {
	return ErrNotImplemented
}
