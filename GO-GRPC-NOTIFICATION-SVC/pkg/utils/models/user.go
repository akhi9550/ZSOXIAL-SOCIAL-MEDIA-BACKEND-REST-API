package models

type UserData struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	Profile  string `json:"profile"`
}
