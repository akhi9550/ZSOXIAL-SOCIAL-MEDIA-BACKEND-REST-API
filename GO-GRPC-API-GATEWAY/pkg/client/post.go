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

func (p *PostClient) CreatePost(userID int, req models.PostRequest, file []*multipart.FileHeader, users models.Tags) (models.PostResponse, error) {
	var fileDataList []byte

	// Loop through each file
	for _, file := range file {
		// Open the file
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

		// Append the file data to the list
		fileDataList = append(fileDataList, fileData...)
	}
	files := &pb.PostImage{Postimages: fileDataList}
	tags := &pb.Tags{
		User1: int64(users.User1),
		User2: int64(users.User2),
		User3: int64(users.User3),
		User4: int64(users.User4),
		User5: int64(users.User5),
	}
	data, err := p.Client.CreatePost(context.Background(), &pb.CreatePostRequest{
		Userid:  int64(userID),
		Caption: req.Caption,
		Post:    files,
		Typeid:  int64(req.TypeId),
		Tag:     tags,
	})
	if err != nil {
		return models.PostResponse{}, err
	}
	user := models.UserData{
		UserId:   uint(userID),
		Username: data.User.Username,
		Profile:  data.User.Imageurl,
	}
	tag := models.Tags{
		User1: uint(data.Tag.User1),
		User2: uint(data.Tag.User2),
		User3: uint(data.Tag.User3),
		User4: uint(data.Tag.User4),
		User5: uint(data.Tag.User5),
	}
	var imageUrls []models.Url
	for _, url := range data.Url {
		imageUrls = append(imageUrls, models.Url{ImageUrls: url.Imageurl})
	}
	fmt.Println("👺", imageUrls)
	return models.PostResponse{
		ID:        uint(data.Id),
		Author:    user,
		Caption:   data.Caption,
		Tag:       tag,
		ImageUrls: imageUrls,
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
	tag := models.Tags{
		User1: uint(data.Postresponse.Tag.User1),
		User2: uint(data.Postresponse.Tag.User2),
		User3: uint(data.Postresponse.Tag.User3),
		User4: uint(data.Postresponse.Tag.User4),
		User5: uint(data.Postresponse.Tag.User5),
	}
	var imageUrls []models.Url
	for _, url := range data.Postresponse.Url {
		imageUrls = append(imageUrls, models.Url{ImageUrls: url.Imageurl})
	}
	return models.PostResponse{
		ID:        uint(data.Postresponse.Id),
		Author:    user,
		Caption:   data.Postresponse.Caption,
		Tag:       tag,
		ImageUrls: imageUrls,
		Likes:     uint(data.Postresponse.Like),
		Comments:  uint(data.Postresponse.Comment),
		CreatedAt: data.Postresponse.CreatedAt.AsTime(),
	}, nil
}

func (p *PostClient) UpdatePost(userID int, req models.UpdatePostReq) (models.PostResponse, error) {
	data, err := p.Client.UpdatePost(context.Background(), &pb.UpdatePostRequest{
		Userid:  int64(userID),
		Postid:  int64(req.PostID),
		Caption: req.Caption,
		Typeid:  int64(req.TypeID),
		Tag: &pb.Tags{
			User1: int64(req.Tags.User1),
			User2: int64(req.Tags.User2),
			User3: int64(req.Tags.User3),
			User4: int64(req.Tags.User4),
			User5: int64(req.Tags.User5),
		},
	})
	if err != nil {
		return models.PostResponse{}, err
	}
	user := models.UserData{
		UserId:   uint(userID),
		Username: data.Updatepost.User.Username,
		Profile:  data.Updatepost.User.Imageurl,
	}
	tag := models.Tags{
		User1: uint(data.Updatepost.Tag.User1),
		User2: uint(data.Updatepost.Tag.User2),
		User3: uint(data.Updatepost.Tag.User3),
		User4: uint(data.Updatepost.Tag.User4),
		User5: uint(data.Updatepost.Tag.User5),
	}
	var imageUrls []models.Url
	for _, url := range data.Updatepost.Url {
		imageUrls = append(imageUrls, models.Url{ImageUrls: url.Imageurl})
	}
	return models.PostResponse{
		ID:        uint(data.Updatepost.Id),
		Author:    user,
		Caption:   data.Updatepost.Caption,
		Tag:       tag,
		ImageUrls: imageUrls,
		Likes:     uint(data.Updatepost.Like),
		Comments:  uint(data.Updatepost.Comment),
		CreatedAt: data.Updatepost.CreatedAt.AsTime(),
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
