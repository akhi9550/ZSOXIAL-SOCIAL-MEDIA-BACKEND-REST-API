package client

import (
	"context"
	"fmt"

	"github.com/akhi9550/post-svc/pkg/config"
	pb "github.com/akhi9550/post-svc/pkg/pb/auth"
	"github.com/akhi9550/post-svc/pkg/utils/models"
	"google.golang.org/grpc"
)

type clientAuth struct {
	Client pb.AuthServiceClient
}

func NewAuthClient(c *config.Config) *clientAuth {
	cc, err := grpc.Dial(c.AuthSvcUrl, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	pbClient := pb.NewAuthServiceClient(cc)

	return &clientAuth{
		Client: pbClient,
	}
}

func (c *clientAuth) CheckUserAvalilabilityWithUserID(userID int) bool {
	ok, _ := c.Client.CheckUserAvalilabilityWithUserID(context.Background(), &pb.CheckUserAvalilabilityWithUserIDRequest{
		Id: int64(userID),
	})
	return ok.Valid
}

func (c *clientAuth) UserData(userID int) (models.UserData, error) {
	data, err := c.Client.UserData(context.Background(), &pb.UserDataRequest{
		Id: int64(userID),
	})
	if err != nil {
		return models.UserData{}, err
	}
	return models.UserData{
		UserId:   uint(data.Id),
		Username: data.Username,
		Profile:  data.ProfilePhoto,
	}, nil
}

func (c *clientAuth) CheckUserAvalilabilityWithTagUserID(users []models.Tag) bool {
	var tags []string
	for _, user := range users {
		tags = append(tags, user.User)
	}
	tag := &pb.Tag{User: tags}
	ok, _ := c.Client.CheckUserAvalilabilityWithTagUserID(context.Background(), &pb.CheckUserAvalilabilityWithTagUserIDRequest{
		Tag: tag,
	})
	return ok.Valid
}

func (c *clientAuth) GetUserNameWithTagUserID(users []models.Tag) ([]models.Tag, error) {
	var tags []string
	for _, user := range users {
		tags = append(tags, user.User)
	}
	tag := &pb.Tag{User: tags}
	data, err := c.Client.GetUserNameWithTagUserID(context.Background(), &pb.GetUserNameWithTagUserIDRequest{
		Tag: tag,
	})
	if err != nil {
		return nil, err
	}

	var tagUsers []models.Tag
	for _, pbTagUser := range data.Name {
		tagUser := models.Tag{
			User: pbTagUser.Username,
		}
		tagUsers = append(tagUsers, tagUser)
	}

	return tagUsers, nil
}

func (c *clientAuth) GetFollowingUsers(userID int) ([]models.Users, error) {
	data, err := c.Client.GetFollowingUsers(context.Background(), &pb.GetFollowingUsersRequest{
		UserID: int64(userID),
	})
	if err != nil {
		return []models.Users{}, err
	}
	var followUsers []models.Users
	for _, pbUser := range data.User {
		followUser := models.Users{
			FollowingUser: int(pbUser.Followinguser),
		}
		followUsers = append(followUsers, followUser)
	}

	return followUsers, nil
}
