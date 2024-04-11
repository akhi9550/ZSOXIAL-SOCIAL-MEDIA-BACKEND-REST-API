package client

import (
	"context"
	"fmt"
	"time"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/config"
	pb "github.com/akhi9550/api-gateway/pkg/pb/chat"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"google.golang.org/grpc"
)

type ChatClient struct {
	Client pb.ChatServiceClient
}

func NewChatClient(cfg config.Config) interfaces.ChatClient {
	grpcConnection, err := grpc.Dial(cfg.ChatSvcUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect", err)
	}

	grpcClient := pb.NewChatServiceClient(grpcConnection)

	return &ChatClient{
		Client: grpcClient,
	}
}

func (ch *ChatClient) GetAllChats(userId uint) ([]models.ChatResponse, error) {
	data, err := ch.Client.GetAllChats(context.Background(), &pb.GetAllChatsRequest{
		Userid: int64(userId),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get chats")
	}

	var chatResponses []models.ChatResponse

	for _, chat := range data.Response {
		var users []models.UserData
		for _, user := range chat.User {
			users = append(users, models.UserData{
				UserId:   uint(user.Userid),
				Username: user.Username,
				Profile:  user.Profile,
			})
		}

		lastMessageTime := time.Unix(chat.ChatSeconds, int64(chat.LastMessageTime.Nanos))

		chatResponse := models.ChatResponse{
			Chat: models.Chat{
				ID:              primitive.ObjectID(chat.ID),
				Users:           chat.User,
				LastMessage:     chat.lat,
				LastMessageTime: lastMessageTime,
			},
			User: users,
		}

		chatResponses = append(chatResponses, chatResponse)
	}

	return chatResponses, nil
}

func (ch *ChatClient) GetMessages(chatId primitive.ObjectID) ([]models.Messages, error) {
	data, err := ch.Client.GetMessages(context.Background(), &pb.GetMessagesRequest{
		Chatid: chatId.Hex(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get messages")
	}
	var messages []models.Messages

	for _, msg := range data.Response {
		timestamp := time.Unix(msg.Timestamp.Seconds, int64(msg.Timestamp.Nanos))
		id, err := primitive.ObjectIDFromHex(msg.Id)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert ID to ObjectID")
		}

		message := models.Messages{
			ID:             id,
			SenderID:       uint(msg.SenderId),
			ChatID:         chatId,
			Seen:           msg.Seen,
			Image:          msg.Image,
			MessageContent: msg.MessageContent,
			Timestamp:      timestamp,
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (ch *ChatClient) SaveMessage(chatId primitive.ObjectID, senderId uint, message string) (primitive.ObjectID, error) {
	data, err := ch.Client.SaveMessage(context.Background(), &pb.SaveMessageRequest{
		Chatid:   chatId.Hex(),
		Senderid: int64(senderId),
		Message:  message,
	})
	if err != nil {
		return primitive.NilObjectID, errors.Wrap(err, "failed to save message")
	}
	objectID, err := primitive.ObjectIDFromHex(data.Id)
	if err != nil {
		return primitive.NilObjectID, errors.Wrap(err, "failed to convert ID from hex")
	}
	return objectID, nil
}

func (ch *ChatClient) ReadMessage(userId uint, chatId primitive.ObjectID) (int64, error) {
	data, err := ch.Client.ReadMessage(context.Background(), &pb.ReadMessageRequest{
		Chatid: chatId.Hex(),
		Userid: int64(userId),
	})
	if err != nil {
		return 0, err
	}
	return data.Id, nil
}
func (ch *ChatClient) FetchRecipient(chatId primitive.ObjectID, userId uint) (uint, error) {
	data, err := ch.Client.FetchRecipient(context.Background(), &pb.FetchRecipientRequest{
		Chatid: chatId.Hex(),
		Userid: int64(userId),
	})
	if err != nil {
		return 0, err
	}
	return uint(data.Id), nil
}
