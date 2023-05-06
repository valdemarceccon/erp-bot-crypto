package store

import (
	"database/sql"
	"log"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/store/query"
)

type ApiKeyPsql struct {
	db *sql.DB
}

func NewApiKeyPsql(db *sql.DB) ApiKey {
	return &ApiKeyPsql{
		db: db,
	}
}

func (api *ApiKeyPsql) New(apiKey *model.ApiKey) error {
	return api.db.QueryRow(query.NewApiKey,
		apiKey.UserId,
		apiKey.ApiKeyName,
		apiKey.Exchange,
		apiKey.ApiKey,
		apiKey.ApiSecret,
		apiKey.Status).Scan(&apiKey.Id)
}
func (api *ApiKeyPsql) Get(id, userId uint32) (*model.ApiKey, error) {
	row := api.db.QueryRow(`
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
func (api *ApiKeyPsql) Save(apiKey *model.ApiKey) error {
	query := `UPDATE api_key
	SET api_key_name = $1, exchange = $2, api_key = $3, api_secret = $4, status = $5
	WHERE id = $6`

	res, err := api.db.Exec(query, apiKey.ApiKeyName, apiKey.Exchange, apiKey.ApiKey, apiKey.ApiSecret, apiKey.Status, apiKey.Id)
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
func (api *ApiKeyPsql) ListActive(userId uint32) ([]model.ApiKey, error) {
	query := `
	SELECT
		a.id,
		a.user_id,
		a.api_key_name,
		a.exchange,
		a.api_key,
		a.api_secret,
		a.status
	FROM
		users u
		inner join api_key a on
					u.id = a.user_id
			and a.deleted_at is null
			and u.deleted_at is null
	WHERE
				(u.id = $1 OR $1 = 0)
		and a.status in ($2,$3);`

	rows, err := api.db.Query(query,
		userId, model.ApiKeyStatusActive,
		model.ApiKeyStatusWaitingDeactivation)
	if err != nil {
		return nil, err
	}

	ret := make([]model.ApiKey, 0)

	for rows.Next() {
		var apiKey model.ApiKey

		rows.Scan(
			&apiKey.Id,
			&apiKey.UserId,
			&apiKey.ApiKeyName,
			&apiKey.Exchange,
			&apiKey.ApiKey,
			&apiKey.ApiSecret,
			&apiKey.Status,
		)

		ret = append(ret, apiKey)
	}

	return ret, nil
}

func (api *ApiKeyPsql) FromUser(userId uint32) ([]model.ApiKey, error) {
	row, err := api.db.Query(query.GetApiFromUser, userId)
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

func (api *ApiKeyPsql) List() ([]model.ApiKey, error) {
	row, err := api.db.Query(query.ListApiKeys)
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

func (api *ApiKeyPsql) ListUsersPermission(userId uint32) ([]model.Permission, error) {
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

	rows, err := api.db.Query(query, userId)
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
