package repository

import (
	"database/sql"
	"errors"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"golang.org/x/crypto/bcrypt"
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
	row := ur.db.QueryRow("INSERT INTO users(email,telegram,username,fullname,hashed_password) VALUES ($1,$2,$3,$4,$5) RETURNING id;", user.Email, user.Telegram, user.Username, user.Name, user.Password)

	return row.Scan(&user.Id)
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

func (ur *UserRepositoryPsql) SearchByUsername(username string) (*model.User, error) {
	row := ur.db.QueryRow("select id,	email,	telegram,	username,	fullname,	hashed_password from users where username = $1", username)
	resp := &model.User{}
	err := row.Scan(&resp.Id, &resp.Email, &resp.Telegram, &resp.Username, &resp.Name, &resp.Password)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (ur *UserRepositoryPsql) ValidateUser(username, password string) (*model.User, error) {
	dbUser, err := ur.SearchByUsername(username)

	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password)); err != nil {
		return nil, err
	}
	return dbUser, nil
}
