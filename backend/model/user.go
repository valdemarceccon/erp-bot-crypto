package model

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
	Inactive ApiKeyStatus = iota
	WaitingActivation
	Active
	WaitingDeactivation
)

type ApiKey struct {
	Id         uint32       `json:"Id"`
	UserId     uint32       `json:"user_id"`
	ApiKeyName string       `json:"api_key_name"`
	Exchange   string       `json:"exchange"`
	ApiKey     string       `json:"api_key"`
	ApiSecret  string       `json:"api_secret"`
	Status     ApiKeyStatus `json:"status"`
}