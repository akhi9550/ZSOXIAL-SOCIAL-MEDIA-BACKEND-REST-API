package service

import (
	"context"

	pb "github.com/akhi9550/post-svc/pkg/pb/post"
	interfaces "github.com/akhi9550/post-svc/pkg/usecase/interface"
	"github.com/akhi9550/post-svc/pkg/utils/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostServer struct {
	postUseCase interfaces.PostUseCase
	pb.UnimplementedPostServiceServer
}

func NewPostServer(UseCasePost interfaces.PostUseCase) pb.PostServiceServer {
	return &PostServer{
		postUseCase: UseCasePost,
	}
}
func (p *PostServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	userID := req.Userid
	createPost := models.PostRequest{
		Caption: req.Caption,
		TypeId:  uint(req.Typeid),
	}
	File := req.Post.Url
	users := models.Tags{
		User1: uint(req.Tag.User1),
		User2: uint(req.Tag.User2),
		User3: uint(req.Tag.User3),
		User4: uint(req.Tag.User4),
		User5: uint(req.Tag.User5),
	}
	data, err := p.postUseCase.CreatePost(int(userID), createPost, File, users)
	if err != nil {
		return &pb.CreatePostResponse{}, err
	}
	Users := &pb.UserData{
		Userid:   userID,
		Username: data.Author.Username,
		Imageurl: data.Author.Profile,
	}
	tags := &pb.Tags{
		User1: int64(data.Tag.User1),
		User2: int64(data.Tag.User2),
		User3: int64(data.Tag.User3),
		User4: int64(data.Tag.User4),
		User5: int64(data.Tag.User5),
	}

	return &pb.CreatePostResponse{
		Id:        int64(data.ID),
		User:      Users,
		Caption:   data.Caption,
		Tag:       tags,
		Url:       data.Url,
		Like:      int64(data.Likes),
		Comment:   int64(data.Comments),
		CreatedAt: timestamppb.New(data.CreatedAt),
	}, nil
}

func (p *PostServer) GetPost(ctx context.Context, req *pb.GetPostRequest) (*pb.GetPostResponse, error) {
	userID, postID := req.Userid, req.Postid
	data, err := p.postUseCase.GetPost(int(userID), int(postID))
	if err != nil {
		return &pb.GetPostResponse{}, err
	}
	Users := &pb.UserData{
		Userid:   userID,
		Username: data.Author.Username,
		Imageurl: data.Author.Profile,
	}
	tags := &pb.Tags{
		User1: int64(data.Tag.User1),
		User2: int64(data.Tag.User2),
		User3: int64(data.Tag.User3),
		User4: int64(data.Tag.User4),
		User5: int64(data.Tag.User5),
	}
	details := &pb.CreatePostResponse{
		Id:        int64(data.ID),
		User:      Users,
		Caption:   data.Caption,
		Tag:       tags,
		Url:       data.Url,
		Like:      int64(data.Likes),
		Comment:   int64(data.Comments),
		CreatedAt: timestamppb.New(data.CreatedAt),
	}
	return &pb.GetPostResponse{
		Postresponse: details,
	}, nil
}

func (p *PostServer) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.UpdatePostResponse, error) {
	userID := req.Userid
	user := models.UpdatePostReq{
		PostID:  uint(req.Postid),
		Caption: req.Caption,
		TypeID:  uint(req.Typeid),
	}
	var users []models.Tag
	for _, user := range req.Tag.User {
		tag := models.Tag{User: user}
		users = append(users, tag)
	}
	data, err := p.postUseCase.UpdatePost(int(userID), user, users)
	if err != nil {
		return &pb.UpdatePostResponse{}, err
	}
	Users := &pb.UserData{
		Userid:   userID,
		Username: data.Author.Username,
		Imageurl: data.Author.Profile,
	}
	var tags []string
	for _, tag := range data.Tag {
		tags = append(tags, tag.User)
	}
	return &pb.UpdatePostResponse{
		Id:        int64(data.ID),
		User:      Users,
		Caption:   data.Caption,
		Tag:       tags,
		Url:       data.Url,
		Like:      int64(data.Likes),
		Comment:   int64(data.Comments),
		CreatedAt: timestamppb.New(data.CreatedAt),
	}, nil
}

func (p *PostServer) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	userID, postID := req.Userid, req.Postid
	err := p.postUseCase.DeletePost(int(userID), int(postID))
	if err != nil {
		return &pb.DeletePostResponse{}, err
	}
	return &pb.DeletePostResponse{}, nil
}
