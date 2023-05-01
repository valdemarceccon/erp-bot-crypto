package schema

import "github.com/valdemarceccon/crypto-bot-erp/backend/model"

type UserResponse struct {
	Id       uint32
	Name     string
	Username string
	Email    string
}

func FromUserModel(user *model.User) *UserResponse {
	return &UserResponse{
		Id:       user.Id,
		Name:     user.Fullname,
		Username: user.Username,
		Email:    user.Email,
	}
}
