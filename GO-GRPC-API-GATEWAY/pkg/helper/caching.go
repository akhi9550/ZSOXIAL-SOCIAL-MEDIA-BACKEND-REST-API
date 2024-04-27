package helper

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	pb "github.com/akhi9550/api-gateway/pkg/pb/auth"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/redis/go-redis/v9"
)

type RedisCaching struct {
	redis      *redis.Client
	authClient pb.AuthServiceClient
}

func NewRedisCaching(redis *redis.Client, authClient pb.AuthServiceClient) *RedisCaching {
	return &RedisCaching{
		redis:      redis,
		authClient: authClient,
	}
}

func (r *RedisCaching) GetUserDetails(userID int) (models.UsersDetails, error) {
	userProfile := r.redis.Get(context.Background(), "user:-"+strconv.Itoa(userID))
	var userProfileFinalResult models.UsersDetails

	if userProfile.Val() == "" {
		var err error
		userProfileFinalResult, err = r.SetUserProfile(userID)
		if err != nil {
			return models.UsersDetails{}, err
		}
	} else {
		err := r.jsonUnmarshel(&userProfileFinalResult, []byte(userProfile.Val()))
		if err != nil {
			return models.UsersDetails{}, err
		}
	}

	return userProfileFinalResult, nil
}

func (r *RedisCaching) SetUserProfile(userID int) (models.UsersDetails, error) {

	contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userProfileFromService, err := r.authClient.UserDetails(contextTimeout, &pb.UserDetailsRequest{Id: int64(userID)})
	if err != nil {
		return models.UsersDetails{}, err
	}

	profileByte, err := r.structMarshel(userProfileFromService)
	if err != nil {
		return models.UsersDetails{}, err
	}

	result := r.redis.Set(context.Background(), "user:-"+strconv.Itoa(userID), profileByte, time.Hour)
	if result.Err() != nil {
		return models.UsersDetails{}, err
	}

	return models.UsersDetails{
		Firstname: userProfileFromService.Responsedata.Firstname,
		Lastname:  userProfileFromService.Responsedata.Lastname,
		Username:  userProfileFromService.Responsedata.Username,
		Dob:       userProfileFromService.Responsedata.Dob,
		Gender:    userProfileFromService.Responsedata.Gender,
		Phone:     userProfileFromService.Responsedata.Phone,
		Email:     userProfileFromService.Responsedata.Email,
		Bio:       userProfileFromService.Responsedata.Bio,
		Imageurl:  userProfileFromService.Responsedata.ProfilePhoto,
		Following: int(userProfileFromService.ResponseFollowigs.Following),
		Follower:  int(userProfileFromService.ResponseFollowigs.Follower),
	}, nil
}

func (r *RedisCaching) structMarshel(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (r *RedisCaching) jsonUnmarshel(model interface{}, data []byte) error {
	return json.Unmarshal(data, model)
}
