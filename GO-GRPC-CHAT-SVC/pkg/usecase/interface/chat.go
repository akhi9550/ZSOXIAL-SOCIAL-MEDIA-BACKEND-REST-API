package interfaces

import (
	"github.com/akhi9550/chat-svc/pkg/utils/models"
)

type ChatUseCase interface {
	MessageConsumer()
	GetFriendChat(string, string, models.Pagination) ([]models.Message, error)
	GetGroupMessages(groupID string, limit, offset int) ([]models.Message, error)
}
