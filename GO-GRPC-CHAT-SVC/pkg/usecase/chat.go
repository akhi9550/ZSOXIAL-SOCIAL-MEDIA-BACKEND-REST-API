package usecase

import (
	"errors"
	"time"

	authclientinterfaces "github.com/akhi9550/chat-svc/pkg/client/interface"
	"github.com/akhi9550/chat-svc/pkg/domain"
	interfaces "github.com/akhi9550/chat-svc/pkg/repository/interface"
	services "github.com/akhi9550/chat-svc/pkg/usecase/interface"
	"github.com/akhi9550/chat-svc/pkg/utils/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatUseCase struct {
	chatRepository interfaces.ChatRepository
	authClient     authclientinterfaces.NewauthClient
}

func NewChatUseCase(repository interfaces.ChatRepository, authclient authclientinterfaces.NewauthClient) services.ChatUseCase {
	return &ChatUseCase{
		chatRepository: repository,
		authClient:     authclient,
	}
}

func (c *ChatUseCase) GetAllChats(userId uint) ([]models.ChatResponse, error) {
	chats, err := c.chatRepository.GetAllChats(userId)
	if err != nil {
		return nil, err
	}
	chatResponses := make([]models.ChatResponse, 0)
	for _, chat := range chats {
		userIDInt := int(chat.Users[0])
		user, err := c.authClient.UserData(userIDInt)
		if err != nil {
			return nil, err
		}
		chatResponse := models.ChatResponse{
			Chat: chat,
			User: user,
		}
		chatResponses = append(chatResponses, chatResponse)
	}
	return chatResponses, nil
}

func (c *ChatUseCase) GetMessages(chatId primitive.ObjectID) ([]domain.Messages, error) {
	isValid, err := c.chatRepository.IsValidChatId(chatId)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, errors.New("chatId is not existing")
	}
	messages, err := c.chatRepository.GetMessages(chatId)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *ChatUseCase) SaveMessage(chatId primitive.ObjectID, senderId uint, message string) (primitive.ObjectID, error) {
	isValid, err := c.chatRepository.IsValidChatId(chatId)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	if !isValid {
		return primitive.ObjectID{}, errors.New("chatId is not existing")
	}
	messages := domain.Messages{
		SenderID:       senderId,
		ChatID:         chatId,
		MessageContent: message,
		Timestamp:      time.Now(),
	}
	res, err := c.chatRepository.SaveMessage(messages)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	err = c.chatRepository.UpdateLastMessageAndTime(chatId, message, messages.Timestamp)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res, nil
}

func (c *ChatUseCase) ReadMessage(userId uint, chatId primitive.ObjectID) (int64, error) {

	isValid, err := c.chatRepository.IsValidChatId(chatId)
	if err != nil {
		return 0, err
	}
	if !isValid {
		return 0, errors.New("chatId is not existing")
	}
	senderId, err := c.chatRepository.FetchRecipient(chatId, userId)

	if err != nil {
		return 0, err
	}

	res, err := c.chatRepository.ReadMessage(chatId, senderId)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func (c *ChatUseCase) FetchRecipient(chatId primitive.ObjectID, userId uint) (uint, error) {
	isValid, err := c.chatRepository.IsValidChatId(chatId)
	if err != nil {
		return 0, err
	}
	if !isValid {
		return 0, errors.New("chatId is not existing")
	}
	res, err := c.chatRepository.FetchRecipient(chatId, userId)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (c *ChatUseCase) CreateChatRoom(user1, user2 int64) error {
	create := c.chatRepository.CreateChatRoom(uint(user1), uint(user2))
	if create != nil {
		return create
	}
	return nil
}
