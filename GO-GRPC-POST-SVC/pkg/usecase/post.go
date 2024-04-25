package usecase

import (
	"errors"
	"fmt"

	authclientinterfaces "github.com/akhi9550/post-svc/pkg/client/interface"
	"github.com/akhi9550/post-svc/pkg/helper"
	interfaces "github.com/akhi9550/post-svc/pkg/repository/interface"
	services "github.com/akhi9550/post-svc/pkg/usecase/interface"
	"github.com/akhi9550/post-svc/pkg/utils/models"
	"github.com/google/uuid"
)

type postUseCase struct {
	postRepository interfaces.PostRepository
	authClient     authclientinterfaces.NewauthClient
}

func NewPostUseCase(repository interfaces.PostRepository, authclient authclientinterfaces.NewauthClient) services.PostUseCase {
	return &postUseCase{
		postRepository: repository,
		authClient:     authclient,
	}

}

func (p *postUseCase) CreatePost(userID int, data models.PostRequest, file []byte, users []models.Tag) (models.PostResponse, error) {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(int(userID))
	if !userExist {
		return models.PostResponse{}, errors.New("user doesn't exist")
	}
	mediatype := p.postRepository.CheckMediaAvalilabilityWithID(int(data.TypeId))
	if !mediatype {
		return models.PostResponse{}, errors.New("type doesn't exist")
	}
	fileUID := uuid.New()
	fileName := fileUID.String()
	s3Path := helper.Formated(int(data.TypeId), fileName)
	url, err := helper.AddImageToAwsS3(file, s3Path)
	if err != nil {
		return models.PostResponse{}, err
	}
	usersExist := p.authClient.CheckUserAvalilabilityWithTagUserID(users)
	if !usersExist {
		return models.PostResponse{}, errors.New("user doesn't exist")
	}
	post, tag, err := p.postRepository.CreatePost(userID, data.Caption, int(data.TypeId), url, users)
	if err != nil {
		return models.PostResponse{}, err
	}
	username, err := p.authClient.GetUserNameWithTagUserID(tag)
	if err != nil {
		return models.PostResponse{}, err
	}
	var Users []models.Tag
	for _, user := range username {
		tag := models.Tag{
			User: user.User,
		}
		Users = append(Users, tag)
	}
	userData, err := p.authClient.UserData(int(userID))
	if err != nil {
		return models.PostResponse{}, err
	}
	return models.PostResponse{
		ID:        post.ID,
		Author:    userData,
		Caption:   post.Caption,
		Tag:       Users,
		Url:       post.Url,
		Likes:     post.Likes,
		Comments:  post.Comments,
		CreatedAt: post.CreatedAt,
	}, nil
}

