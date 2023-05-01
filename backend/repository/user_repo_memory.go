package repository

import (
	"errors"
	"time"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryInMemory struct {
	data   map[uint32]*User
	nextId uint32
}

func NewUserRepositoryInMemory() UserRepository {
	return &UserRepositoryInMemory{
		data:   make(map[uint32]*User),
		nextId: 1,
	}
}

func (repo *UserRepositoryInMemory) Create(user *model.User) error {
	nextId := repo.nextId
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user.Id = nextId

	now := time.Now()
	newUser := toDBModel(user)
	newUser.CreatedAt = &now
	newUser.UpdatedAt = &now
	newUser.Password = string(hashedPassword)
	repo.data[nextId] = newUser

	repo.nextId += 1
	return nil
}

func (repo *UserRepositoryInMemory) Get(id uint32) (*model.User, error) {
	if user, ok := repo.data[id]; ok {
		return toDomainUser(user), nil
	}
	return nil, ErrUserNotFound
}

func (repo *UserRepositoryInMemory) GetAll() ([]model.User, error) {
	allUsers := make([]model.User, 0)
	for _, user := range repo.data {
		allUsers = append(allUsers, model.User{
			Id:       user.Id,
			Name:     user.Name,
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
			return toDomainUser(v), nil
		}
	}
	return nil, ErrUserNotFound
}

func (repo *UserRepositoryInMemory) ValidateUser(username, password string) (*model.User, error) {
	user, err := repo.SearchByUsername(username)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
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
