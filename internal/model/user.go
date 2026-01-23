package model

// User 用户表模型
type User struct {
	Id          int64  `db:"id" json:"id"`
	Password    string `db:"password" json:"-"`
	PhoneNumber string `db:"phoneNumber" json:"phoneNumber"`
	Status      int    `db:"status" json:"status"`
	ValidTime   string `db:"validTime" json:"validTime"`
	CreatedAt   string `db:"created_at" json:"createdAt"`
	UpdatedAt   string `db:"updated_at" json:"updatedAt"`
}
