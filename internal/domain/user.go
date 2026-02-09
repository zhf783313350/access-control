package domain

import (
	"errors"
	"strings"
)

// User 领域模型
type User struct {
	ID          int64
	PhoneNumber string
	Status      UserStatus
	ValidTime   string
}

type UserStatus int

const (
	UserStatusNormal  UserStatus = 1
	UserStatusBlocked UserStatus = 2
)

// Validate 领域验证逻辑
func (u *User) Validate() error {
	if strings.TrimSpace(u.PhoneNumber) == "" {
		return errors.New("phone number cannot be empty")
	}
	// 这里可以添加更复杂的校验逻辑，如正则匹配
	return nil
}

// IsExpired 检查用户是否过期 (示例逻辑)
func (u *User) IsExpired() bool {
	// 实现基于 ValidTime 的判断
	return false
}
