package store

import (
	"database/sql"
	"log"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
)

type RolePsql struct {
	db *sql.DB
}

// New(user *model.Role) error
// Get(id uint32) (*model.Role, error)
// List() ([]model.Role, error)
// Save(user *model.Role) error
// Delete(id uint32) error
// ByName(string) (*model.Role, error)
// FromUser(userId uint32) ([]model.Permission, error)

func NewRolePsql(db *sql.DB) Role {
	return &RolePsql{
		db: db,
	}
}

func (r *RolePsql) New(user *model.Role) error {

	return ErrNotImplemented
}

func (r *RolePsql) Get(id uint32) (*model.Role, error) {

	return nil, ErrNotImplemented
}

func (r *RolePsql) List() ([]model.Role, error) {

	return nil, ErrNotImplemented
}

func (r *RolePsql) Save(user *model.Role) error {

	return ErrNotImplemented
}
func (r *RolePsql) Delete(id uint32) error {

	return ErrNotImplemented
}
func (r *RolePsql) ByName(string) (*model.Role, error) {

	return nil, ErrNotImplemented
}

func (r *RolePsql) FromUser(userId uint32) ([]model.Permission, error) {
	// TODO: deleted_at case
	rows, err := r.db.Query(`
		SELECT
			p.permission_name
		FROM
			user_roles ur
			inner join role_permission rp on
				ur.role_id = rp.role_id
			inner join permission p on
				rp.permission_id = p.id
		WHERE
			ur.user_id = $1
	`, userId)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	ret := make([]model.Permission, 0)

	for rows.Next() {
		var permissionName string
		rows.Scan(&permissionName)
		ret = append(ret, model.Permission(permissionName))
	}

	return ret, nil
}
