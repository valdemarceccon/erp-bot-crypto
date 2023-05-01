package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
	"github.com/valdemarceccon/crypto-bot-erp/backend/repository"
)

func init() {
	goose.AddMigration(upInitialData, downInitialData)
}

var allTransactions []repository.Transactions = []repository.Transactions{
	repository.ListUsers,
}

func upInitialData(tx *sql.Tx) error {

	row := tx.QueryRow("INSERT INTO roles(name) VALUES ('Admin') RETURNING id;")
	var adminId uint32
	err := row.Scan(&adminId)
	if err != nil {
		return err
	}
	for _, v := range allTransactions {
		row = tx.QueryRow("INSERT INTO permissions (name) VALUES ($1) RETURNING id;", v)
		var transctionId uint32
		err = row.Scan(&transctionId)
		if err != nil {
			return err
		}

		if _, err := tx.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)", adminId, transctionId); err != nil {
			return err
		}
	}

	return nil
}

func downInitialData(tx *sql.Tx) error {
	tx.Exec("DELETE FROM role_permissions;")
	tx.Exec("DELETE FROM permissions;")
	tx.Exec("DELETE FROM roles;")

	return nil
}
