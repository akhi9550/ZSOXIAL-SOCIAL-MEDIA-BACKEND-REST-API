package usecase

import (
	"errors"

	authclientinterfaces "github.com/akhi9550/post-svc/pkg/client/interface"
	"github.com/akhi9550/post-svc/pkg/helper"
	interfaces "github.com/akhi9550/post-svc/pkg/repository/interface"
	services "github.com/akhi9550/post-svc/pkg/usecase/interface"
	"github.com/akhi9550/post-svc/pkg/utils/models"
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
	filename := "posted"
	url, err := helper.AddImageToAwsS3(file, filename)
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
func(p *postUseCase)GetAllPost(userID int) ([]models.PostResponse, error){
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(int(userID))
	if !userExist {
		return models.PostResponse{}, errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithID(postID)
	if !ok {
		return models.PostResponse{}, errors.New("post doesn't exist")
	}
	post, tag, err := p.postRepository.GetPost(userID)
	if err != nil {
		return []models.PostResponse{}, err
	}
	username, err := p.authClient.GetUserNameWithTagUserID(tag)
	if err != nil {
		return []models.PostResponse{}, err
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
		return []models.PostResponse{}, err
	}
	return []models.PostResponse{
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
	post, tag, err := p.postRepository.GetPost(userID, postID)
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
	userData, err := p.authClient.UserData(userID)
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
