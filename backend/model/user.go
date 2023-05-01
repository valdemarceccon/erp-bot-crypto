package model

type User struct {
	Id       uint32 `json:"id"`
	Fullname string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Telegram string `json:"telegram"`
}
