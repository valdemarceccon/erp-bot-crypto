package schema

import "github.com/valdemarceccon/crypto-bot-erp/backend/model"

type UserResponse struct {
	Id       uint32 `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type PermissionResponse struct {
	Name string `json:"name"`
}

type UserMeResponse struct {
	UserResponse
	Permissions []PermissionResponse `json:"permissions"`
}

type ApiKeyResponse struct {
	Id           uint32             `json:"id"`
	UserId       uint32             `json:"user_id"`
	Username     string             `json:"username"`
	ApiKeyName   string             `json:"api_key_name"`
	Exchange     string             `json:"exchange"`
	ApiKeySecret string             `json:"api_secret"`
	ApiKey       string             `json:"api_key"`
	Status       model.ApiKeyStatus `json:"status"`
}

type ApiKeyRequest struct {
	ApiKeyName   string             `json:"api_key_name"`
	Exchange     string             `json:"exchange"`
	ApiKey       string             `json:"api_key"`
	ApiKeySecret string             `json:"api_secret"`
	Status       model.ApiKeyStatus `json:"status"`
}

func FromUserModel(user *model.User) *UserResponse {
	return &UserResponse{
		Id:       user.Id,
		Fullname: user.Fullname,
		Username: user.Username,
		Email:    user.Email,
	}
}

func FromApiKeyModel(apiKey *model.ApiKey) *ApiKeyResponse {
	return &ApiKeyResponse{
		Id:           apiKey.Id,
		UserId:       apiKey.UserId,
		Username:     apiKey.Username,
		ApiKeyName:   apiKey.ApiKeyName,
		ApiKeySecret: apiKey.ApiSecret,
		ApiKey:       apiKey.ApiKey,
		Exchange:     apiKey.Exchange,
		Status:       apiKey.Status,
	}
}
