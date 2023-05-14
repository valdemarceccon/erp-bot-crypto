package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	Id       uint32 `json:"id"`
	Fullname string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Telegram string `json:"telegram"`
}

type ApiKeyStatus uint8

const (
	ApiKeyStatusInactive ApiKeyStatus = iota
	ApiKeyStatusWaitingActivation
	ApiKeyStatusActive
	ApiKeyStatusWaitingDeactivation
)

type ApiKey struct {
	Id         uint32       `json:"id"`
	UserId     uint32       `json:"user_id"`
	Username   string       `json:"username"`
	ApiKeyName string       `json:"api_key_name"`
	Exchange   string       `json:"exchange"`
	ApiKey     string       `json:"api_key"`
	ApiSecret  string       `json:"api_secret"`
	Status     ApiKeyStatus `json:"status"`
}

type ApiKeyRun struct {
	Id           uint32           `json:"id"`
	ApiKeyId     uint32           `json:"api_key_id"`
	StartTime    *time.Time       `json:"start_time"`
	StopTime     *time.Time       `json:"stop_time"`
	StartBalance *decimal.Decimal `json:"start_balance"`
	StopBalanace *decimal.Decimal `json:"stop_balanace"`
}
