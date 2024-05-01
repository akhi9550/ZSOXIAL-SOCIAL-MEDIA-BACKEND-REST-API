package models

type NotificationPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type NotificationResponse struct {
	UserID    int    `json:"user_id" gorm:"column:sender_id"`
	Username  string `json:"username"`
	Profile   string `json:"profile"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
}
