package model

// type Transaction struct {
// 	Id   uint32
// 	Name string
// }

type Transaction string

const (
	ListUsers   Transaction = "ListUsers"
	ListApiKeys Transaction = "ListApiKeys"
)

type Role struct {
	Id           uint32
	Name         string
	Transactions []Transaction
}

func NewAdminRole() *Role {
	return &Role{
		Id:           0,
		Name:         "Admin",
		Transactions: []Transaction{ListApiKeys, ListUsers},
	}
}
