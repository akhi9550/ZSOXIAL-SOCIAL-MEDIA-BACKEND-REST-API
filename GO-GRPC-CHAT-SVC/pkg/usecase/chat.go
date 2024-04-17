package usecase

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	authclientinterfaces "github.com/akhi9550/chat-svc/pkg/client/interface"
	"github.com/akhi9550/chat-svc/pkg/config"
	"github.com/akhi9550/chat-svc/pkg/helper"
	interfaces "github.com/akhi9550/chat-svc/pkg/repository/interface"
	services "github.com/akhi9550/chat-svc/pkg/usecase/interface"
	"github.com/akhi9550/chat-svc/pkg/utils/models"
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

func (c *ChatUseCase) MessageConsumer() {
	fmt.Println("Starting Kafka consumer")

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	configs := sarama.NewConfig()
	consumer, err := sarama.NewConsumer([]string{cfg.KafkaPort}, configs)
	if err != nil {
		fmt.Println("Error creating Kafka consumer:", err)
		return
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(cfg.KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("Error creating partition consumer:", err)
		return
	}
	defer partitionConsumer.Close()

	fmt.Println("Kafka consumer started")

	for {
		select {
		case message := <-partitionConsumer.Messages():
			fmt.Println(string(message.Value), "ddddd❌❌❌")
			msg, err := c.UnmarshelChatMessage(message.Value)
			if err != nil {
				fmt.Println("Error unmarshalling message:", err)
				continue
			}
			fmt.Println("Received message:", msg)
			err = c.chatRepository.StoreFriendsChat(*msg)
			if err != nil {
				fmt.Println("Error storing message in repository:", err)
				continue
			}
		case err := <-partitionConsumer.Errors():
			fmt.Println("Kafka consumer error:", err)
		}
	}
}

func (c *ChatUseCase) UnmarshelChatMessage(data []byte) (*models.MessageReq, error) {
	var message models.MessageReq
	err := json.Unmarshal(data, &message)
	if err != nil {
		return nil, err
	}
	message.Timestamp = time.Now()
	return &message, nil
}

func (c *ChatUseCase) GetFriendChat(userID, friendID string, pagination models.Pagination) ([]models.Message, error) {
	var err error
	pagination.OffSet, err = helper.Pagination(pagination.Limit, pagination.OffSet)
	if err != nil {
		return nil, err
	}
	_ = c.chatRepository.UpdateReadAsMessage(userID, friendID)
	return c.chatRepository.GetFriendChat(userID, friendID, pagination)
}

// func (c *ChatUseCase) GetAllChats(userId uint) ([]models.ChatResponse, error) {
// 	chats, err := c.chatRepository.GetAllChats(userId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	chatResponses := make([]models.ChatResponse, 0)
// 	for _, chat := range chats {
// 		userIDInt := int(chat.Users[0])
// 		user, err := c.authClient.UserData(userIDInt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		chatResponse := models.ChatResponse{
// 			Chat: chat,
// 			User: user,
// 		}
// 		chatResponses = append(chatResponses, chatResponse)
// 	}
// 	return chatResponses, nil
// }

// func (c *ChatUseCase) GetMessages(chatId primitive.ObjectID) ([]domain.Messages, error) {
// 	isValid, err := c.chatRepository.IsValidChatId(chatId)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if !isValid {
// 		return nil, errors.New("chatId is not existing")
// 	}
// 	messages, err := c.chatRepository.GetMessages(chatId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return messages, nil
// }

// func (c *ChatUseCase) SaveMessage(chatId primitive.ObjectID, senderId uint, message string) (primitive.ObjectID, error) {
// 	isValid, err := c.chatRepository.IsValidChatId(chatId)
// 	if err != nil {
// 		return primitive.ObjectID{}, err
// 	}
// 	if !isValid {
// 		return primitive.ObjectID{}, errors.New("chatId is not existing")
// 	}
// 	messages := domain.Messages{
// 		SenderID:       senderId,
// 		ChatID:         chatId,
// 		MessageContent: message,
// 		Timestamp:      time.Now(),
// 	}
// 	res, err := c.chatRepository.SaveMessage(messages)
// 	if err != nil {
// 		return primitive.ObjectID{}, err
// 	}
// 	err = c.chatRepository.UpdateLastMessageAndTime(chatId, message, messages.Timestamp)
// 	if err != nil {
// 		return primitive.ObjectID{}, err
// 	}
// 	return res, nil
// }

// func (c *ChatUseCase) ReadMessage(userId uint, chatId primitive.ObjectID) (int64, error) {

// 	isValid, err := c.chatRepository.IsValidChatId(chatId)
// 	if err != nil {
// 		return 0, err
// 	}
// 	if !isValid {
// 		return 0, errors.New("chatId is not existing")
// 	}
// 	senderId, err := c.chatRepository.FetchRecipient(chatId, userId)

// 	if err != nil {
// 		return 0, err
// 	}

// 	res, err := c.chatRepository.ReadMessage(chatId, senderId)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return res, nil
// }

// func (c *ChatUseCase) FetchRecipient(chatId primitive.ObjectID, userId uint) (uint, error) {
// 	isValid, err := c.chatRepository.IsValidChatId(chatId)
// 	if err != nil {
// 		return 0, err
// 	}
// 	if !isValid {
// 		return 0, errors.New("chatId is not existing")
// 	}
// 	res, err := c.chatRepository.FetchRecipient(chatId, userId)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return res, nil
// }

// func (c *ChatUseCase) CreateChatRoom(user1, user2 int64) error {
// 	create := c.chatRepository.CreateChatRoom(uint(user1), uint(user2))
// 	if create != nil {
// 		return create
// 	}
// 	return nil
// }
