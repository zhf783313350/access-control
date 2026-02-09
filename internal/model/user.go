package model

// User 用户表模型
type User struct {
	Id             int64  `db:"id" json:"id"`
	PhoneNumber    string `db:"phonenumber" json:"phoneNumber"`
	Status         int    `db:"status" json:"status"`
	ValidTime      string `db:"validtime" json:"validTime"`
	Organization   string `db:"organization" json:"organization"`
	Client         string `db:"client" json:"client"`
	OrganizationID string `db:"organization_id" json:"organizationId"`
	ClientID       string `db:"client_id" json:"clientId"`
}
