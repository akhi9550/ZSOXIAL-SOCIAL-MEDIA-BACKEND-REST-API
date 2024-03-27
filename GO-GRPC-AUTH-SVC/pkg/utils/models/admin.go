package models

import "time"

type AdminLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type AdminResponse struct {
	Id       uint   `json:"id"`
	Email    string `json:"email"`
	Imageurl string `json:"imageurl"`
	Isadmin  bool   `json:"is_admin"`
}
type AdminResponsewithPassword struct {
	Id       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Imageurl string `json:"imageurl"`
	Isadmin  bool   `json:"is_admin"`
}
type AdminReponseWithToken struct {
	Users        AdminResponse
	AccessToken  string
	RefreshToken string
}
type UserDetailsAtAdmin struct {
	Id        uint      `json:"id"`
	Firstname string    `json:"firstname" validate:"gte=3"`
	Lastname  string    `json:"lastname" validate:"gte=1"`
	Username  string    `json:"username" validate:"gte=3"`
	Dob       string    `json:"dob" gorm:"validate:required"`
	Gender    string    `json:"gender" gorm:"validate:required"`
	Phone     string    `json:"phone" validate:"e164"`
	Email     string    `json:"email" validate:"email"`
	Imageurl  string    `json:"imageurl" gorm:"validate:required"`
	CreatedAt time.Time `json:"created_at"`
	Blocked   bool      `json:"blocked" gorm:"default:false"`
}
