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
	var users []models.Tag
	for _, user := range req.Tag.User {
		tag := models.Tag{User: user}
		users = append(users, tag)
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
	var tags []string
	for _, tag := range data.Tag {
		tags = append(tags, tag.User)
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
	var tags []string
	for _, tag := range data.Tag {
		tags = append(tags, tag.User)
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
func (p *PostServer) GetAllPost(ctx context.Context, req *pb.GetAllPostRequest) (*pb.GetAllPostResponse, error) {
	userID := req.Userid
	posts, err := p.postUseCase.GetAllPost(int(userID))
	if err != nil {
		return nil, err
	}

	var allPostResponses []*pb.CreatePostResponse
	for _, post := range posts {
		userData := &pb.UserData{
			Userid:   userID,
			Username: post.Author.Username,
			Imageurl: post.Author.Profile,
		}
		// var tags []string
		// for _, tag := range post.Tag {
		// 	tags = append(tags, tag.User)
		// }

		details := &pb.CreatePostResponse{
			Id:        int64(post.ID),
			User:      userData,
			Caption:   post.Caption,
			Url:       post.Url,
			Like:      int64(post.Likes),
			Comment:   int64(post.Comments),
			CreatedAt: timestamppb.New(post.CreatedAt),
		}

		allPostResponses = append(allPostResponses, details)
	}

	return &pb.GetAllPostResponse{
		Allpost: allPostResponses,
	}, nil
}

func (p *PostServer) ArchivePost(ctx context.Context, req *pb.ArchivePostRequest) (*pb.ArchivePostResponse, error) {
	userID, postID := req.Userid, req.Postid
	err := p.postUseCase.ArchivePost(int(userID), int(postID))
	if err != nil {
		return &pb.ArchivePostResponse{}, err
	}
	return &pb.ArchivePostResponse{}, nil
}

func (p *PostServer) UnArchivePost(ctx context.Context, req *pb.UnArchivePostrequest) (*pb.UnArchivePostResponse, error) {
	userID, postID := req.Userid, req.Postid
	err := p.postUseCase.UnArchivePost(int(userID), int(postID))
	if err != nil {
		return &pb.UnArchivePostResponse{}, err
	}
	return &pb.UnArchivePostResponse{}, nil
}

func (p *PostServer) GetAllArchivePost(ctx context.Context, req *pb.GetAllArchivePostRequest) (*pb.GetAllArchivePostResponse, error) {
	userID := req.Userid
	posts, err := p.postUseCase.GetAllArchivePost(int(userID))
	if err != nil {
		return nil, err
	}
	var allPostResponses []*pb.ArchivePostResponses
	for _, post := range posts {
		details := &pb.ArchivePostResponses{
			Id:        int64(post.ID),
			Caption:   post.Caption,
			Url:       post.Url,
			Like:      int64(post.Likes),
			Comment:   int64(post.Comments),
			CreatedAt: timestamppb.New(post.CreatedAt),
		}

		allPostResponses = append(allPostResponses, details)
	}
	return &pb.GetAllArchivePostResponse{
		Allpost: allPostResponses,
	}, nil
}

func (p *PostServer) LikePost(ctx context.Context, req *pb.LikePostRequest) (*pb.LikePostResponse, error) {
	userID, postID := req.Userid, req.Postid
	data, err := p.postUseCase.LikePost(int(userID), int(postID))
	if err != nil {
		return &pb.LikePostResponse{}, err
	}
	return &pb.LikePostResponse{
		Userid:    int64(data.UserID),
		LikedUser: data.LikedUser,
		Posturl:   data.Profile,
		CreatedAt: timestamppb.New(data.CreatedAt),
	}, nil
}

func (p *PostServer) UnLinkPost(ctx context.Context, req *pb.UnLikePostRequest) (*pb.UnLikePostResponse, error) {
	userID, postID := req.Userid, req.Postid
	err := p.postUseCase.UnLinkPost(int(userID), int(postID))
	if err != nil {
		return &pb.UnLikePostResponse{}, err
	}
	return &pb.UnLikePostResponse{}, nil
}

func (p *PostServer) PostComment(ctx context.Context, req *pb.PostCommentRequest) (*pb.PostCommentResponse, error) {
	userID := req.Userid
	reqdata := models.PostCommentReq{
		PostID:  uint(req.Postid),
		Comment: req.Comment,
	}
	data, err := p.postUseCase.PostComment(int(userID), reqdata)
	if err != nil {
		return &pb.PostCommentResponse{}, err
	}
	return &pb.PostCommentResponse{
		Userid:        int64(data.UserID),
		CommentedUser: data.CommentedUser,
		Posturl:       data.Profile,
		Comment:       data.Comment,
		CreatedAt:     timestamppb.New(data.CreatedAt),
	}, nil
}
