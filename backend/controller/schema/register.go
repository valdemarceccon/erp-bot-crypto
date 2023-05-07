package schema

type RegisterRequest struct {
	Username string `validate:"required" json:"username"`
	Email    string `validate:"required" json:"email"`
	Password string `validate:"required" json:"password"`
	Telegram string `json:"telegram"`
	Fullname string `validate:"required" json:"fullname"`
}
