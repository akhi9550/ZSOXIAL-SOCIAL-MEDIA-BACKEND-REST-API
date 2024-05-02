package repository

import (
	"math/rand"
	"sort"
	"time"

	interfaces "github.com/akhi9550/post-svc/pkg/repository/interface"
	"github.com/akhi9550/post-svc/pkg/utils/models"
	"gorm.io/gorm"
)

type postRepository struct {
	DB *gorm.DB
}

func NewPostRepository(DB *gorm.DB) interfaces.PostRepository {
	return &postRepository{
		DB: DB,
	}
}

func (p *postRepository) CheckUserAvalilabilityWithUserID(userID int) bool {
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM users WHERE id = ?`, userID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (p *postRepository) CheckMediaAvalilabilityWithID(typeid int) bool {
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM post_types WHERE id = ?`, typeid).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (p *postRepository) CheckPostAvalilabilityWithID(postID int) bool {
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM posts WHERE id = ?`, postID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func(p *postRepository) CheckPostedUserID(userID, PostID int)bool{
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM posts WHERE id = ? AND user_id = ?`, PostID,userID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}


func (p *postRepository) CreatePost(userID int, Caption string, TypeId int, file string, users []models.Tag) (models.Response, []models.Tag, error) {
	var post models.Response
	var tag []models.Tag
	err := p.DB.Raw(`INSERT INTO posts (user_id, url, caption, type_id, created_at) VALUES (?, ?,?, ?, NOW()) RETURNING id,url, caption, likes_count, comments_count, created_at`, userID, file, Caption, TypeId).Scan(&post).Error
	if err != nil {
		return models.Response{}, []models.Tag{}, err
	}
	for _, i := range users {
		err := p.DB.Exec(`INSERT INTO tags(user_id,post_id,taguser) VALUES ( ?,?,? )`, userID, post.ID, i.User).Error
		if err != nil {
			return models.Response{}, []models.Tag{}, err
		}
	}
	err = p.DB.Raw(`SELECT taguser FROM tags WHERE post_id = ?`, post.ID).Scan(&tag).Error
	if err != nil {
		return models.Response{}, []models.Tag{}, err
	}
	return post, tag, nil
}

func (p *postRepository) UserData(userID int) (models.UserData, error) {
	var user models.UserData
	err := p.DB.Raw(`SELECT user_id,username,url FROM users WHERE id = ?`, userID).Scan(&user).Error
	if err != nil {
		return models.UserData{}, err
	}
	return user, nil
}

func (p *postRepository) GetPost(postID int) (models.Responses, error) {
	var post models.Responses
	err := p.DB.Raw(`SELECT id,url,caption,user_id,likes_count, comments_count,created_at FROM posts WHERE  id = ? AND is_archive = 'false' `, postID).Scan(&post).Error
	if err != nil {
		return models.Responses{}, err
	}
	if post.ID == 0 {
		return models.Responses{}, err
	}

	return post, nil

}

func (p *postRepository) GetTagUser(postID int) ([]models.Tag, error) {
	var tag []models.Tag
	err := p.DB.Raw(`SELECT taguser FROM tags WHERE post_id = ?`, postID).Scan(&tag).Error
	if err != nil {
		return []models.Tag{}, err
	}
	return tag, nil
}

