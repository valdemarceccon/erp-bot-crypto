package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
)

type UserPsql struct {
	db *sql.DB
}

func NewUserPsql(db *sql.DB) User {
	return &UserPsql{
		db: db,
	}
}

var (
	ErrNotImplemented = errors.New("method not implemented")
)

func (ur *UserPsql) Create(user *model.User) error {
	queryExists, err := ur.db.Query("SELECT 1 FROM users WHERE username = $1 or email = $2;", user.Username, user.Email)
	if err != nil {
		log.Println(err)
		return err
	}
	defer queryExists.Close()

	if queryExists.Next() {
		return ErrUserOrEmailInUse
	}

	row := ur.db.QueryRow("INSERT INTO users(email,telegram,username,fullname,hashed_password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,now(),now()) RETURNING id;", user.Email, user.Telegram, user.Username, user.Fullname, user.Password)

	return row.Scan(&user.Id)
}

func (ur *UserPsql) Get(id uint32) (*model.User, error) {

	row := ur.db.QueryRow(`
	SELECT 	id,
					email,
					username,
					fullname,
					hashed_password
	FROM users
	WHERE id = $1
		AND deleted_at is null;`, id)

	var ret model.User

	err := row.Scan(
		&ret.Id,
		&ret.Email,
		&ret.Username,
		&ret.Fullname,
		&ret.Password,
	)

	if err == sql.ErrNoRows {
		log.Println(err)
		return nil, ErrUserNotFound
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &ret, nil
}

func (ur *UserPsql) GetAll() ([]model.User, error) {
	rows, err := ur.db.Query(`
		SELECT
			id,
			email,
			username,
			fullname,
			hashed_password
		FROM
			users
		WHERE
			deleted_at is null;`)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	ret := make([]model.User, 0)

	for rows.Next() {
		var user model.User

		err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.Fullname, &user.Password)

		if err != nil {
			log.Println(err)
			return nil, err
		}

		ret = append(ret, user)
	}

	return ret, nil
}

func (ur *UserPsql) Update(user *model.User) error {

	return ErrNotImplemented
}

func (ur *UserPsql) Delete(id uint32) error {

	return ErrNotImplemented
}

func (ur *UserPsql) SearchByUsername(username string) (*model.User, error) {
	row := ur.db.QueryRow(`
		select
			id,
			email,
			telegram,
			username,
			fullname,
			hashed_password
		from users
		where username = $1
			and deleted_at is null`, username)
	resp := &model.User{}
	err := row.Scan(&resp.Id, &resp.Email, &resp.Telegram, &resp.Username, &resp.Fullname, &resp.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return resp, nil
}

func (ur *UserPsql) ListApiKeys() ([]model.ApiKey, error) {
	row, err := ur.db.Query(`
		select
			id,
			user_id,
			api_key_name,
			exchange,
			api_key,
			api_secret,
			status
		from api_key
		where deleted_at is null`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer row.Close()
	resp := make([]model.ApiKey, 0)
	for row.Next() {
		var apiKey model.ApiKey

		err = row.Scan(&apiKey.Id, &apiKey.UserId, &apiKey.ApiKeyName, &apiKey.Exchange, &apiKey.ApiKey, &apiKey.ApiSecret, &apiKey.Status)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		resp = append(resp, apiKey)
	}

	return resp, nil
}
