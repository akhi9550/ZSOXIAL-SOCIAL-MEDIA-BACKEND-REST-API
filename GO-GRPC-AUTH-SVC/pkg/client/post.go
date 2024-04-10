package client

import (
	"context"
	"fmt"

	"github.com/akhi9550/auth-svc/pkg/config"
	pb "github.com/akhi9550/auth-svc/pkg/pb/post"
	"github.com/akhi9550/auth-svc/pkg/utils/models"
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

func (p *clientPost) ShowPostReports(page, count int) ([]models.PostReports, error) {
	data, err := p.Client.ShowPostReports(context.Background(), &pb.ShowPostReportsRequest{
		Page:  int64(page),
		Count: int64(count),
	})
	if err != nil {
		return []models.PostReports{}, err
	}
	var Report []models.PostReports
	for _, report := range data.Reports {
		reports := models.PostReports{
			ReportUserID: uint(report.RepostedUserid),
			PostID: uint(report.Postid),
			Report: report.Report,
		}

		Report = append(Report, reports)
	}

	return Report, nil
}

func (p *clientPost) GetAllPosts(page, count int) ([]models.PostResponse, error) {
	data, err := p.Client.GetAllposts(context.Background(), &pb.GetAllpostsRequest{
		Page:  int64(page),
		Count: int64(count),
	})
	if err != nil {
		return []models.PostResponse{}, err
	}

	var postResponses []models.PostResponse
	for _, post := range data.Posts {
		user := models.UserData{
			UserId:   uint(post.User.Userid),
			Username: post.User.Username,
			Profile:  post.User.Imageurl,
		}
		postResponse := models.PostResponse{
			ID:        uint(post.Id),
			Author:    user,
			Caption:   post.Caption,
			Url:       post.Url,
			Likes:     uint(post.Like),
			Comments:  uint(post.Comment),
			CreatedAt: post.CreatedAt.AsTime(),
		}

		postResponses = append(postResponses, postResponse)
	}

	return postResponses, nil
}

func (p *clientPost) CheckPostIDByID(postID int) bool {
	ok, _ := p.Client.CheckPostIDByID(context.Background(), &pb.CheckPostIDByIDRequest{
		PostID: int64(postID),
	})
	return ok.Exist
}

func (p *clientPost) RemovePost(postID int) error {
	_, err := p.Client.RemovePost(context.Background(), &pb.RemovePostRequest{
		PostID: int64(postID),
	})
	if err != nil {
		return err
	}
	return nil
}
