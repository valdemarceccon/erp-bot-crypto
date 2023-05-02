package model

// type Transaction struct {
// 	Id   uint32
// 	Name string
// }

type Permission string

const (
	ListUsersPermission   Permission = "ListUsers"
	ListApiKeysPermission Permission = "ListApiKeys"
)

type Role struct {
	Id           uint32
	Name         string
	Transactions []Permission
}