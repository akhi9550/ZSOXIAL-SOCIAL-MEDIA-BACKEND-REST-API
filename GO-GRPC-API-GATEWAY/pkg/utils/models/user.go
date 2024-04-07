package models

type UserSignUpRequest struct {
	Firstname string `json:"firstname" validate:"gte=3"`
	Lastname  string `json:"lastname" validate:"gte=1"`
	Username  string `json:"username" validate:"gte=3"`
	Phone     string `json:"phone" validate:"min=10,max=10"`
	Email     string `json:"email" validate:"email"`
	Password  string `json:"password" validate:"min=6,max=20"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Imageurl string `json:"imageurl"`
	Isadmin  bool   `json:"is_admin"`
}

type UserResponsewithPassword struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Imageurl string `json:"imageurl"`
	Isadmin  bool   `json:"is_admin"`
}

type ReponseWithToken struct {
	Users        UserResponse
	AccessToken  string
	RefreshToken string
}

type OTPData struct {
	PhoneNumber string `json:"phone" `
}

type VerifyData struct {
	User *OTPData `json:"user" binding:"required" validate:"required"`
	Code string   `json:"code" binding:"required" validate:"required"`
}

type ForgotPasswordSend struct {
	Phone string `json:"phone"`
}

type ForgotVerify struct {
	Phone       string `json:"phone" binding:"required" validate:"required"`
	Otp         string `json:"otp" binding:"required"`
	NewPassword string `json:"newpassword" binding:"required" validate:"min=6,max=20"`
}

type UpdatePassword struct {
	OldPassword        string `json:"old_password" binding:"required"`
	NewPassword        string `json:"new_password" binding:"required"`
	ConfirmNewPassword string `json:"confirm_new_password" binding:"required"`
}

type UsersProfileDetail struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Dob       string `json:"dob"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
}

type UsersProfileDetails struct {
	Firstname string `json:"firstname" validate:"gte=3"`
	Lastname  string `json:"lastname" validate:"gte=1"`
	Username  string `json:"username" validate:"gte=3"`
	Dob       string `json:"dob" gorm:"validate:required"`
	Gender    string `json:"gender" gorm:"validate:required"`
	Phone     string `json:"phone" validate:"min=10,max=10"`
	Email     string `json:"email" validate:"email"`
	Bio       string `json:"bio"`
	Imageurl  string `json:"imageurl" gorm:"validate:required"`
}

type ChangePassword struct {
	Oldpassword string `json:"old_password"`
	Password    string `json:"password"`
	Repassword  string `json:"re_password"`
}

type ReportRequest struct {
	UserID uint   `json:"user_id"`
	Report string `json:"report"`
}