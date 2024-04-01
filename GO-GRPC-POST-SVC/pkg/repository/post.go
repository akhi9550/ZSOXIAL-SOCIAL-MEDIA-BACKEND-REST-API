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

func (p *postRepository) CreatePost(userID int, Caption string, TypeId int, file string, users models.Tags) (models.Response, models.Tags, error) {
	var post models.Response
	var tag models.Tags
	err := p.DB.Raw(`INSERT INTO posts (user_id, url, caption, type_id, created_at) VALUES (?, ?,?, ?, NOW()) RETURNING id,url, caption, likes_count, comments_count, created_at`, userID, file, Caption, TypeId).Scan(&post).Error
	if err != nil {
		return models.Response{}, models.Tags{}, err
	}
	err = p.DB.Exec(`INSERT INTO tags(user_id,post_id,user1,user2,user3,user4,user5) VALUES ( ?,?,?,?,?,?,? )`, userID, post.ID, users.User1, users.User2, users.User3, users.User4, users.User5).Error
	if err != nil {
		return models.Response{}, models.Tags{}, err
	}
	err = p.DB.Raw(`SELECT user1, user2, user3, user4, user5 FROM tags WHERE post_id = ?`, post.ID).Scan(&tag).Error
	if err != nil {
		return models.Response{}, models.Tags{}, err
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
func (p *postRepository) GetPost(userID, postID int) (models.Response, models.Tags, error) {
	var post models.Response
	var tag models.Tags
	err := p.DB.Raw(`SELECT user1,user2,user3,user4,user5 FROM tags WHERE post_id = ?`, postID).Scan(&tag).Error
	if err != nil {
		return models.Response{}, models.Tags{}, err
	}
	err = p.DB.Raw(`SELECT id,url,caption,likes_count, comments_count,created_at FROM posts WHERE user_id = ? AND id = ?`, userID, postID).Scan(&post).Error
	if err != nil {
		return models.Response{}, models.Tags{}, err
	}
	return post, tag, nil

}

func (ur *postRepository) UpdateCaption(postID, userID int, caption string) error {
	err := ur.DB.Exec("UPDATE posts SET caption= ? WHERE id = ? AND user_id = ?", caption, postID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *postRepository) UpdateTypeID(postID, userID, typeID int) error {
	err := ur.DB.Exec("UPDATE posts SET type_id= ? WHERE id = ? AND user_id = ?", typeID, postID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (ur *postRepository) UpdateTags(postID, userID int, tag []models.Tag) error {
	for i := range tag {
		err := ur.DB.Exec("UPDATE tags SET taguser = ? WHERE post_id = ? AND user_id = ?", i, postID, userID).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *postRepository) PostDetails(userID, postID int) (models.UpdateResponse, error) {
	var post models.UpdateResponse
	var tag []models.Tag
	err := p.DB.Raw(`SELECT taguser FROM tags WHERE post_id = ?`, postID).Scan(&tag).Error
	if err != nil {
		return models.UpdateResponse{}, err
	}
	err = p.DB.Raw(`SELECT id,url,caption,likes_count, comments_count,created_at FROM posts WHERE user_id = ? AND id = ?`, userID, postID).Scan(&post).Error
	if err != nil {
		return models.UpdateResponse{}, err
	}
	post.Tag = tag
	return post, nil

}

func (p *postRepository) DeletePost(userID, postID int) error {
	err := p.DB.Exec(`DELETE * FROM posts WHERE user_id = ? AND post_id = ?`, userID, postID).Error
	if err != nil {
		return err
	}
	err = p.DB.Exec(`DELETE * FROM tags WHERE user_id = ? AND post_id = ?`, userID, postID).Error
	if err != nil {
		return err
	}
	return nil
}
