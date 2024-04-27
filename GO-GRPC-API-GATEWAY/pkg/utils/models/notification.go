package models

type NotificationPagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
