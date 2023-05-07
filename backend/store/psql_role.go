package store

import (
	"database/sql"
	"log"

	"github.com/valdemarceccon/crypto-bot-erp/backend/model"
	"github.com/valdemarceccon/crypto-bot-erp/backend/store/query"
)

type RolePsql struct {
	db *sql.DB
}

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
	rows, err := r.db.Query(query.PermissionFromUser, userId)

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

// func (api *RolePsql) ListUsersPermission(userId uint32) ([]model.Permission, error) {
// 	query := `
// 	SELECT p.permission_name
// 	FROM users AS u
// 	JOIN user_roles AS ur ON u.id = ur.user_id
// 	JOIN roles AS r ON ur.role_id = r.id
// 	JOIN role_permission AS rp ON r.id = rp.role_id
// 	JOIN permission AS p ON rp.permission_id = p.id
// 	WHERE u.id = $1 AND u.deleted_at IS NULL AND r.deleted_at IS NULL
// 	  AND rp.deleted_at IS NULL AND ur.deleted_at IS NULL;
// 	`

// 	rows, err := api.db.Query(query, userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var permissions []model.Permission
// 	for rows.Next() {
// 		var permissionName model.Permission
// 		err := rows.Scan(&permissionName)
// 		if err != nil {
// 			return nil, err
// 		}
// 		permissions = append(permissions, permissionName)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return permissions, nil
// }
