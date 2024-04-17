package interfaces

import (
	"github.com/akhi9550/chat-svc/pkg/utils/models"
)

type ChatUseCase interface {
	// GetAllChats(userId uint) ([]models.ChatResponse, error)
	// GetMessages(chatId primitive.ObjectID) ([]domain.Messages, error)
	// SaveMessage(chatId primitive.ObjectID, senderId uint, message string) (primitive.ObjectID, error)
	// ReadMessage(userId uint, chatId primitive.ObjectID) (int64, error)
	// FetchRecipient(chatId primitive.ObjectID, userId uint) (uint, error)
	// CreateChatRoom(user1, user2 int64) error
	MessageConsumer()
	GetFriendChat(string, string, models.Pagination) ([]models.Message, error)
}
