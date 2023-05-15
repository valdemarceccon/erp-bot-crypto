package store

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/store/query"
)

type apiKeyPsql struct {
	db *sql.DB
}

func NewApiKeyPsql(db *sql.DB) ApiKey {
	return &apiKeyPsql{
		db: db,
	}
}

func (api *apiKeyPsql) New(apiKey *model.ApiKey) error {
	return api.db.QueryRow(query.NewApiKey,
		apiKey.UserId,
		apiKey.ApiKeyName,
		apiKey.Exchange,
		apiKey.ApiKey,
		apiKey.ApiSecret,
		apiKey.Status).Scan(&apiKey.Id)
}

func (api *apiKeyPsql) Get(id, userId uint32) (*model.ApiKey, error) {
	row := api.db.QueryRow(query.GetApiKey,
		id,
		userId)
	var apiKey model.ApiKey
	err := row.Scan(&apiKey.Id,
		&apiKey.UserId,
		&apiKey.ApiKeyName,
		&apiKey.Exchange,
		&apiKey.ApiKey,
		&apiKey.ApiSecret,
		&apiKey.Status)
	if err == sql.ErrNoRows {
		return nil, ErrApiKeyNotFound
	}

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (api *apiKeyPsql) Save(apiKey *model.ApiKey) error {
	res, err := api.db.Exec(query.SaveApiKey,
		apiKey.ApiKeyName,
		apiKey.Exchange,
		apiKey.ApiKey,
		apiKey.ApiSecret,
		apiKey.Status,
		apiKey.Id)
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
func (api *apiKeyPsql) ListActive(userId uint32) ([]model.ApiKey, error) {
	rows, err := api.db.Query(query.ListActiveApiKey,
		userId,
		model.ApiKeyStatusActive,
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

func (api *apiKeyPsql) FromUser(userId uint32) ([]model.ApiKey, error) {
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

func (api *apiKeyPsql) List() ([]model.ApiKey, error) {
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

func (api *apiKeyPsql) GetBotRunsStartStop(userId uint32) ([]model.ApiKeyRun, error) {
	row, err := api.db.Query(query.BotStartStopApiKey, userId)
	if err != nil {
		log.Println(fmt.Errorf("api key store: %w", err))
		return nil, err
	}
	defer row.Close()
	resp := make([]model.ApiKeyRun, 0)

	for row.Next() {
		var run model.ApiKeyRun

		err = row.Scan(
			&run.Id,
			&run.UserId,
			&run.Username,
			&run.ApiKeyId,
			&run.StartTime,
			&run.StartBalance,
			&run.StopTime,
			&run.StopBalanace,
		)

		if err != nil {
			return nil, err
		}

		resp = append(resp, run)
	}

	return resp, nil
}
