package repository

import (
	"database/sql"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
)

type RolePsql struct {
	db *sql.DB
}

func NewRolePsql(db *sql.DB) Role {
	return &RolePsql{
		db: db,
	}
}

func (r *RolePsql) Create(user *model.Role) error {

	return ErrNotImplemented
}

func (r *RolePsql) Get(id uint32) (*model.Role, error) {

	return nil, ErrNotImplemented
}

func (r *RolePsql) GetAll() ([]model.Role, error) {

	return nil, ErrNotImplemented
}

func (r *RolePsql) Update(user *model.Role) error {

	return ErrNotImplemented
}
func (r *RolePsql) Delete(id uint32) error {

	return ErrNotImplemented
}
func (r *RolePsql) SearchByName(string) (*model.Role, error) {

	return nil, ErrNotImplemented
}

func (r *RolePsql) UserPermissions(userId uint32) ([]model.Permission, error) {
	// 	CREATE TABLE user_roles (
	//     user_id INTEGER NOT NULL REFERENCES users(id),
	//     role_id INTEGER NOT NULL REFERENCES roles(id),
	//     PRIMARY KEY (user_id, role_id),
	//     created_at timestamp not null,
	//     updated_at timestamp not null,
	//     deleted_at timestamp null
	// );
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
