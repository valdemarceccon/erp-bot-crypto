package store

import (
	"database/sql"
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

func (ur *UserPsql) List() ([]model.User, error) {
	rows, err := ur.db.Query(query.ListUsers)

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

func (r *UserPsql) SaveClosedPnL(userId, apiKeyId uint32, data []bybit.V5GetClosedPnLItem) error {

	query := `
	INSERT INTO closed_pnl (
		user_id, api_key_id, symbol, orderId, side, qty, orderPrice,
		orderType, execType, closedSize, cumEntryValue, avgEntryPrice,
		cumExitValue, avgExitPrice, closedPnl, fillCount, leverage,
		createdTime, updatedTime, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12,
		$13, $14, $15, $16, $17, $18, $19, now(), now()
	);`
	tx, err := r.db.Begin()

	if err != nil {
		return fmt.Errorf("failed to insert closed pnl data: %w", err)
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, val := range data {
		_, err := stmt.Exec(
			userId, apiKeyId, val.Symbol, val.OrderID, val.Side, val.Qty, val.OrderPrice,
			val.OrderType, val.ExecType, val.ClosedSize, val.CumEntryValue, val.AvgEntryPrice,
			val.CumExitValue, val.AvgExitPrice, val.ClosedPnl, val.FillCount, val.Leverage,
			val.CreatedTime, val.UpdatedTime,
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

func (r *UserPsql) StartBot(apikey *model.ApiKey, balance *big.Float) error {
	query := `
		INSERT INTO bot_start(
			api_key_id,
			start_time,
			wallet_balance
		)
		values (
			$1, now(), $2
		)
	`
	_, err := r.db.Exec(query, apikey.Id, balance)

	return err
}

func (r *UserPsql) StopBot(apikey *model.ApiKey, balance *big.Float) error {

	query := `
	INSERT INTO bot_stop (
		stop_time,
		start_time_id,
		wallet_balance
	) SELECT
		NOW(),
		start.id,
		$1
	FROM
		bot_start start
	WHERE
		start.api_key_id = $2
		AND NOT EXISTS (
			select * from bot_stop stop where stop.start_time_id = start.id
		);
`

	res, err := r.db.Exec(query, balance, apikey.Id)

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
