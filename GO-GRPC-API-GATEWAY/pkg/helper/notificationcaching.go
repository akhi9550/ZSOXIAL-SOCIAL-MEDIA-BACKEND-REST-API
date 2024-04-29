package helper

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/redis/go-redis/v9"
)

type RedisNotificationCaching struct {
	redis              *redis.Client
	notificationClient interfaces.NotificationClient
}

func NewRedisNotificationCaching(redis *redis.Client, notificationClient interfaces.NotificationClient) *RedisNotificationCaching {
	return &RedisNotificationCaching{
		redis:              redis,
		notificationClient: notificationClient,
	}
}

func (r *RedisNotificationCaching) structMarshel(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (r *RedisNotificationCaching) jsonUnmarshel(model interface{}, data []byte) error {
	return json.Unmarshal(data, model)
}

func (r *RedisNotificationCaching) GetNotification(userID int, req models.NotificationPagination) ([]models.NotificationResponse, error) {
	res := r.redis.Get(context.Background(), "notification:"+strconv.Itoa(userID))
	var data []models.NotificationResponse

	if res.Val() == "" {
		var err error
		data, err = r.SetGetNotification(userID, req)
		if err != nil {
			return []models.NotificationResponse{}, err
		}
	} else {
		err := r.jsonUnmarshel(&res, []byte(res.Val()))
		if err != nil {
			return []models.NotificationResponse{}, err
		}
	}
	return data, nil
}

func (r *RedisNotificationCaching) SetGetNotification(userID int, req models.NotificationPagination) ([]models.NotificationResponse, error) {
	data, err := r.notificationClient.GetNotification(userID, req)
	if err != nil {
		return []models.NotificationResponse{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.NotificationResponse{}, err
	}

	result := r.redis.Set(context.Background(),  "notification:"+strconv.Itoa(userID), profileByte, time.Hour)
	if result.Err() != nil {
		return []models.NotificationResponse{}, err
	}

	return data, nil
}
