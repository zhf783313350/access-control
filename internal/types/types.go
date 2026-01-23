package types

import "accesscontrol/internal/model"

// 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 登录请求
type LoginRequest struct {
	PhoneNumber string `json:"phoneNumber"`
}

// 登录响应
type LoginResponse struct {
	AccessToken  string     `json:"accessToken"`
	AccessExpire int64      `json:"accessExpire"`
	UserInfo     model.User `json:"userInfo"`
}

// 删除用户请求
type DeleteUserRequest struct {
	PhoneNumber string `json:"phoneNumber"`
}

// 用户列表请求
type ListUsersRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}
