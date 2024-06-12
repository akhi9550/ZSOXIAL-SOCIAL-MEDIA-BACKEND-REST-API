package service

import (
	"context"

	pb "github.com/akhi9550/post-svc/pkg/pb/post"
	interfaces "github.com/akhi9550/post-svc/pkg/usecase/interface"
	"github.com/akhi9550/post-svc/pkg/utils/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostServer struct {
	postUseCase  interfaces.PostUseCase
	storyUseCase interfaces.StoryUseCase
	pb.UnimplementedPostServiceServer
}

func NewPostServer(UseCasePost interfaces.PostUseCase, UseCaseStory interfaces.StoryUseCase) pb.PostServiceServer {
	return &PostServer{
		postUseCase:  UseCasePost,
		storyUseCase: UseCaseStory,
	}
}
func (p *PostServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.CreatePostResponse, error) {
	userID := req.Userid
	createPost := models.PostRequest{
		Caption: req.Caption,
		TypeId:  req.Typeid,
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
		TypeID:  req.Typeid,
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

func (p *PostServer) UnLikePost(ctx context.Context, req *pb.UnLikePostRequest) (*pb.UnLikePostResponse, error) {
	userID, postID := req.Userid, req.Postid
	err := p.postUseCase.UnLikePost(int(userID), int(postID))
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

func (p *PostServer) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	userID, commentID := req.Userid, req.Commentid
	err := p.postUseCase.DeleteComment(int(userID), int(commentID))
	if err != nil {
		return &pb.DeleteCommentResponse{}, err
	}
	return &pb.DeleteCommentResponse{}, nil
}

func (p *PostServer) GetAllPostComments(ctx context.Context, req *pb.GetAllCommentsRequest) (*pb.GetAllCommentsResponse, error) {
	postID := req.Postid
	data, err := p.postUseCase.GetAllPostComments(int(postID))
	if err != nil {
		return nil, err
	}
	var allCommentResponses []*pb.PostCommentResponses
	for _, post := range data {
		details := &pb.PostCommentResponses{
			Userid:        int64(post.UserID),
			CommentedUser: post.CommentedUser,
			Posturl:       post.Profile,
			Commentid:     int64(post.CommentID),
			Comment:       post.Comment,
			CreatedAt:     timestamppb.New(post.CreatedAt),
		}
		allCommentResponses = append(allCommentResponses, details)
	}
	return &pb.GetAllCommentsResponse{
		Allcomments: allCommentResponses,
	}, nil
}

func (p *PostServer) ReplyComment(ctx context.Context, req *pb.ReplyCommentRequest) (*pb.ReplyCommentResponse, error) {
	userID := req.Replyuserid
	reqdata := models.ReplyCommentReq{
		CommentID: uint(req.Commentid),
		Reply:     req.Replies,
	}
	data, err := p.postUseCase.ReplyComment(int(userID), reqdata)
	if err != nil {
		return &pb.ReplyCommentResponse{}, err
	}
	comment := &pb.PostCommentResponse{
		Userid:        int64(data.Comment.UserID),
		CommentedUser: data.Comment.CommentedUser,
		Posturl:       data.Comment.Profile,
		Comment:       data.Comment.Comment,
		CreatedAt:     timestamppb.New(data.Comment.CreatedAt),
	}
	reply := &pb.Reply{
		Replyuserid: int64(data.Reply.UserID),
		Replieduser: data.Reply.ReplyUser,
		Posturl:     data.Reply.Profile,
		Replies:     data.Reply.Reply,
		CreatedAt:   timestamppb.New(data.Reply.CreatedAt),
	}
	return &pb.ReplyCommentResponse{
		Comment: comment,
		Reply:   reply,
	}, nil
}

func (p *PostServer) ShowAllPostComments(ctx context.Context, req *pb.ShowAllPostCommentsRequest) (*pb.ShowAllPostCommentsResponse, error) {
	postID := req.Postid
	data, err := p.postUseCase.ShowAllPostComments(int(postID))
	if err != nil {
		return nil, err
	}

	var allCommentsAndReplies []*pb.AllPostCommentsResponse
	for _, commentData := range data {
		comment := &pb.AllPostCommentsResponse{
			CommentedUser: commentData.CommentUser,
			Posturl:       commentData.Profile,
			Comment:       commentData.Comment,
			CreatedAt:     timestamppb.New(commentData.CreatedAt),
			Reply:         make([]*pb.Replies, len(commentData.Reply)),
		}

		for i, replyData := range commentData.Reply {
			comment.Reply[i] = &pb.Replies{
				Replieduser: replyData.ReplyUser,
				Posturl:     replyData.Profile,
				Replies:     replyData.Reply,
				CreatedAt:   timestamppb.New(replyData.CreatedAt),
			}
		}

		allCommentsAndReplies = append(allCommentsAndReplies, comment)
	}

	return &pb.ShowAllPostCommentsResponse{
		Comments: allCommentsAndReplies,
	}, nil
}

func (p *PostServer) ReportPost(ctx context.Context, req *pb.ReportPostRequest) (*pb.ReportPostResponse, error) {
	ReportUser := req.RepostedUserid
	reportReq := models.ReportRequest{
		PostID: uint(req.Postid),
		Report: req.Report,
	}
	err := p.postUseCase.ReportPost(int(ReportUser), reportReq)
	if err != nil {
		return &pb.ReportPostResponse{}, err
	}
	return &pb.ReportPostResponse{}, nil
}

func (p *PostServer) SavedPost(ctx context.Context, req *pb.SavedPostRequest) (*pb.SavedPostResponse, error) {
	userID, postID := req.Userid, req.Postid
	err := p.postUseCase.SavedPost(int(userID), int(postID))
	if err != nil {
		return &pb.SavedPostResponse{}, err
	}
	return &pb.SavedPostResponse{}, nil
}

func (p *PostServer) UnSavedPost(ctx context.Context, req *pb.UnSavedPostRequest) (*pb.UnSavedPostResponse, error) {
	userID, postID := req.Userid, req.Postid
	err := p.postUseCase.UnSavedPost(int(userID), int(postID))
	if err != nil {
		return &pb.UnSavedPostResponse{}, err
	}
	return &pb.UnSavedPostResponse{}, nil
}

func (p *PostServer) GetSavedPost(ctx context.Context, req *pb.GetSavedPostRequest) (*pb.GetSavedPostResponse, error) {
	userID := req.Userid
	posts, err := p.postUseCase.GetSavedPost(int(userID))
	if err != nil {
		return nil, err
	}

	var allPostResponses []*pb.CreatePostResponse
	for _, post := range posts {
		userData := &pb.UserData{
			Userid:   int64(post.Author.UserId),
			Username: post.Author.Username,
			Imageurl: post.Author.Profile,
		}
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
	return &pb.GetSavedPostResponse{
		Allpost: allPostResponses,
	}, nil
}

func (p *PostServer) CreateStory(ctx context.Context, req *pb.CreateStoryRequest) (*pb.CreateStoryResponse, error) {
	userID := req.Userid
	file := req.Story.Postimages
	data, err := p.storyUseCase.CreateStory(int(userID), file)
	if err != nil {
		return nil, err
	}
	user := &pb.UserData{
		Userid:   int64(data.Author.UserId),
		Username: data.Author.Username,
		Imageurl: data.Author.Profile,
	}
	return &pb.CreateStoryResponse{
		User:      user,
		Story:     data.Story,
		CreatedAt: timestamppb.New(data.CreatedAt),
	}, nil
}

func (p *PostServer) GetStory(ctx context.Context, req *pb.GetStoryRequest) (*pb.GetStoryResponse, error) {
	userID, viewer := req.Userid, req.Viewer
	data, err := p.storyUseCase.GetStory(int(userID), int(viewer))
	if err != nil {
		return nil, err
	}

	var allStoryResponses []*pb.CreateStoryResponses
	for _, story := range data {
		userData := &pb.UserData{
			Userid:   int64(story.Author.UserId),
			Username: story.Author.Username,
			Imageurl: story.Author.Profile,
		}
		details := &pb.CreateStoryResponses{
			User:      userData,
			StoryID:   int64(story.StoryID),
			Story:     story.Story,
			CreatedAt: timestamppb.New(story.CreatedAt),
		}
		allStoryResponses = append(allStoryResponses, details)
	}
	return &pb.GetStoryResponse{
		Stories: allStoryResponses,
	}, nil
}

func (p *PostServer) DeleteStory(ctx context.Context, req *pb.DeleteStoryRequest) (*pb.DeleteStoryResponse, error) {
	userID := req.UserID
	storyID := req.Storyid
	err := p.storyUseCase.DeleteStory(int(userID), int(storyID))
	if err != nil {
		return &pb.DeleteStoryResponse{}, err
	}
	return &pb.DeleteStoryResponse{}, nil
}

func (p *PostServer) LikeStory(ctx context.Context, req *pb.LikeStoryRequest) (*pb.LikeStoryResponse, error) {
	userID := req.Userid
	storyID := req.Storyid
	err := p.storyUseCase.LikeStory(int(userID), int(storyID))
	if err != nil {
		return &pb.LikeStoryResponse{}, err
	}
	return &pb.LikeStoryResponse{}, nil
}

func (p *PostServer) UnLikeStory(ctx context.Context, req *pb.LikeStoryRequest) (*pb.LikeStoryResponse, error) {
	userID := req.Userid
	storyID := req.Storyid
	err := p.storyUseCase.UnLikeStory(int(userID), int(storyID))
	if err != nil {
		return &pb.LikeStoryResponse{}, err
	}
	return &pb.LikeStoryResponse{}, nil
}

func (p *PostServer) StoryDetails(ctx context.Context, req *pb.StoryDetailsRequest) (*pb.StoryDetailsResponse, error) {
	userID, storyID := req.UserID, req.StoryID
	data, err := p.storyUseCase.StoryDetails(int(userID), int(storyID))
	if err != nil {
		return &pb.StoryDetailsResponse{}, err
	}
	var viewers []*pb.Viewer
	for _, v := range data.Viewer {
		viewers = append(viewers, &pb.Viewer{
			Viewer:  v.ViewUser,
			Profile: v.Profile,
		})
	}

	var likedUsers []*pb.Likeduser
	for _, l := range data.LikedUser {
		likedUsers = append(likedUsers, &pb.Likeduser{
			Likeduser: l.LikeUser,
			Profile:   l.Profile,
		})
	}

	return &pb.StoryDetailsResponse{
		StoryID:   int64(data.StoryID),
		Viewer:    viewers,
		LikedUser: likedUsers,
	}, nil
}

func (p *PostServer) ShowPostReports(ctx context.Context, req *pb.ShowPostReportsRequest) (*pb.ShowPostReportsResponse, error) {
	page, count := req.Page, req.Count
	data, err := p.postUseCase.ShowPostReports(int(page), int(count))
	if err != nil {
		return &pb.ShowPostReportsResponse{}, err
	}
	var reports []*pb.ReportPostResponses
	for _, v := range data {
		report := &pb.ReportPostResponses{
			RepostedUserid: int64(v.ReportUserID),
			Postid:         int64(v.PostID),
			Report:         v.Report,
		}
		reports = append(reports, report)
	}
	return &pb.ShowPostReportsResponse{
		Reports: reports,
	}, nil
}

func (p *PostServer) GetAllPosts(ctx context.Context, req *pb.GetAllpostsRequest) (*pb.GetAllpostsResponse, error) {
	page, count := req.Page, req.Count
	data, err := p.postUseCase.GetAllPosts(int(page), int(count))
	if err != nil {
		return &pb.GetAllpostsResponse{}, err
	}
	var posts []*pb.PostResponse
	for _, v := range data {
		Users := &pb.UserData{
			Userid:   int64(v.Author.UserId),
			Username: v.Author.Username,
			Imageurl: v.Author.Profile,
		}
		details := &pb.PostResponse{
			Id:        int64(v.ID),
			User:      Users,
			Caption:   v.Caption,
			Url:       v.Url,
			Like:      int64(v.Likes),
			Comment:   int64(v.Comments),
			CreatedAt: timestamppb.New(v.CreatedAt),
		}
		posts = append(posts, details)
	}
	return &pb.GetAllpostsResponse{
		Posts: posts,
	}, nil
}

func (p *PostServer) CheckPostIDByID(ctx context.Context, req *pb.CheckPostIDByIDRequest) (*pb.CheckPostIDByIDResponse, error) {
	postID := req.PostID
	ok := p.postUseCase.CheckPostIDByID(int(postID))
	return &pb.CheckPostIDByIDResponse{
		Exist: ok,
	}, nil
}

func (p *PostServer) RemovePost(ctx context.Context, req *pb.RemovePostRequest) (*pb.RemovePostResponse, error) {
	postID := req.PostID
	err := p.postUseCase.RemovePost(int(postID))
	if err != nil {
		return &pb.RemovePostResponse{}, err
	}
	return &pb.RemovePostResponse{}, nil
}

func (p *PostServer) Home(ctx context.Context, req *pb.HomeRequest) (*pb.HomeResponse, error) {
	data, err := p.postUseCase.Home(int(req.UserID))
	if err != nil {
		return &pb.HomeResponse{}, err
	}
	var allPostResponses []*pb.CreatePostResponse
	for _, post := range data {
		userData := &pb.UserData{
			Userid:   int64(post.Author.UserId),
			Username: post.Author.Username,
			Imageurl: post.Author.Profile,
		}

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

	return &pb.HomeResponse{
		Allpost: allPostResponses,
	}, nil
}

func (p *PostServer) CreatePostType(ctx context.Context, req *pb.CreatePostTypeRequest) (*pb.CreatePostTypeResponse, error) {
	postType := req.PostType
	err := p.postUseCase.CreatePostType(postType)
	if err != nil {
		return &pb.CreatePostTypeResponse{}, err
	}
	return &pb.CreatePostTypeResponse{}, nil
}

func (p *PostServer) ShowPostType(ctx context.Context, req *pb.ShowPostTypeRequest) (*pb.ShowPostTypeResponse, error) {
	data, err := p.postUseCase.ShowPostType()
	if err != nil {
		return &pb.ShowPostTypeResponse{}, err
	}
	var result []*pb.PostTypes
	for _, v := range data {
		posttypes := &pb.PostTypes{
			Id:       int64(v.ID),
			PostType: v.Type,
		}
		result = append(result, posttypes)
	}
	return &pb.ShowPostTypeResponse{
		Types: result,
	}, nil
}

func (p *PostServer) DeletePostType(ctx context.Context, req *pb.DeletePostTypeRequest) (*pb.DeletePostTypeResponse, error) {
	postTypeID := req.PostTypeID
	err := p.postUseCase.DeletePostType(int(postTypeID))
	if err != nil {
		return &pb.DeletePostTypeResponse{}, err
	}
	return &pb.DeletePostTypeResponse{}, nil
}
