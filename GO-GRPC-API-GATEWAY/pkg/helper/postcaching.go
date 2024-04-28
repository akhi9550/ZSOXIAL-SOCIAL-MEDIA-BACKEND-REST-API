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

type RedisPostCaching struct {
	redis      *redis.Client
	postClient interfaces.PostClient
}

func NewRedisPostCaching(redis *redis.Client, postClient interfaces.PostClient) *RedisPostCaching {
	return &RedisPostCaching{
		redis:      redis,
		postClient: postClient,
	}
}

func (r *RedisPostCaching) structMarshel(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (r *RedisPostCaching) jsonUnmarshel(model interface{}, data []byte) error {
	return json.Unmarshal(data, model)
}

func (r *RedisPostCaching) GetUserPost(userID int) ([]models.PostResponse, error) {
	res := r.redis.Get(context.Background(), "userpost:"+strconv.Itoa(userID))
	var data []models.PostResponse

	if res.Val() == "" {
		var err error
		data, err = r.SetGetUserPost(userID)
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

func (r *RedisPostCaching) SetGetUserPost(userID int) ([]models.PostResponse, error) {
	data, err := r.postClient.GetUserPost(userID)
	if err != nil {
		return []models.PostResponse{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.PostResponse{}, err
	}

	result := r.redis.Set(context.Background(), "userpost:"+strconv.Itoa(userID), profileByte, time.Hour)
	if result.Err() != nil {
		return []models.PostResponse{}, err
	}

	return data, nil
}

func (r *RedisPostCaching) GetAllPost(userID int) ([]models.PostResponse, error) {
	res := r.redis.Get(context.Background(), "getpost:"+strconv.Itoa(userID))
	var data []models.PostResponse

	if res.Val() == "" {
		var err error
		data, err = r.SetGetAllPost(userID)
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

func (r *RedisPostCaching) SetGetAllPost(userID int) ([]models.PostResponse, error) {
	data, err := r.postClient.GetAllPost(userID)
	if err != nil {
		return []models.PostResponse{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.PostResponse{}, err
	}

	result := r.redis.Set(context.Background(), "getpost:"+strconv.Itoa(userID), profileByte, time.Hour)
	if result.Err() != nil {
		return []models.PostResponse{}, err
	}

	return data, nil
}

func (r *RedisPostCaching) GetPost(userID, postID int) (models.PostResponse, error) {
	res := r.redis.Get(context.Background(), "getposts:"+strconv.Itoa(userID)+":"+strconv.Itoa(postID))
	var data models.PostResponse

	if res.Val() == "" {
		var err error
		data, err = r.SetGetPost(userID, postID)
		if err != nil {
			return models.PostResponse{}, err
		}
	} else {
		err := r.jsonUnmarshel(&data, []byte(res.Val()))
		if err != nil {
			return models.PostResponse{}, err
		}
	}
	return data, nil
}

func (r *RedisPostCaching) SetGetPost(userID, postID int) (models.PostResponse, error) {
	data, err := r.postClient.GetPost(userID, postID)
	if err != nil {
		return models.PostResponse{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return models.PostResponse{}, err
	}

	result := r.redis.Set(context.Background(), "getposts:"+strconv.Itoa(userID)+":"+strconv.Itoa(postID), profileByte, time.Hour)
	if result.Err() != nil {
		return models.PostResponse{}, err
	}

	return data, nil
}

func (r *RedisPostCaching) GetSavedPost(userID int) ([]models.PostResponse, error) {
	res := r.redis.Get(context.Background(), "savepost:"+strconv.Itoa(userID))
	var data []models.PostResponse

	if res.Val() == "" {
		var err error
		data, err = r.SetGetSavedPost(userID)
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

func (r *RedisPostCaching) SetGetSavedPost(userID int) ([]models.PostResponse, error) {
	data, err := r.postClient.GetSavedPost(userID)
	if err != nil {
		return []models.PostResponse{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.PostResponse{}, err
	}

	result := r.redis.Set(context.Background(), "savepost:"+strconv.Itoa(userID), profileByte, time.Hour)
	if result.Err() != nil {
		return []models.PostResponse{}, err
	}

	return data, nil
}

func (r *RedisPostCaching) GetAllArchivePost(userID int) ([]models.ArchivePostResponse, error) {
	res := r.redis.Get(context.Background(), "archivepost:"+strconv.Itoa(userID))
	var data []models.ArchivePostResponse

	if res.Val() == "" {
		var err error
		data, err = r.SetGetAllArchivePost(userID)
		if err != nil {
			return []models.ArchivePostResponse{}, err
		}
	} else {
		err := r.jsonUnmarshel(&data, []byte(res.Val()))
		if err != nil {
			return []models.ArchivePostResponse{}, err
		}
	}
	return data, nil
}

func (r *RedisPostCaching) SetGetAllArchivePost(userID int) ([]models.ArchivePostResponse, error) {
	data, err := r.postClient.GetAllArchivePost(userID)
	if err != nil {
		return []models.ArchivePostResponse{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.ArchivePostResponse{}, err
	}

	result := r.redis.Set(context.Background(), "archivepost:"+strconv.Itoa(userID), profileByte, time.Hour)
	if result.Err() != nil {
		return []models.ArchivePostResponse{}, err
	}

	return data, nil
}

func (r *RedisPostCaching) GetAllPostComments(postID int) ([]models.PostCommentResponse, error) {
	res := r.redis.Get(context.Background(), "allpost:"+strconv.Itoa(postID))
	var data []models.PostCommentResponse

	if res.Val() == "" {
		var err error
		data, err = r.SetGetAllPostComments(postID)
		if err != nil {
			return []models.PostCommentResponse{}, err
		}
	} else {
		err := r.jsonUnmarshel(&data, []byte(res.Val()))
		if err != nil {
			return []models.PostCommentResponse{}, err
		}
	}
	return data, nil
}

func (r *RedisPostCaching) SetGetAllPostComments(postID int) ([]models.PostCommentResponse, error) {
	data, err := r.postClient.GetAllPostComments(postID)
	if err != nil {
		return []models.PostCommentResponse{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.PostCommentResponse{}, err
	}

	result := r.redis.Set(context.Background(), "allpost:"+strconv.Itoa(postID), profileByte, time.Hour)
	if result.Err() != nil {
		return []models.PostCommentResponse{}, err
	}

	return data, nil
}

func (r *RedisPostCaching) ShowAllPostComments(postID int) ([]models.AllCommentsAndReplies, error) {
	res := r.redis.Get(context.Background(), "showallpost:"+strconv.Itoa(postID))
	var data []models.AllCommentsAndReplies

	if res.Val() == "" {
		var err error
		data, err = r.SetShowAllPostComments(postID)
		if err != nil {
			return []models.AllCommentsAndReplies{}, err
		}
	} else {
		err := r.jsonUnmarshel(&data, []byte(res.Val()))
		if err != nil {
			return []models.AllCommentsAndReplies{}, err
		}
	}
	return data, nil
}

func (r *RedisPostCaching) SetShowAllPostComments(postID int) ([]models.AllCommentsAndReplies, error) {
	data, err := r.postClient.ShowAllPostComments(postID)
	if err != nil {
		return []models.AllCommentsAndReplies{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.AllCommentsAndReplies{}, err
	}

	result := r.redis.Set(context.Background(), "showallpost:"+strconv.Itoa(postID), profileByte, time.Hour)
	if result.Err() != nil {
		return []models.AllCommentsAndReplies{}, err
	}
	return data, nil
}

func (r *RedisPostCaching) GetStory(userID,viewUser int) ([]models.CreateStoryResponses, error) {
	res := r.redis.Get(context.Background(), "viewstory:"+strconv.Itoa(viewUser)+":"+strconv.Itoa(userID))
	var data []models.CreateStoryResponses

	if res.Val() == "" {
		var err error
		data, err = r.SetGetStory(userID,viewUser)
		if err != nil {
			return []models.CreateStoryResponses{}, err
		}
	} else {
		err := r.jsonUnmarshel(&data, []byte(res.Val()))
		if err != nil {
			return []models.CreateStoryResponses{}, err
		}
	}
	return data, nil
}

func (r *RedisPostCaching) SetGetStory(userID,viewUser int) ([]models.CreateStoryResponses, error) {
	data, err := r.postClient.GetStory(userID,viewUser)
	if err != nil {
		return []models.CreateStoryResponses{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return []models.CreateStoryResponses{}, err
	}

	result := r.redis.Set(context.Background(),"viewstory:"+strconv.Itoa(viewUser)+":"+strconv.Itoa(userID), profileByte, time.Hour)
	if result.Err() != nil {
		return []models.CreateStoryResponses{}, err
	}
	return data, nil
}

func (r *RedisPostCaching) StoryDetails(userID,storyID int) (models.StoryDetails, error) {
	res := r.redis.Get(context.Background(), "story:"+strconv.Itoa(storyID))
	var data models.StoryDetails

	if res.Val() == "" {
		var err error
		data, err = r.SetStoryDetails(userID,storyID)
		if err != nil {
			return models.StoryDetails{}, err
		}
	} else {
		err := r.jsonUnmarshel(&data, []byte(res.Val()))
		if err != nil {
			return models.StoryDetails{}, err
		}
	}
	return data, nil
}

func (r *RedisPostCaching) SetStoryDetails(userID,storyID int) (models.StoryDetails, error) {
	data, err := r.postClient.StoryDetails(userID,storyID)
	if err != nil {
		return models.StoryDetails{}, err
	}

	profileByte, err := r.structMarshel(data)
	if err != nil {
		return models.StoryDetails{}, err
	}

	result := r.redis.Set(context.Background(), "story:"+strconv.Itoa(userID), profileByte, time.Hour)
	if result.Err() != nil {
		return models.StoryDetails{}, err
	}
	return data, nil
}
