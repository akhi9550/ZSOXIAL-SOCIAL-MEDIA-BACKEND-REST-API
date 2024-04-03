package client

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/config"
	pb "github.com/akhi9550/api-gateway/pkg/pb/post"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"google.golang.org/grpc"
)

type PostClient struct {
	Client pb.PostServiceClient
}

func NewPostClient(cfg config.Config) interfaces.PostClient {
	grpcConnection, err := grpc.Dial(cfg.PostSvcUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect", err)
	}

	grpcClient := pb.NewPostServiceClient(grpcConnection)

	return &PostClient{
		Client: grpcClient,
	}
}

func (p *PostClient) CreatePost(userID int, req models.PostRequest, file *multipart.FileHeader, user []string) (models.PostResponse, error) {
	f, err := file.Open()
	if err != nil {
		return models.PostResponse{}, err
	}
	defer f.Close()

	// Read the file content
	fileData, err := io.ReadAll(f)
	if err != nil {
		return models.PostResponse{}, err
	}
	files := &pb.PostPhoto{Url: fileData}
	tag := &pb.Tag{User: user}
	data, err := p.Client.CreatePost(context.Background(), &pb.CreatePostRequest{
		Userid:  int64(userID),
		Caption: req.Caption,
		Post:    files,
		Typeid:  int64(req.TypeId),
		Tag:     tag,
	})
	if err != nil {
		return models.PostResponse{}, err
	}
	users := models.UserData{
		UserId:   uint(userID),
		Username: data.User.Username,
		Profile:  data.User.Imageurl,
	}
	var tags []models.Tag
	for _, tagStr := range data.Tag {
		var tag models.Tag
		tag.User = tagStr
		tags = append(tags, tag)
	}

	return models.PostResponse{
		ID:        uint(data.Id),
		Author:    users,
		Caption:   data.Caption,
		Tag:       tags,
		Url:       data.Url,
		Likes:     uint(data.Like),
		Comments:  uint(data.Comment),
		CreatedAt: data.CreatedAt.AsTime(),
	}, nil

}
func (p *PostClient) GetPost(userID int, postID int) (models.PostResponse, error) {
	data, err := p.Client.GetPost(context.Background(), &pb.GetPostRequest{
		Userid: int64(userID),
		Postid: int64(postID),
	})
	if err != nil {
		return models.PostResponse{}, err
	}
	user := models.UserData{
		UserId:   uint(userID),
		Username: data.Postresponse.User.Username,
		Profile:  data.Postresponse.User.Imageurl,
	}
	var tags []models.Tag
	for _, tagStr := range data.Postresponse.Tag {
		var tag models.Tag
		tag.User = tagStr
		tags = append(tags, tag)
	}

	return models.PostResponse{
		ID:        uint(data.Postresponse.Id),
		Author:    user,
		Caption:   data.Postresponse.Caption,
		Tag:       tags,
		Url:       data.Postresponse.Url,
		Likes:     uint(data.Postresponse.Like),
		Comments:  uint(data.Postresponse.Comment),
		CreatedAt: data.Postresponse.CreatedAt.AsTime(),
	}, nil
}

func (p *PostClient) UpdatePost(userID int, req models.UpdatePostReq, user []string) (models.UpdateResponse, error) {
	tag := &pb.Tag{User: user}
	data, err := p.Client.UpdatePost(context.Background(), &pb.UpdatePostRequest{
		Userid:  int64(userID),
		Postid:  int64(req.PostID),
		Caption: req.Caption,
		Typeid:  int64(req.TypeID),
		Tag:     tag,
	})
	if err != nil {
		return models.UpdateResponse{}, err
	}
	users := models.UserData{
		UserId:   uint(userID),
		Username: data.User.Username,
		Profile:  data.User.Imageurl,
	}

	var tags []models.Tag
	for _, tagStr := range data.Tag {
		var tag models.Tag
		tag.User = tagStr
		tags = append(tags, tag)
	}
	return models.UpdateResponse{
		ID:        uint(data.Id),
		Author:    users,
		Caption:   data.Caption,
		Tag:       tags,
		Url:       data.Url,
		Likes:     uint(data.Like),
		Comments:  uint(data.Comment),
		CreatedAt: data.CreatedAt.AsTime(),
	}, nil
}

