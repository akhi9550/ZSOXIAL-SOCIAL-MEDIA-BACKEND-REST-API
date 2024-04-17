package interfaces

import (
	"github.com/akhi9550/chat-svc/pkg/utils/models"
)

type ChatRepository interface {
	// CreateChatRoom(user1, user2 uint) error
	// GetAllChats(userId uint) ([]models.Chat, error)
	// GetMessages(id primitive.ObjectID) ([]domain.Messages, error)
	// UpdateLastMessageAndTime(chatId primitive.ObjectID, lastMessage string, time time.Time) error
	// IsChatExist(user1, user2 uint) (bool, error)
	// IsValidChatId(chatId primitive.ObjectID) (bool, error)
	// SaveMessage(message domain.Messages) (primitive.ObjectID, error)
	// ReadMessage(chatId primitive.ObjectID, senderId uint) (int64, error)
	// FetchRecipient(chatId primitive.ObjectID, userId uint) (uint, error)
	// DeleteChatsAndMessagesByUserID(userID uint) error
	StoreFriendsChat(models.MessageReq) error
	GetLastMessage(string, string) (*models.Message, error)
	GetMessageCount(string, string) (int, error)
	GetFriendChat(string, string, models.Pagination) ([]models.Message, error)
	UpdateReadAsMessage(string, string) error
}
