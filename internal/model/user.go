package model

// User 用户表模型
type User struct {
	Id          int64  `db:"id" json:"id"`
	PhoneNumber string `db:"phoneNumber" json:"phoneNumber"`
	Status      int    `db:"status" json:"status"`
	ValidTime   string `db:"validTime" json:"validTime"`
}
