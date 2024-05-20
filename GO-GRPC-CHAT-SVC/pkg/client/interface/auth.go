package interfaces

import "github.com/akhi9550/chat-svc/pkg/utils/models"

type NewauthClient interface {
	CheckUserAvalilabilityWithUserID(userID int) bool
	UserData(userID int) (models.UserData, error)
}