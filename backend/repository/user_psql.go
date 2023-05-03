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
			deleted_at is null
		order by
			id;`)

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

func (ur *UserPsql) ListUserApiKeys(userId uint32) ([]model.ApiKey, error) {
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
		where deleted_at is null
		and user_id = $1
		order by id, user_id;`, userId)
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

func (ur *UserPsql) ListApiKeys() ([]model.ApiKey, error) {
	row, err := ur.db.Query(`
		select
			ak.id,
			user_id,
			u.username,
			api_key_name,
			exchange,
			api_key,
			api_secret,
			status
		from api_key ak
		inner join users u on
		u.id = ak.user_id
		where ak.deleted_at is null
		and u.deleted_at is null
		order by id, user_id;`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer row.Close()
	resp := make([]model.ApiKey, 0)
	for row.Next() {
		var apiKey model.ApiKey

		err = row.Scan(&apiKey.Id, &apiKey.UserId, &apiKey.Username, &apiKey.ApiKeyName, &apiKey.Exchange, &apiKey.ApiKey, &apiKey.ApiSecret, &apiKey.Status)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		resp = append(resp, apiKey)
	}

	return resp, nil
}

func (ur *UserPsql) AddApiKey(apiKey *model.ApiKey) error {
	row := ur.db.QueryRow(`INSERT INTO api_key (
		user_id,
		api_key_name,
		exchange,
		api_key,
		api_secret,
		status,
		created_at,
		updated_at,
		deleted_at
	) VALUES ($1,$2,$3,$4,$5,$6,now(),now(),null) RETURNING id`, apiKey.UserId, apiKey.ApiKeyName, apiKey.Exchange, apiKey.ApiKey, apiKey.ApiSecret, apiKey.Status)

	return row.Scan(&apiKey.Id)
}

func (ur *UserPsql) GetApiKey(id, userId uint32) (*model.ApiKey, error) {
	row := ur.db.QueryRow(`
	SELECT
		id,
		user_id,
		api_key_name,
		exchange,
		api_key,
		api_secret,
		status
	FROM
		api_key
	WHERE
		id = $1
		AND user_id = $2
	order by
		id;`, id, userId)
	var apiKey model.ApiKey
	err := row.Scan(&apiKey.Id, &apiKey.UserId, &apiKey.ApiKeyName, &apiKey.Exchange, &apiKey.ApiKey, &apiKey.ApiSecret, &apiKey.Status)
	if err == sql.ErrNoRows {
		return nil, ErrApiKeyNotFound
	}

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (ur *UserPsql) SaveApiKey(apiKey *model.ApiKey) error {
	query := `UPDATE api_key
		SET api_key_name = $1, exchange = $2, api_key = $3, api_secret = $4, status = $5
		WHERE id = $6`

	res, err := ur.db.Exec(query, apiKey.ApiKeyName, apiKey.Exchange, apiKey.ApiKey, apiKey.ApiSecret, apiKey.Status, apiKey.Id)
	if err != nil {
		log.Println(err)
		return ErrCouldNotUpdateApikey
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("could not retrieve rows affected: %v", err)
		return ErrCouldNoteRetrieveAffectedRows
	}

	if rowsAffected == 0 {
		log.Printf("api_key with id %d not found", apiKey.Id)
		return ErrApiKeyNotFound
	}

	return nil
}

func (ur *UserPsql) ListUsersPermission(userId uint32) ([]model.Permission, error) {
	query := `
	SELECT p.permission_name
	FROM users AS u
	JOIN user_roles AS ur ON u.id = ur.user_id
	JOIN roles AS r ON ur.role_id = r.id
	JOIN role_permission AS rp ON r.id = rp.role_id
	JOIN permission AS p ON rp.permission_id = p.id
	WHERE u.id = $1 AND u.deleted_at IS NULL AND r.deleted_at IS NULL
	  AND rp.deleted_at IS NULL AND ur.deleted_at IS NULL;
	`

	rows, err := ur.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var permissionName model.Permission
		err := rows.Scan(&permissionName)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permissionName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
