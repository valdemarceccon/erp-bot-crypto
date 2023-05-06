package controller

import "github.com/valdemarceccon/crypto-bot-erp/backend/store"

type DataCollector struct {
	store.All
}

func NewDataCollector(userStore store.User, roleStore store.Role) *DataCollector {
	return &DataCollector{
		store.All{
			User: userStore,
			Role: roleStore,
		},
	}
}
