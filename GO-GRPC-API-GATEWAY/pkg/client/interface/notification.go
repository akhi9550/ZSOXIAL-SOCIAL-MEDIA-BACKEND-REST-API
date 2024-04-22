package interfaces

import "github.com/akhi9550/api-gateway/pkg/utils/models"

type NotificationClient interface {
	SendLikeNotification(userID, postID int) (string, error)
	ConsumeKafkaLikeMessages(userID int) (models.Responses,error)
	ConsumeKafkaCommentMessages(userID int)(models.Responses,error)
}