func (p *PostClient) DeletePost(userID int, postID int) error {
	_, err := p.Client.DeletePost(context.Background(), &pb.DeletePostRequest{
		Userid: int64(userID),
		Postid: int64(postID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *PostClient) GetAllPost(userID int) ([]models.PostResponse, error) {
	data, err := p.Client.GetAllPost(context.Background(), &pb.GetAllPostRequest{
		Userid: int64(userID),
	})
	if err != nil {
		return nil, err
	}

	var postResponses []models.PostResponse
	for _, post := range data.Allpost {
		user := models.UserData{
			UserId:   uint(userID),
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

func (p *PostClient) ArchivePost(userID, PostID int) error {
	_, err := p.Client.ArchivePost(context.Background(), &pb.ArchivePostRequest{
		Userid: int64(userID),
		Postid: int64(PostID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *PostClient) UnArchivePost(userID, PostID int) error {
	_, err := p.Client.UnArchivePost(context.Background(), &pb.UnArchivePostrequest{
		Userid: int64(userID),
		Postid: int64(PostID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *PostClient) GetAllArchivePost(userID int) ([]models.ArchivePostResponse, error) {
	data, err := p.Client.GetAllArchivePost(context.Background(), &pb.GetAllArchivePostRequest{
		Userid: int64(userID),
	})
	if err != nil {
		return []models.ArchivePostResponse{}, err
	}
	var postResponses []models.ArchivePostResponse
	for _, post := range data.Allpost {
		postResponse := models.ArchivePostResponse{
			ID:        uint(post.Id),
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

func (p *PostClient) LikePost(userID int, postID int) (models.LikePostResponse, error) {
	data, err := p.Client.LikePost(context.Background(), &pb.LikePostRequest{
		Userid: int64(userID),
		Postid: int64(postID),
	})
	if err != nil {
		return models.LikePostResponse{}, err
	}
	return models.LikePostResponse{
		UserID:    uint(data.Userid),
		LikedUser: data.LikedUser,
		Profile:   data.Posturl,
		CreatedAt: data.GetCreatedAt().AsTime(),
	}, nil
}

func (p *PostClient) UnLinkPost(userID int, postID int) error {
	_, err := p.Client.UnLinkPost(context.Background(), &pb.UnLikePostRequest{
		Userid: int64(userID),
		Postid: int64(postID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *PostClient) PostComment(userID int, postData models.PostCommentReq) (models.PostCommentResponse, error) {
	data, err := p.Client.PostComment(context.Background(), &pb.PostCommentRequest{
		Userid:  int64(userID),
		Postid:  int64(postData.PostID),
		Comment: postData.Comment,
	})
	if err != nil {
		return models.PostCommentResponse{}, err
	}
	return models.PostCommentResponse{
		UserID:      uint(data.Userid),
		CommentUser: data.CommentedUser,
		Profile:     data.Posturl,
		Comment:     data.Comment,
		CreatedAt:   data.GetCreatedAt().AsTime(),
	}, nil

}

func (p *PostClient) SavedPost(userID, postID int) error {
	_, err := p.Client.SavedPost(context.Background(), &pb.SavedPostRequest{
		Userid: int64(userID),
		Postid: int64(postID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *PostClient) UnSavedPost(userID, postID int) error {
	_, err := p.Client.UnSavedPost(context.Background(), &pb.UnSavedPostRequest{
		Userid: int64(userID),
		Postid: int64(postID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (p *PostClient) GetSavedPost(userID int) ([]models.PostResponse, error) {
	data, err := p.Client.GetSavedPost(context.Background(), &pb.GetSavedPostRequest{
		Userid: int64(userID),
	})
	if err != nil {
		return []models.PostResponse{}, err
	}
	var postResponses []models.PostResponse
	for _, post := range data.Allpost {
		user := models.UserData{
			UserId:   uint(userID),
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
