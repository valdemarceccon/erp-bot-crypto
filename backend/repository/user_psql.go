package repository

import (
	"database/sql"
	"errors"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
)

type UserRepositoryPsql struct {
	db *sql.DB
}

func NewUserPsql(db *sql.DB) User {
	return &UserRepositoryPsql{
		db: db,
	}
}

var (
	ErrNotImplemented = errors.New("method not implemented")
)

func (ur *UserRepositoryPsql) Create(user *model.User) error {
	row := ur.db.QueryRow("INSERT INTO users(email,telegram,username,fullname,hashed_password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,now(),now()) RETURNING id;", user.Email, user.Telegram, user.Username, user.Fullname, user.Password)

	return row.Scan(&user.Id)
}

func (ur *UserRepositoryPsql) Get(id uint32) (*model.User, error) {

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
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func (ur *UserRepositoryPsql) GetAll() ([]model.User, error) {
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
		return nil, err
	}
	defer rows.Close()

	ret := make([]model.User, 0)

	for rows.Next() {
		var user model.User

		err := rows.Scan(&user.Id, &user.Email, &user.Username, &user.Fullname, &user.Password)

		if err != nil {
			return nil, err
		}

		ret = append(ret, user)
	}

	return ret, nil
}

func (ur *UserRepositoryPsql) Update(user *model.User) error {

	return ErrNotImplemented
}

func (ur *UserRepositoryPsql) Delete(id uint32) error {

	return ErrNotImplemented
}

func (ur *UserRepositoryPsql) SearchByUsername(username string) (*model.User, error) {
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
		return nil, err
	}
	return resp, nil
}
