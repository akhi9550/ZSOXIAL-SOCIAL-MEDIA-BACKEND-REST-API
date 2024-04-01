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
	fmt.Println("valid", ok)
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