func (p *postUseCase) GetPost(userID int, postID int) (models.PostResponse, error) {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(int(userID))
	if !userExist {
		return models.PostResponse{}, errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithID(postID)
	if !ok {
		return models.PostResponse{}, errors.New("post doesn't exist")
	}
	post, err := p.postRepository.GetPost(postID)
	if err != nil {
		return models.PostResponse{}, err
	}
	if post.ID == 0 {
		return models.PostResponse{}, errors.New(" post archived")
	}
	tag, err := p.postRepository.GetTagUser(postID)
	if err != nil {
		return models.PostResponse{}, err
	}
	username, err := p.authClient.GetUserNameWithTagUserID(tag)
	if err != nil {
		return models.PostResponse{}, err
	}
	var Users []models.Tag
	for _, user := range username {
		tag := models.Tag{
			User: user.User,
		}
		Users = append(Users, tag)
	}
	userData, err := p.authClient.UserData(int(post.UserID))
	if err != nil {
		return models.PostResponse{}, err
	}
	return models.PostResponse{
		ID:        post.ID,
		Author:    userData,
		Caption:   post.Caption,
		Tag:       Users,
		Url:       post.Url,
		Likes:     post.Likes,
		Comments:  post.Comments,
		CreatedAt: post.CreatedAt,
	}, nil
}

func (p *postUseCase) UpdatePost(userID int, data models.UpdatePostReq, tag []models.Tag) (models.UpdateResponse, error) {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return models.UpdateResponse{}, errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithID(int(data.PostID))
	if !ok {
		return models.UpdateResponse{}, errors.New("post doesn't exist")
	}
	err := p.postRepository.UpdateCaption(userID, int(data.PostID), data.Caption)
	if err != nil {
		return models.UpdateResponse{}, err
	}
	mediatype := p.postRepository.CheckMediaAvalilabilityWithID(int(data.TypeID))
	if !mediatype {
		return models.UpdateResponse{}, errors.New("type doesn't exist")
	}
	err = p.postRepository.UpdateTypeID(userID, int(data.PostID), int(data.TypeID))
	if err != nil {
		return models.UpdateResponse{}, err
	}
	usersExist := p.authClient.CheckUserAvalilabilityWithTagUserID(tag)
	if !usersExist {
		return models.UpdateResponse{}, errors.New("users doesn't exist")
	}
	err = p.postRepository.UpdateTags(userID, int(data.PostID), tag)
	if err != nil {
		return models.UpdateResponse{}, err
	}
	post, tags, err := p.postRepository.PostDetails(int(data.PostID), userID)
	if err != nil {
		return models.UpdateResponse{}, err
	}
	username, err := p.authClient.GetUserNameWithTagUserID(tags)
	if err != nil {
		return models.UpdateResponse{}, err
	}
	var Users []models.Tag
	for _, user := range username {
		tag := models.Tag{
			User: user.User,
		}
		Users = append(Users, tag)
	}
	userData, err := p.authClient.UserData(userID)
	if err != nil {
		return models.UpdateResponse{}, err
	}
	return models.UpdateResponse{
		ID:        post.ID,
		Author:    userData,
		Caption:   post.Caption,
		Tag:       Users,
		Url:       post.Url,
		Likes:     post.Likes,
		Comments:  post.Comments,
		CreatedAt: post.CreatedAt,
	}, nil
}

func (p *postUseCase) DeletePost(userID int, postID int) error {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithID(postID)
	if !ok {
		return errors.New("post doesn't exist")
	}
	err := p.postRepository.DeletePost(userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (p *postUseCase) GetAllPost(userID int) ([]models.PostResponse, error) {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(int(userID))
	if !userExist {
		return []models.PostResponse{}, errors.New("user doesn't exist")
	}
	post, err := p.postRepository.GetPostAll(userID)
	if err != nil {
		return []models.PostResponse{}, err
	}

	userData, err := p.authClient.UserData(userID)
	if err != nil {
		return []models.PostResponse{}, err
	}
	var postResponses []models.PostResponse
	for _, singlePost := range post {
		postResponses = append(postResponses, models.PostResponse{
			ID:        singlePost.ID,
			Author:    userData,
			Caption:   singlePost.Caption,
			Url:       singlePost.Url,
			Likes:     singlePost.Likes,
			Comments:  singlePost.Comments,
			CreatedAt: singlePost.CreatedAt,
		})
	}

	return postResponses, nil
}

func (p *postUseCase) ArchivePost(userID, postID int) error {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithUserID(postID, userID)
	if !ok {
		return errors.New("post doesn't exist")
	}
	err := p.postRepository.ArchivePost(userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (p *postUseCase) UnArchivePost(userID, postID int) error {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithUserID(postID, userID)
	if !ok {
		return errors.New("post doesn't exist")
	}
	err := p.postRepository.UnArchivePost(userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (p *postUseCase) GetAllArchivePost(userID int) ([]models.ArchivePostResponse, error) {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(int(userID))
	if !userExist {
		return []models.ArchivePostResponse{}, errors.New("user doesn't exist")
	}
	post, err := p.postRepository.GetAllArchivePost(userID)
	if err != nil {
		return []models.ArchivePostResponse{}, err
	}
	var postResponses []models.ArchivePostResponse
	for _, singlePost := range post {
		postResponses = append(postResponses, models.ArchivePostResponse{
			ID:        singlePost.ID,
			Caption:   singlePost.Caption,
			Url:       singlePost.Url,
			Likes:     singlePost.Likes,
			Comments:  singlePost.Comments,
			CreatedAt: singlePost.CreatedAt,
		})
	}

	return postResponses, nil
}

func (p *postUseCase) LikePost(userID int, postID int) (models.LikePostResponse, error) {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return models.LikePostResponse{}, errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithID(postID)
	if !ok {
		return models.LikePostResponse{}, errors.New("post doesn't exist")
	}
	ok = p.postRepository.CheckAlreadyLiked(userID, postID)
	if ok {
		return models.LikePostResponse{}, errors.New("already liked")
	}
	data, err := p.postRepository.LikePost(userID, postID)
	if err != nil {
		return models.LikePostResponse{}, err
	}
	userData, err := p.authClient.UserData(userID)
	if err != nil {
		return models.LikePostResponse{}, err
	}
	postedUserID, err := p.postRepository.GetPostedUserID(postID)
	if err != nil {
		return models.LikePostResponse{}, err
	}
	msg := fmt.Sprintf("%s Liked PostID %d", userData.Username, postID)
	helper.SendLikeNotification(models.Notification{
		UserID:      postedUserID,
		LikedUserID: userID,
		PostID:      postID,
	}, []byte(msg))

	return models.LikePostResponse{
		UserID:    data.UserID,
		LikedUser: userData.Username,
		Profile:   userData.Profile,
		CreatedAt: data.CreatedAt,
	}, nil
}

func (p *postUseCase) UnLinkPost(userID int, postID int) error {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithID(postID)
	if !ok {
		return errors.New("post doesn't exist")
	}
	ok = p.postRepository.CheckAlreadyLiked(userID, postID)
	if !ok {
		p.postRepository.LikePost(userID, postID)
		return errors.New("")
	}
	err := p.postRepository.UnLikePost(userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (p *postUseCase) PostComment(userID int, data models.PostCommentReq) (models.PostComment, error) {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return models.PostComment{}, errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithID(int(data.PostID))
	if !ok {
		return models.PostComment{}, errors.New("post doesn't exist")
	}
	result, err := p.postRepository.PostComment(userID, data)
	if err != nil {
		return models.PostComment{}, err
	}
	userData, err := p.authClient.UserData(userID)
	if err != nil {
		return models.PostComment{}, err
	}
	postedUserID, err := p.postRepository.GetPostedUserID(int(data.PostID))
	if err != nil {
		return models.PostComment{}, err
	}

	msg := fmt.Sprintf("%s comment on post %d Comment %s", userData.Username, data.PostID, data.Comment)
	helper.SendCommentNotification(models.Notification{
		UserID:      postedUserID,
		LikedUserID: userID,
		PostID:      int(data.PostID),
	}, []byte(msg))

	return models.PostComment{
		UserID:        result.UserID,
		CommentedUser: userData.Username,
		Profile:       userData.Profile,
		Comment:       result.Comment,
		CreatedAt:     result.CreatedAt,
	}, nil
}

func (p *postUseCase) DeleteComment(userID, CommentID int) error {
	userExist := p.postRepository.CheckUserWithUserID(userID)
	if !userExist {
		return errors.New("user couldn't delete this")
	}
	ok := p.postRepository.CheckCommentWithID(CommentID)
	if !ok {
		return errors.New("comment doesn't exist")
	}
	err := p.postRepository.DeleteComment(userID, CommentID)
	if err != nil {
		return err
	}
	return nil
}
func (p *postUseCase) GetAllPostComments(PostID int) ([]models.PostCommentResponse, error) {
	postExist := p.postRepository.CheckPostAvalilabilityWithID(PostID)
	if !postExist {
		return []models.PostCommentResponse{}, errors.New("post doesn't exist")
	}
	data, err := p.postRepository.GetAllPostComments(PostID)
	if err != nil {
		return []models.PostCommentResponse{}, err
	}
	var comments []models.PostCommentResponse
	for _, post := range data {
		userData, err := p.authClient.UserData(int(post.UserID))
		if err != nil {
			return nil, err
		}
		comment := models.PostCommentResponse{
			UserID:        userData.UserId,
			CommentedUser: userData.Username,
			Profile:       userData.Profile,
			CommentID:     post.CommentID,
			Comment:       post.Comment,
			CreatedAt:     post.CreatedAt,
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (p *postUseCase) ReplyComment(userID int, req models.ReplyCommentReq) (models.ReplyReposne, error) {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return models.ReplyReposne{}, errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckCommentWithID(int(req.CommentID))
	if !ok {
		return models.ReplyReposne{}, errors.New("comment doesn't exist")
	}
	// alreadyReplied := p.postRepository.AllReadyExistReply(userID, int(req.CommentID))
	// if alreadyReplied {
	// 	return models.ReplyReposne{}, errors.New("already replied the comment")
	// }
	com, rep, err := p.postRepository.ReplyComment(userID, req)
	if err != nil {
		return models.ReplyReposne{}, err
	}
	commetUserData, err := p.authClient.UserData(int(com.UserID))
	if err != nil {
		return models.ReplyReposne{}, err
	}
	replyUserData, err := p.authClient.UserData(int(rep.UserID))
	if err != nil {
		return models.ReplyReposne{}, err
	}
	comment := models.PostComment{
		UserID:        com.UserID,
		CommentedUser: commetUserData.Username,
		Profile:       commetUserData.Profile,
		Comment:       com.Comment,
		CreatedAt:     com.CreatedAt,
	}
	reply := models.ReplyPostCommentResponse{
		UserID:    rep.UserID,
		ReplyUser: replyUserData.Username,
		Profile:   replyUserData.Profile,
		Reply:     rep.Reply,
		CreatedAt: rep.CreatedAt,
	}

	return models.ReplyReposne{
		Comment: comment,
		Reply:   reply,
	}, nil
}

func (p *postUseCase) ShowAllPostComments(PostID int) ([]models.AllCommentsAndReplies, error) {
	postExist := p.postRepository.CheckPostAvalilabilityWithID(PostID)
	if !postExist {
		return []models.AllCommentsAndReplies{}, errors.New("post doesn't exist")
	}
	comments, err := p.postRepository.GetCommentsByPostID(PostID)
	if err != nil {
		return []models.AllCommentsAndReplies{}, err
	}
	var Allcomments []models.AllCommentsAndReplies
	for _, comment := range comments {
		userData, err := p.authClient.UserData(int(comment.UserID))
		if err != nil {
			return nil, err
		}
		Reply, err := p.postRepository.GetRepliesByID(PostID, int(comment.CommentID))
		if err != nil {
			return nil, err
		}

		var replies []models.AllReplies
		for _, reply := range Reply {
			ReplyuserData, err := p.authClient.UserData(int(reply.UserID))
			if err != nil {
				return nil, err
			}

			repliess := models.AllReplies{
				UserID:    ReplyuserData.UserId,
				ReplyUser: ReplyuserData.Username,
				Profile:   ReplyuserData.Profile,
				Reply:     reply.Reply,
				CreatedAt: reply.CreatedAt,
			}
			replies = append(replies, repliess)
		}
		commentWithReplies := models.AllCommentsAndReplies{
			CommentUser: userData.Username,
			Profile:     userData.Profile,
			Comment:     comment.Comment,
			CreatedAt:   comment.CreatedAt,
			Reply:       replies,
		}
		Allcomments = append(Allcomments, commentWithReplies)
	}
	return Allcomments, nil
}

func (p *postUseCase) ReportPost(userID int, req models.ReportRequest) error {
	ReportuserExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !ReportuserExist {
		return errors.New("user doesn't exist")
	}
	postExist := p.postRepository.CheckPostAvalilabilityWithID(int(req.PostID))
	if !postExist {
		return errors.New("post doesn't exist")
	}
	Isreport := p.postRepository.AlreadyReported(userID, int(req.PostID))
	if Isreport {
		return errors.New("already reported")
	}
	err := p.postRepository.ReportPost(userID, req)
	if err != nil {
		return err
	}
	return nil
}

func (p *postUseCase) SavedPost(userID, postID int) error {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithID(postID)
	if !ok {
		return errors.New("post doesn't exist")
	}
	saved := p.postRepository.AllReadyExistPost(userID, postID)
	if saved {
		return errors.New("already saved")
	}
	err := p.postRepository.SavedPost(userID, postID)
	if err != nil {
		return err
	}
	return nil
}

func (p *postUseCase) UnSavedPost(userID, postID int) error {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithID(postID)
	if !ok {
		return errors.New("post doesn't exist")
	}

	err := p.postRepository.UnSavedPost(userID, postID)
	if err != nil {
		return err
	}
	return nil
}
func (p *postUseCase) GetSavedPost(userID int) ([]models.PostResponse, error) {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(int(userID))
	if !userExist {
		return []models.PostResponse{}, errors.New("user doesn't exist")
	}
	posts, err := p.postRepository.GetSavedPost(userID)
	if err != nil {
		return nil, err
	}
	var postResponses []models.PostResponse
	for _, post := range posts {
		userData, err := p.authClient.UserData(int(post.UserID))
		if err != nil {
			return nil, err
		}
		postResponse := models.PostResponse{
			ID:        post.ID,
			Author:    userData,
			Caption:   post.Caption,
			Url:       post.Url,
			Likes:     post.Likes,
			Comments:  post.Comments,
			CreatedAt: post.CreatedAt,
		}
		postResponses = append(postResponses, postResponse)
	}
	return postResponses, nil
}

func (p *postUseCase) ShowPostReports(page, count int) ([]models.PostReports, error) {
	reports, err := p.postRepository.ShowPostReports(page, count)
	if err != nil {
		return []models.PostReports{}, err
	}
	return reports, nil
}

func (p *postUseCase) GetAllPosts(page, count int) ([]models.PostResponse, error) {
	posts, err := p.postRepository.GetAllPosts(page, count)
	if err != nil {
		return []models.PostResponse{}, err
	}
	var AllPosts []models.PostResponse
	for _, v := range posts {
		userData, err := p.authClient.UserData(int(v.UserID))
		if err != nil {
			return []models.PostResponse{}, err
		}
		detaiils := models.PostResponse{
			ID:        v.ID,
			Author:    userData,
			Caption:   v.Caption,
			Url:       v.Url,
			Likes:     v.Likes,
			Comments:  v.Comments,
			CreatedAt: v.CreatedAt,
		}
		AllPosts = append(AllPosts, detaiils)
	}
	return AllPosts, nil
}
func (p *postUseCase) CheckPostIDByID(postID int) bool {
	ok := p.postRepository.CheckPostIDByID(postID)
	return ok
}

func (p *postUseCase) RemovePost(postID int) error {
	postExist := p.postRepository.CheckPostIDByID(postID)
	if !postExist {
		return errors.New("postID doesn't exist")
	}
	err := p.postRepository.RemovePost(postID)
	if err != nil {
		return err
	}
	return nil
}
