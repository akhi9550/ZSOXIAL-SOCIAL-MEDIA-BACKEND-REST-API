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

type RedisAuthCaching struct {
	redis      *redis.Client
	authClient interfaces.AuthClient
}

func NewRedisAuthCaching(redis *redis.Client, authClient interfaces.AuthClient) *RedisAuthCaching {
	return &RedisAuthCaching{
		redis:      redis,
		authClient: authClient,
	}
}

func (r *RedisAuthCaching) structMarshel(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (r *RedisAuthCaching) jsonUnmarshel(model interface{}, data []byte) error {
	return json.Unmarshal(data, model)
}

func (r *RedisAuthCaching) GetUserDetails(userID int) (models.UsersDetails, error) {
	userProfile := r.redis.Get(context.Background(), "user:"+strconv.Itoa(userID))
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

func (r *RedisAuthCaching) SetUserProfile(userID int) (models.UsersDetails, error) {
	userProfileFromService, err := r.authClient.UserDetails(userID)
	if err != nil {
		return models.UsersDetails{}, err
	}

	profileByte, err := r.structMarshel(userProfileFromService)
	if err != nil {
		return models.UsersDetails{}, err
	}

	result := r.redis.Set(context.Background(), "user:"+strconv.Itoa(userID), profileByte, time.Hour)
	if result.Err() != nil {
		return models.UsersDetails{}, err
	}

	return userProfileFromService, nil
}

func (r *RedisAuthCaching) GetSpecificUserDetails(userID int) (models.UsersDetails, error) {
	userProfile := r.redis.Get(context.Background(), "specificuser:"+strconv.Itoa(userID))
	var userProfileFinalResult models.UsersDetails

	if userProfile.Val() == "" {
		var err error
		userProfileFinalResult, err = r.SetSpecificUserProfile(userID)
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

func (r *RedisAuthCaching) SetSpecificUserProfile(userID int) (models.UsersDetails, error) {
	userProfileFromService, err := r.authClient.SpecificUserDetails(userID)
	if err != nil {
		return models.UsersDetails{}, err
	}

	profileByte, err := r.structMarshel(userProfileFromService)
	if err != nil {
		return models.UsersDetails{}, err
	}

	result := r.redis.Set(context.Background(), "specificuser:"+strconv.Itoa(userID), profileByte, time.Hour)
	if result.Err() != nil {
		return models.UsersDetails{}, err
	}

	return userProfileFromService, nil
}

func (r *RedisAuthCaching) ShowAllUsers(page, pageSize int) ([]models.UserDetailsAtAdmin, error) {
	res := r.redis.Get(context.Background(), "adminshowuser")
	var data []models.UserDetailsAtAdmin

	if res.Val() == "" {
		var err error
		data, err = r.SetShowAllUsers(page, pageSize)
		if err != nil {
			return []models.UserDetailsAtAdmin{}, err
		}
	} else {
		err := r.jsonUnmarshel(&data, []byte(res.Val()))
		if err != nil {
			return []models.UserDetailsAtAdmin{}, err
		}
	}
	return data, nil
}

func (r *RedisAuthCaching) SetShowAllUsers(page, pageSize int) ([]models.UserDetailsAtAdmin, error) {
	data, err := r.authClient.ShowAllUsers(page, pageSize)
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}

	result := r.redis.Set(context.Background(), "adminshowuser", profileByte, time.Hour)
	if result.Err() != nil {
		return []models.UserDetailsAtAdmin{}, err
	}
	return data, nil
}

func (r *RedisAuthCaching) GetAllPosts(page, pageSize int) ([]models.PostResponse, error) {
	res := r.redis.Get(context.Background(), "adminshowallposts")
	var data []models.PostResponse

	if res.Val() == "" {
		var err error
		data, err = r.SetGetAllPosts(page, pageSize)
		if err != nil {
			return []models.PostResponse{}, err
		}
	} else {
		err := r.jsonUnmarshel(&data, []byte(res.Val()))
		if err != nil {
			return []models.PostResponse{}, err
		}
	}
	return data, nil
}

func (r *RedisAuthCaching) SetGetAllPosts(page, pageSize int) ([]models.PostResponse, error) {
	data, err := r.authClient.GetAllPosts(page, pageSize)
	if err != nil {
		return []models.PostResponse{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.PostResponse{}, err
	}

	result := r.redis.Set(context.Background(), "adminshowallposts", profileByte, time.Hour)
	if result.Err() != nil {
		return []models.PostResponse{}, err
	}
	return data, nil
}

// func (r *RedisAuthCaching) ShowPostReports(page, pageSize int) ([]models.PostReports, error) {
// 	res := r.redis.Get(context.Background(), "adminshowpostreports")
// 	var data []models.PostReports

// 	if res.Val() == "" {
// 		var err error
// 		data, err = r.SetShowPostReports(page, pageSize)
// 		if err != nil {
// 			return []models.PostReports{}, err
// 		}
// 	} else {
// 		err := r.jsonUnmarshel(&data, []byte(res.Val()))
// 		if err != nil {
// 			return []models.PostReports{}, err
// 		}
// 	}
// 	return data, nil
// }

// func (r *RedisAuthCaching) SetShowPostReports(page, pageSize int) ([]models.PostReports, error) {
// 	data, err := r.authClient.ShowPostReports(page, pageSize)
// 	if err != nil {
// 		return []models.PostReports{}, err
// 	}

// 	profileByte, err := r.structMarshel(data)
// 	if err != nil {
// 		return []models.PostReports{}, err
// 	}

// 	result := r.redis.Set(context.Background(), "adminshowpostreports", profileByte, time.Hour)
// 	if result.Err() != nil {
// 		return []models.PostReports{}, err
// 	}
// 	return data, nil
// }

// func (r *RedisAuthCaching) ShowUserReports(page, pageSize int) ([]models.UserReports, error) {
// 	res := r.redis.Get(context.Background(), "adminshowuserreports")
// 	var data []models.UserReports

// 	if res.Val() == "" {
// 		var err error
// 		data, err = r.SetShowUserReports(page, pageSize)
// 		if err != nil {
// 			return []models.UserReports{}, err
// 		}
// 	} else {
// 		err := r.jsonUnmarshel(&data, []byte(res.Val()))
// 		if err != nil {
// 			return []models.UserReports{}, err
// 		}
// 	}
// 	return data, nil
// }

// func (r *RedisAuthCaching) SetShowUserReports(page, pageSize int) ([]models.UserReports, error) {
// 	data, err := r.authClient.ShowUserReports(page, pageSize)
// 	if err != nil {
// 		return []models.UserReports{}, err
// 	}

// 	profileByte, err := r.structMarshel(data)
// 	if err != nil {
// 		return []models.UserReports{}, err
// 	}

// 	result := r.redis.Set(context.Background(), "adminshowuserreports", profileByte, time.Hour)
// 	if result.Err() != nil {
// 		return []models.UserReports{}, err
// 	}
// 	return data, nil
// }

// func (r *RedisAuthCaching) ShowFollowREQ(userID int) ([]models.FollowingRequests, error) {
// 	userProfile := r.redis.Get(context.Background(), "followreq:"+strconv.Itoa(userID))
// 	var data []models.FollowingRequests

// 	if userProfile.Val() == "" {
// 		var err error
// 		data, err = r.SetShowFollowREQ(userID)
// 		if err != nil {
// 			return []models.FollowingRequests{}, err
// 		}
// 	} else {
// 		err := r.jsonUnmarshel(&data, []byte(userProfile.Val()))
// 		if err != nil {
// 			return []models.FollowingRequests{}, err
// 		}
// 	}
// 	return data, nil
// }

// func (r *RedisAuthCaching) SetShowFollowREQ(userID int) ([]models.FollowingRequests, error) {
// 	data, err := r.authClient.ShowFollowREQ(userID)
// 	if err != nil {
// 		return []models.FollowingRequests{}, err
// 	}

// 	profileByte, err := r.structMarshel(data)
// 	if err != nil {
// 		return []models.FollowingRequests{}, err
// 	}

// 	result := r.redis.Set(context.Background(), "followreq:"+strconv.Itoa(userID), profileByte, time.Hour)
// 	if result.Err() != nil {
// 		return []models.FollowingRequests{}, err
// 	}
// 	return data, nil
// }

// func (r *RedisAuthCaching) Following(userID int) ([]models.FollowingResponse, error) {
// 	userProfile := r.redis.Get(context.Background(), "following:"+strconv.Itoa(userID))
// 	var data []models.FollowingResponse

// 	if userProfile.Val() == "" {
// 		var err error
// 		data, err = r.SetFollowing(userID)
// 		if err != nil {
// 			return []models.FollowingResponse{}, err
// 		}
// 	} else {
// 		err := r.jsonUnmarshel(&data, []byte(userProfile.Val()))
// 		if err != nil {
// 			return []models.FollowingResponse{}, err
// 		}
// 	}
// 	return data, nil
// }

// func (r *RedisAuthCaching) SetFollowing(userID int) ([]models.FollowingResponse, error) {
// 	data, err := r.authClient.Following(userID)
// 	if err != nil {
// 		return []models.FollowingResponse{}, err
// 	}

// 	profileByte, err := r.structMarshel(data)
// 	if err != nil {
// 		return []models.FollowingResponse{}, err
// 	}

// 	result := r.redis.Set(context.Background(), "following:"+strconv.Itoa(userID), profileByte, time.Hour)
// 	if result.Err() != nil {
// 		return []models.FollowingResponse{}, err
// 	}
// 	fmt.Println("dafta", data)
// 	return data, nil
// }

// func (r *RedisAuthCaching) Follower(userID int) ([]models.FollowingResponse, error) {
// 	userProfile := r.redis.Get(context.Background(), "follower:"+strconv.Itoa(userID))
// 	var data []models.FollowingResponse

// 	if userProfile.Val() == "" {
// 		var err error
// 		data, err = r.SetFollower(userID)
// 		if err != nil {
// 			return []models.FollowingResponse{}, err
// 		}
// 	} else {
// 		err := r.jsonUnmarshel(&data, []byte(userProfile.Val()))
// 		if err != nil {
// 			return []models.FollowingResponse{}, err
// 		}
// 	}
// 	return data, nil
// }

// func (r *RedisAuthCaching) SetFollower(userID int) ([]models.FollowingResponse, error) {
// 	data, err := r.authClient.Follower(userID)
// 	if err != nil {
// 		return []models.FollowingResponse{}, err
// 	}

// 	profileByte, err := r.structMarshel(data)
// 	if err != nil {
// 		return []models.FollowingResponse{}, err
// 	}

// 	result := r.redis.Set(context.Background(), "follower:"+strconv.Itoa(userID), profileByte, time.Hour)
// 	if result.Err() != nil {
// 		return []models.FollowingResponse{}, err
// 	}
// 	return data, nil
// }
