package helper

import (
	"context"
	"encoding/json"
	"time"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/redis/go-redis/v9"
)

type RedisChatCaching struct {
	redis      *redis.Client
	chatClient interfaces.ChatClient
}

func NewRedisChatCaching(redis *redis.Client, chatClient interfaces.ChatClient) *RedisChatCaching {
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
	res := r.redis.Get(context.Background(), "chat:"+userID+":"+req.FriendID)
	var data []models.Message

	if res.Val() == "" {
		var err error
		data, err = r.SetGetChat(userID, req)
		if err != nil {
			return []models.Message{}, err
		}
	} else {
		err := r.jsonUnmarshel(&res, []byte(res.Val()))
		if err != nil {
			return []models.Message{}, err
		}
	}
	return data, nil
}

func (r *RedisChatCaching) SetGetChat(userID string, req models.ChatRequest) ([]models.Message, error) {
	data, err := r.chatClient.GetChat(userID, req)
	if err != nil {
		return []models.Message{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.Message{}, err
	}

	result := r.redis.Set(context.Background(), "chat:"+userID+":"+req.FriendID, profileByte, time.Hour)
	if result.Err() != nil {
		return []models.Message{}, err
	}

	return data, nil
}
