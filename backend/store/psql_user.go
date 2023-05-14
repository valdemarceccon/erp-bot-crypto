package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/big"

	"github.com/hirokisan/bybit/v2"
	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/store/query"
)

type UserPsql struct {
	db *sql.DB
}

func NewUserPsql(db *sql.DB) User {
	return &UserPsql{
		db: db,
	}
}

func (ur *UserPsql) New(user *model.User) error {
	queryExists, err := ur.db.Query(query.UserExistsUsernameEmail,
		user.Username,
		user.Email)
	if err != nil {
		return fmt.Errorf("user: %w", err)
	}
	defer queryExists.Close()

	if queryExists.Next() {
		return ErrUserOrEmailInUse
	}

	row := ur.db.QueryRow(query.NewUser, user.Email, user.Username, user.Fullname, user.Password)

	return row.Scan(&user.Id)
}

func (ur *UserPsql) Get(id uint32) (*model.User, error) {
	row := ur.db.QueryRow(query.GetUser, id)

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

func (ur *UserPsql) List() ([]model.User, error) {
	rows, err := ur.db.Query(query.ListUser)

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

func (ur *UserPsql) Save(user *model.User) error {

	return ErrNotImplemented
}

func (ur *UserPsql) Delete(id uint32) error {

	return ErrNotImplemented
}

func (ur *UserPsql) ByUsername(username string) (*model.User, error) {
	row := ur.db.QueryRow(query.GetUsernameUser, username)
	resp := &model.User{}
	err := row.Scan(&resp.Id, &resp.Email, &resp.Username, &resp.Fullname, &resp.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("user: %w", err)
	}
	return resp, nil
}

func (r *UserPsql) SaveClosedPnL(userId, apiKeyId uint32, data []bybit.V5GetClosedPnLItem) error {
	tx, err := r.db.Begin()

	if err != nil {
		return fmt.Errorf("failed to insert closed pnl data: %w", err)
	}

	stmt, err := tx.Prepare(query.SaveClosedPnLUser)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, val := range data {
		_, err := stmt.Exec(
			userId,
			apiKeyId,
			val.Symbol,
			val.OrderID,
			val.ExecType,
			val.ClosedPnl,
			val.CreatedTime,
			val.UpdatedTime,
		)

		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert closed pnl data: %w", err)
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit closed pnl data: %w", err)
	}

	return nil
}

func (r *UserPsql) GetClosedPnL(userId, apiKeyId uint32, startTime, endTime int64) ([]bybit.V5GetClosedPnLItem, error) {
	rows, err := r.db.Query(query.GetClosedPnLUser, userId, apiKeyId, startTime, endTime)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := make([]bybit.V5GetClosedPnLItem, 0)
	for rows.Next() {
		var val bybit.V5GetClosedPnLItem
		err = rows.Scan(
			&val.Symbol,
			&val.OrderID,
			&val.ExecType,
			&val.ClosedPnl,
			&val.CreatedTime,
			&val.UpdatedTime,
		)
		if err != nil {
			return nil, err
		}
		ret = append(ret, val)
	}

	return ret, nil
}

func (r *UserPsql) StartBot(apikey *model.ApiKey, balance *big.Float) error {
	_, err := r.db.Exec(query.StarBotUser,
		apikey.Id,
		balance)

	return err
}

func (r *UserPsql) StopBot(apikey *model.ApiKey, balance *big.Float) error {
	res, err := r.db.Exec(query.StopBotUser,
		balance,
		apikey.Id)

	if err != nil {
		return err
	}

	newRows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if newRows != 1 {
		return fmt.Errorf("stop bot: something went wrong, it should have inserted just a row, but %d was inserted", newRows)
	}

	return nil
}
