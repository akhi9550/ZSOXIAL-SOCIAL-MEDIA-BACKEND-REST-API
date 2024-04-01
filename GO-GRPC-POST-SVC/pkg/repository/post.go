package repository

import (
	"fmt"

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

func (p *postRepository) CreatePost(userID int, Caption string, TypeId int, file string, users models.Tags) (models.Response, models.Tags, []models.Url, error) {
	var post models.Response
	var tag models.Tags
	var image []models.Url
	err := p.DB.Raw(`INSERT INTO posts (user_id, caption, type_id, created_at) VALUES (?, ?, ?, NOW()) RETURNING id, caption, likes_count, comments_count, created_at`, userID, Caption, TypeId).Scan(&post).Error
	if err != nil {
		return models.Response{}, models.Tags{}, []models.Url{}, err
	}
	err = p.DB.Exec(`INSERT INTO tags(user_id,post_id,user1,user2,user3,user4,user5) VALUES ( ?,?,?,?,?,?,? )`, userID, post.ID, users.User1, users.User2, users.User3, users.User4, users.User5).Error
	if err != nil {
		return models.Response{}, models.Tags{}, []models.Url{}, err
	}
	err = p.DB.Exec(`INSERT INTO urls(user_id,post_id,url) VALUES ( ?,?,? )`, userID, post.ID, file).Error
	if err != nil {
		return models.Response{}, models.Tags{}, []models.Url{}, err
	}
	err = p.DB.Raw(`SELECT url FROM urls WHERE post_id = ? AND user_id = ?`, post.ID, userID).Scan(&image).Error
	if err != nil {
		return models.Response{}, models.Tags{}, []models.Url{}, err
	}
	err = p.DB.Raw(`SELECT user1, user2, user3, user4, user5 FROM tags WHERE post_id = ?`, post.ID).Scan(&tag).Error
	if err != nil {
		return models.Response{}, models.Tags{}, []models.Url{}, err
	}
	fmt.Println("imagessss:=", image)
	return post, tag, image, nil
}
func (p *postRepository) UserData(userID int) (models.UserData, error) {
	var user models.UserData
	err := p.DB.Raw(`SELECT user_id,username,url FROM users WHERE id = ?`, userID).Scan(&user).Error
	if err != nil {
		return models.UserData{}, err
	}
	return user, nil
}
func (p *postRepository) GetPost(userID, postID int) (models.Response, models.Tags, []models.Url, error) {
	var post models.Response
	var tag models.Tags
	var image []models.Url
	err := p.DB.Raw(`SELECT user1,user2,user3,user4,user5 FROM tags WHERE post_id = ?`, postID).Scan(&tag).Error
	if err != nil {
		return models.Response{}, models.Tags{}, []models.Url{}, err
	}
	err = p.DB.Raw(`SELECT id,caption,likes_count, comments_count,created_at FROM posts WHERE user_id = ? AND id = ?`, userID, postID).Scan(&post).Error
	if err != nil {
		return models.Response{}, models.Tags{}, []models.Url{}, err
	}
	err = p.DB.Raw(`SELECT url FROM urls WHERE post_id = ? AND user_id = ?`, postID, userID).Scan(&image).Error
	if err != nil {
		return models.Response{}, models.Tags{}, []models.Url{}, err
	}
	return post, tag, image, nil

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

func (ur *postRepository) UpdateTags(postID, userID int, data models.Tags) error {
	err := ur.DB.Exec("UPDATE tags SET user1,user2,user3,user4,user5 WHERE post_id = ? AND user_id = ?", data.User1, data.User2, data.User3, data.User4, data.User5, postID, userID).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *postRepository) PostDetails(userID, postID int) (models.Response, models.Tags, []models.Url, error) {
	var post models.Response
	var tag models.Tags
	var image []models.Url
	err := p.DB.Raw(`SELECT user1,user2,user3,user4,user5 FROM tags WHERE post_id = ?`, postID).Scan(&tag).Error
	if err != nil {
		return models.Response{}, models.Tags{}, []models.Url{}, err
	}
	err = p.DB.Raw(`SELECT id,caption,likes_count, comments_count,created_at FROM posts WHERE user_id = ? AND id = ?`, userID, postID).Scan(&post).Error
	if err != nil {
		return models.Response{}, models.Tags{}, []models.Url{}, err
	}
	err = p.DB.Raw(`SELECT url FROM urls  WHERE post_id = ? AND user_id = ?`, postID, userID).Scan(&image).Error
	if err != nil {
		return models.Response{}, models.Tags{}, []models.Url{}, err
	}
	return post, tag, image, nil

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
