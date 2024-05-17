package domain

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"uniquekey; not null"`
	Firstname string    `json:"firstname" gorm:"validate:required"`
	Lastname  string    `json:"lastname" gorm:"validate:required"`
	Username  string    `json:"username" gorm:"validate:required"`
	Dob       string    `json:"dob" gorm:"validate:required"`
	Gender    string    `json:"gender" gorm:"validate:required"`
	Phone     string    `json:"phone" gorm:"validate:required"`
	Email     string    `json:"email" validate:"email"`
	Password  string    `json:"password" validate:"min=6,max=20"`
	Bio       string    `json:"bio"`
	Imageurl  string    `json:"imageurl" gorm:"validate:required"`
	CreatedAt time.Time `json:"created_at"`
	Blocked   bool      `json:"blocked" gorm:"default:false"`
	Isadmin   bool      `json:"is_admin" gorm:"default:false"`
}

type UserReports struct {
	ID           uint   `json:"id" gorm:"uniquekey; not null"`
	ReportUserID uint   `json:"report_user_id"`
	UserID       uint   `json:"user_id"`
	Report       string `json:"reports"`
}

type FollowingRequests struct {
	UserID        uint      `json:"user_id"`
	FollowingUser uint      `json:"following_user"`
	CreatedAt     time.Time `json:"created_at"`
}

type Followings struct {
	UserID        uint      `json:"user_id"`
	FollowingUser uint      `json:"following_user"`
	CreatedAt     time.Time `json:"created_at"`
}

type Followers struct {
	UserID        uint      `json:"user_id"`
	FollowingUser uint      `json:"following_user"`
	CreatedAt     time.Time `json:"created_at"`
}