package interfaces

import (
	"time"

	"github.com/akhi9550/chat-svc/pkg/domain"
	"github.com/akhi9550/chat-svc/pkg/utils/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatRepository interface {
	CreateChatRoom(user1, user2 uint) error
	GetAllChats(userId uint) ([]models.Chat, error)
	GetMessages(id primitive.ObjectID) ([]domain.Messages, error)
	UpdateLastMessageAndTime(chatId primitive.ObjectID, lastMessage string, time time.Time) error
	IsChatExist(user1, user2 uint) (bool, error)
	IsValidChatId(chatId primitive.ObjectID) (bool, error)
	SaveMessage(message domain.Messages) (primitive.ObjectID, error)
	ReadMessage(chatId primitive.ObjectID, senderId uint) (int64, error)
	FetchRecipient(chatId primitive.ObjectID, userId uint) (uint, error)
	DeleteChatsAndMessagesByUserID(userID uint) error
}