func (ur *postRepository) UpdateCaption(userID, postID int, caption string) error {
	err := ur.DB.Exec("UPDATE posts SET caption= $1 WHERE id = $2 AND user_id = $3", caption, postID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *postRepository) UpdateTypeID(userID, postID, typeID int) error {
	err := ur.DB.Exec("UPDATE posts SET type_id= ? WHERE id = ? AND user_id = ?", typeID, postID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *postRepository) UpdateTags(userID, postID int, tag []models.Tag) error {
	err := ur.DB.Exec(`DELETE FROM tags WHERE post_id = ? AND user_id = ?`, postID, userID).Error
	if err != nil {
		return err
	}
	for _, i := range tag {
		err := ur.DB.Exec("INSERT INTO tags (user_id,post_id,taguser) VALUES (?,?,?)", userID, postID, i.User).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *postRepository) PostDetails(postID, userID int) (models.Response, []models.Tag, error) {
	var post models.Response
	var tag []models.Tag
	err := p.DB.Raw(`SELECT taguser FROM tags WHERE post_id = ?`, postID).Scan(&tag).Error
	if err != nil {
		return models.Response{}, []models.Tag{}, err
	}
	err = p.DB.Raw(`SELECT id,url,caption,likes_count, comments_count,created_at FROM posts WHERE user_id = ? AND id = ?`, userID, postID).Scan(&post).Error
	if err != nil {
		return models.Response{}, []models.Tag{}, err
	}
	return post, tag, nil
}

func (p *postRepository) DeletePost(userID, postID int) error {
	err := p.DB.Exec(`DELETE FROM posts WHERE user_id = ? AND id = ?`, userID, postID).Error
	if err != nil {
		return err
	}
	err = p.DB.Exec(`DELETE FROM tags WHERE user_id = ? AND post_id = ?`, userID, postID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) GetPostAll(userID int) ([]models.Response, error) {
	var post []models.Response
	err := p.DB.Raw(`SELECT id,url,caption,likes_count, comments_count,created_at FROM posts WHERE user_id = ?  AND is_archive = 'false'`, userID).Scan(&post).Error
	if err != nil {
		return []models.Response{}, err
	}
	return post, nil
}

func (p *postRepository) CheckPostAvalilabilityWithUserID(postID, userID int) bool {
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM posts WHERE id = ? AND user_id = ?`, postID, userID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (p *postRepository) ArchivePost(userID, postID int) error {
	err := p.DB.Exec(`INSERT INTO archive_posts (user_id,post_id) VALUES (?,?)`, userID, postID).Error
	if err != nil {
		return err
	}
	err = p.DB.Exec(`UPDATE posts SET is_archive = 'true' WHERE id = ? AND user_id = ?`, postID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) UnArchivePost(userID, postID int) error {
	err := p.DB.Exec(`DELETE FROM archive_posts WHERE user_id = ? AND post_id = ?`, userID, postID).Error
	if err != nil {
		return err
	}
	err = p.DB.Exec(`UPDATE posts SET is_archive = 'false' WHERE id = ? AND user_id = ?`, postID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) GetAllArchivePost(userID int) ([]models.ArchivePostResponse, error) {
	var response []models.ArchivePostResponse
	err := p.DB.Raw(`SELECT id,url,caption,likes_count, comments_count,created_at FROM posts WHERE user_id = ? AND is_archive = 'true'`, userID).Scan(&response).Error
	if err != nil {
		return []models.ArchivePostResponse{}, err
	}
	return response, nil
}

func (p *postRepository) CheckAlreadyLiked(userID, PostID int) bool {
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM likes WHERE post_id = ?  AND liked_user = ?`, PostID, userID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (p *postRepository) LikePost(userID, postID int) (models.LikesReponse, error) {
	var response models.LikesReponse
	err := p.DB.Raw(`INSERT INTO likes(post_id,liked_user,created_at) VALUES (?,?,NOW()) RETURNING liked_user,created_at`, postID, userID).Scan(&response).Error
	if err != nil {
		return models.LikesReponse{}, err
	}
	err = p.DB.Exec(`UPDATE posts SET likes_count = likes_count + 1 WHERE id = ?`, postID).Error
	if err != nil {
		return models.LikesReponse{}, err
	}
	return response, err
}

func (p *postRepository) GetPostedUserID(postID int) (int, error) {
	var id int
	err := p.DB.Raw(`SELECT user_id FROM posts WHERE id = ?`, postID).Scan(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *postRepository) UnLikePost(userID, postID int) error {
	err := p.DB.Exec(`UPDATE posts SET likes_count = likes_count - 1 WHERE id = ?`, postID).Error
	if err != nil {
		return err
	}
	err = p.DB.Exec(`DELETE FROM likes WHERE liked_user = ? AND post_id = ?`, userID, postID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) PostComment(userID int, data models.PostCommentReq) (models.PostComments, error) {
	var response models.PostComments
	err := p.DB.Raw(`INSERT INTO comments (post_id,commented_user,comment_data,created_at) VALUES (?,?,?,NOW()) RETURNING commented_user,comment_data,created_at`, data.PostID, userID, data.Comment).Scan(&response).Error
	if err != nil {
		return models.PostComments{}, err
	}
	err = p.DB.Exec(`UPDATE posts SET comments_count = comments_count + 1 WHERE id = ?`, data.PostID).Error
	if err != nil {
		return models.PostComments{}, err
	}
	return response, nil
}

func (p *postRepository) CheckUserWithUserID(userID int) bool {
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM comments WHERE commented_user = ?`, userID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (p *postRepository) CheckCommentWithID(CommentID int) bool {
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM comments WHERE id = ? `, CommentID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (p *postRepository) DeleteComment(userID, CommentID int) error {
	err := p.DB.Exec(`DELETE FROM comments WHERE commented_user = ? AND id = ?`, userID, CommentID).Error
	if err != nil {
		return err
	}
	err = p.DB.Exec(`DELETE FROM comment_repies WHERE comment_id = ?`, CommentID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) GetAllPostComments(PostID int) ([]models.PostCommentResponses, error) {
	var comments []models.PostCommentResponses
	err := p.DB.Raw(`SELECT id,commented_user,comment_data,created_at FROM comments WHERE post_id = ?`, PostID).Scan(&comments).Error
	if err != nil {
		return []models.PostCommentResponses{}, err
	}
	return comments, nil
}

func (p *postRepository) AllReadyExistReply(userID, CommentID int) bool {
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM comment_replies WHERE id = ? AND reply_user = ?`, CommentID, userID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (p *postRepository) ReplyComment(userID int, req models.ReplyCommentReq) (models.PostComments, models.ReplyResponse, error) {
	var a models.Resp
	var comments models.PostComments
	var reply models.ReplyResponse
	err := p.DB.Raw(`SELECT commented_user,post_id FROM comments WHERE id = ?`, req.CommentID).Scan(&a).Error
	if err != nil {
		return models.PostComments{}, models.ReplyResponse{}, err
	}
	err = p.DB.Raw(`INSERT INTO comment_replies (post_id, comment_id,commented_user, reply_user, replies, created_at)
	VALUES (?, ?, ?, ?, ?, NOW())
	RETURNING reply_user, replies, created_at`, a.PostID, req.CommentID, a.UserID, userID, req.Reply).Scan(&reply).Error
	if err != nil {
		return models.PostComments{}, models.ReplyResponse{}, err
	}
	err = p.DB.Raw(`SELECT commented_user,comment_data,created_at FROM comments WHERE id = ?`, req.CommentID).Scan(&comments).Error
	if err != nil {
		return models.PostComments{}, models.ReplyResponse{}, err
	}
	return comments, reply, nil
}

func (p *postRepository) AlreadyReported(userID, postID int) bool {
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM post_reports WHERE report_user_id = ? AND post_id = ?`, userID, postID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (ur *postRepository) ReportPost(userID int, req models.ReportRequest) error {
	err := ur.DB.Exec(`INSERT INTO post_reports (report_user_id,post_id,report) VALUES (?,?,?)`, userID, req.PostID, req.Report).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) SavedPost(userID, postID int) error {
	err := p.DB.Exec(`INSERT INTO saved_posts (post_id,user_id) VALUES (?,?)`, postID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) AllReadyExistPost(userID, postID int) bool {
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM saved_posts WHERE post_id = ? AND user_id = ?`, postID, userID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (p *postRepository) UnSavedPost(userID, postID int) error {
	err := p.DB.Exec(`DELETE FROM saved_posts WHERE post_id = ? AND user_id = ?`, postID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) GetSavedPost(userID int) ([]models.SavedResponse, error) {
	var response []models.SavedResponse
	var id []models.PostID
	err := p.DB.Raw(`SELECT post_id FROM saved_posts WHERE user_id = ?`, userID).Scan(&id).Error
	if err != nil {
		return nil, err
	}
	for _, i := range id {
		var post models.SavedResponse
		err = p.DB.Raw(`SELECT id,user_id,url,caption,likes_count, comments_count,created_at FROM posts WHERE  is_archive = 'false' AND id = ?`, i.PostID).Scan(&post).Error
		if err != nil {
			return nil, err
		}
		response = append(response, post)
	}
	return response, nil
}

func (p *postRepository) GetCommentsByPostID(postID int) ([]models.AllComments, error) {
	var response []models.AllComments
	err := p.DB.Raw(`SELECT id,commented_user,comment_data,created_at FROM comments WHERE post_id = ?`, postID).Scan(&response).Error
	if err != nil {
		return []models.AllComments{}, err
	}
	return response, nil
}

func (p *postRepository) GetRepliesByID(PostID, CommentID int) ([]models.Replies, error) {
	var response []models.Replies
	err := p.DB.Raw(`SELECT reply_user,replies,created_at FROM comment_replies WHERE post_id = ? AND comment_id = ?`, PostID, CommentID).Scan(&response).Error
	if err != nil {
		return []models.Replies{}, err
	}
	return response, nil
}

func (p *postRepository) ShowPostReports(page, count int) ([]models.PostReports, error) {
	var response []models.PostReports
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * count
	err := p.DB.Raw(`SELECT report_user_id,post_id,report FROM post_reports  limit ? offset ?`, count, offset).Scan(&response).Error
	if err != nil {
		return []models.PostReports{}, err
	}
	return response, nil
}

func (p *postRepository) GetAllPosts(page, count int) ([]models.Responses, error) {
	var response []models.Responses
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * count
	err := p.DB.Raw(`SELECT id,url,caption,user_id,likes_count, comments_count,created_at FROM posts WHERE is_archive = 'false'  limit ? offset ?`, count, offset).Scan(&response).Error
	if err != nil {
		return []models.Responses{}, err
	}
	return response, nil
}

func (p *postRepository) CheckPostIDByID(postID int) bool {
	var count int
	err := p.DB.Raw(`SELECT COUNT(*) FROM posts WHERE id = ?`, postID).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}

func (p *postRepository) RemovePost(postID int) error {
	err := p.DB.Exec(`DELETE FROM posts WHERE id = ?`, postID).Error
	if err != nil {
		return err
	}
	err = p.DB.Exec(`DELETE FROM post_reports WHERE post_id = ?`, postID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) Home(users []models.Users) ([]models.Responses, error) {
	var allPosts []models.Responses
	for _, user := range users {
		var userPosts []models.Responses
		err := p.DB.Raw(`SELECT id, url, caption, user_id, likes_count, comments_count, created_at 
						FROM posts 
						WHERE is_archive = 'false' AND user_id = ? AND created_at IS NOT NULL
						ORDER BY created_at DESC`, user.FollowingUser).Scan(&userPosts).Error
		if err != nil {
			return nil, err
		}
		allPosts = append(allPosts, userPosts...)
	}

	sort.SliceStable(allPosts, func(i, j int) bool {
		return allPosts[i].CreatedAt.After(allPosts[j].CreatedAt)
	})

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(users), func(i, j int) {
		users[i], users[j] = users[j], users[i]
	})

	var responses []models.Responses
	for _, user := range users {
		for _, post := range allPosts {
			if post.UserID == uint(user.FollowingUser) {
				responses = append(responses, post)
			}
		}
	}

	return responses, nil
}
