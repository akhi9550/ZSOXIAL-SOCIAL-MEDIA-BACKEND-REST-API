package repository

import (
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
func (p *postRepository) GetPost(userID, postID int) (models.Response, []models.Tag, error) {
	var post models.Response
	var tag []models.Tag
	err := p.DB.Raw(`SELECT taguser FROM tags WHERE post_id = ?`, postID).Scan(&tag).Error
	if err != nil {
		return models.Response{}, []models.Tag{}, err
	}
	err = p.DB.Raw(`SELECT id,url,caption,likes_count, comments_count,created_at FROM posts WHERE user_id = ? AND id = ? AND is_archive = 'false' `, userID, postID).Scan(&post).Error
	if err != nil {
		return models.Response{}, []models.Tag{}, err
	}
	return post, tag, nil

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

func (p *postRepository) PostComment(userID int, data models.PostCommentReq) (models.PostCommentResponses, error) {
	var response models.PostCommentResponses
	err := p.DB.Raw(`INSERT INTO comments (post_id,commented_user,comment_data,created_at) VALUES (?,?,?,NOW()) RETURNING commented_user,comment_data,created_at`, data.PostID, userID, data.Comment).Scan(&response).Error
	if err != nil {
		return models.PostCommentResponses{}, err
	}
	err = p.DB.Exec(`UPDATE posts SET comments_count = comments_count + 1 WHERE id = ?`, data.PostID).Error
	if err != nil {
		return models.PostCommentResponses{}, err
	}
	return response, nil
}
