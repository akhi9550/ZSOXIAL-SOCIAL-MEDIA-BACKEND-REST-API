package client

import (
	"context"
	"fmt"

	"github.com/akhi9550/notification-svc/pkg/config"
	pb "github.com/akhi9550/notification-svc/pkg/pb/post"
	"google.golang.org/grpc"
)

type clientPost struct {
	Client pb.PostServiceClient
}

func NewPostClient(c *config.Config) *clientPost {
	cc, err := grpc.Dial(c.PostSvcUrl, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	pbClient := pb.NewPostServiceClient(cc)

	return &clientPost{
		Client: pbClient,
	}
}

func (p clientPost)GetUserId(postID int) (int, error){
	data,err:=p.Client.GetUserIDFromPost(context.Background(),&pb.GetUserIDRequest{
		PostID: int64(postID),
	})
	if err!=nil{
		return 0,err
	}
	return int(data.UserID),nil
}