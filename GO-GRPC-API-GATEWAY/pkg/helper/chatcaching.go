package helper

import (
	"context"
	"encoding/json"
	"time"

	pb "github.com/akhi9550/api-gateway/pkg/pb/chat"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/redis/go-redis/v9"
)

type RedisChatCaching struct {
	redis      *redis.Client
	chatClient pb.ChatServiceClient
}

func NewRedisChatCaching(redis *redis.Client, chatClient pb.ChatServiceClient) *RedisChatCaching {
	return &RedisChatCaching{
		redis:      redis,
		chatClient: chatClient,
	}
}

func (r *RedisChatCaching) structMarshel(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (r *RedisChatCaching) jsonUnmarshel(model interface{}, data []byte) error {
	return json.Unmarshal(data, model)
}

func (r *RedisChatCaching) GetChat(userID string, req models.ChatRequest) ([]models.Message, error) {
	res := r.redis.Get(context.Background(), "chat:"+userID)
	var data []models.Message

	if res.Val() == "" {
		var err error
		data, err = r.SetGetChat(userID, req)
		if err != nil {
			return []models.Message{}, err
		}
	} else {
		err := r.jsonUnmarshel(&data, []byte(res.Val()))
		if err != nil {
			return []models.Message{}, err
		}
	}
	return data, nil
}

func (r *RedisChatCaching) SetGetChat(userID string, req models.ChatRequest) ([]models.Message, error) {
	contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data, err := r.chatClient.GetFriendChat(contextTimeout, &pb.GetFriendChatRequest{
		UserID:   userID,
		FriendID: req.FriendID,
		OffSet:   req.Offset,
		Limit:    req.Limit,
	})
	if err != nil {
		return []models.Message{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.Message{}, err
	}

	result := r.redis.Set(context.Background(), "chat:"+userID, profileByte, time.Hour)
	if result.Err() != nil {
		return []models.Message{}, err
	}
	var response []models.Message
	for _, v := range data.FriendChat {
		chatResponse := models.Message{
			SenderID:    v.SenderId,
			RecipientID: v.RecipientId,
			Content:     v.Content,
			Timestamp:   v.Timestamp,
		}
		response = append(response, chatResponse)

	}
	return response, nil
}
