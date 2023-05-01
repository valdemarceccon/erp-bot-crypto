package repository

import (
	"database/sql"
	"errors"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
)

type UserRepositoryPsql struct {
	db *sql.DB
}

func NewUserRepositoryPsql(db *sql.DB) UserRepository {
	return &UserRepositoryPsql{
		db: db,
	}
}

var (
	ErrNotImplemented = errors.New("method not implemented")
)

func (ur *UserRepositoryPsql) Create(user *model.User) error {

	return ErrNotImplemented
}

func (ur *UserRepositoryPsql) Get(id uint32) (*model.User, error) {

	return nil, ErrNotImplemented
}

func (ur *UserRepositoryPsql) GetAll() ([]model.User, error) {
	rows, err := ur.db.Query(`
		SELECT
			id,
			email,
			username,
			name,
			hashed_password
		FROM
			users`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret := make([]model.User, 0)

	for rows.Next() {
		var user model.User

		err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.Name, &user.Password)

		if err != nil {
			return nil, err
		}

		ret = append(ret, user)
	}

	return ret, err
}

func (ur *UserRepositoryPsql) Update(user *model.User) error {

	return ErrNotImplemented
}

func (ur *UserRepositoryPsql) Delete(id uint32) error {

	return ErrNotImplemented
}

func (ur *UserRepositoryPsql) SearchByUsername(string) (*model.User, error) {

	return nil, ErrNotImplemented
}

func (ur *UserRepositoryPsql) ValidateUser(username, password string) (*model.User, error) {

	return nil, ErrNotImplemented
}
