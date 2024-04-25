package interfaces

import "github.com/akhi9550/notification-svc/pkg/utils/models"

type NewauthClient interface {
	UserData(userID int) (models.UserData, error)
}
