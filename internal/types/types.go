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
	Password    string `json:"password"`
}

// 登录响应
type LoginResponse struct {
	AccessToken  string   `json:"accessToken"`
	AccessExpire int64    `json:"accessExpire"`
	UserInfo     model.User `json:"userInfo"`
}
