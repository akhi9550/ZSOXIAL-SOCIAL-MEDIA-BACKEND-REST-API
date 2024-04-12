package interfaces

import (
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatClient interface {
	GetAllChats(userId uint) ([]models.ChatResponse, error)
	GetMessages(chatId primitive.ObjectID) ([]models.Messages, error)
	SaveMessage(chatId primitive.ObjectID, senderId uint, message string) (primitive.ObjectID, error)
	ReadMessage(userId uint, chatId primitive.ObjectID) (int64, error)
	FetchRecipient(chatId primitive.ObjectID, userId uint) (uint, error)
	
}
