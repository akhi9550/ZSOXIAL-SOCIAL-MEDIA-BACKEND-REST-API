package helper

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	pb "github.com/akhi9550/api-gateway/pkg/pb/notification"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/redis/go-redis/v9"
)

type RedisNotificationCaching struct {
	redis              *redis.Client
	notificationClient pb.NotificationServiceClient
}

func NewRedisNotificationCaching(redis *redis.Client, notificationClient pb.NotificationServiceClient) *RedisNotificationCaching {
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
		err := r.jsonUnmarshel(&data, []byte(res.Val()))
		if err != nil {
			return []models.NotificationResponse{}, err
		}
	}
	return data, nil
}

func (r *RedisNotificationCaching) SetGetNotification(userID int, req models.NotificationPagination) ([]models.NotificationResponse, error) {
	contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	data, err := r.notificationClient.GetNotification(contextTimeout, &pb.GetNotificationRequest{
		UserID: int64(userID),
		Offset: int64(req.Offset),
		Limit:  int64(req.Limit),
	})
	if err != nil {
		return []models.NotificationResponse{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.NotificationResponse{}, err
	}

	result := r.redis.Set(context.Background(), "notification:"+strconv.Itoa(userID), profileByte, time.Hour)
	if result.Err() != nil {
		return []models.NotificationResponse{}, err
	}

	var response []models.NotificationResponse
	for _, v := range data.Notification {
		notificationResponse := models.NotificationResponse{
			UserID:    int(v.UserID),
			Username:  v.Username,
			Profile:   v.Profile,
			PostID:    int(v.PostID),
			Message:   v.Message,
			CreatedAt: v.Time,
		}
		response = append(response, notificationResponse)
	}
	return response, nil
}
