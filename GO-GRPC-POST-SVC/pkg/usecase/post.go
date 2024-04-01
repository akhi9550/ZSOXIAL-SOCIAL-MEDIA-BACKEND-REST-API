package usecase

import (
	"errors"
	"fmt"

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

func (p *postUseCase) CreatePost(userID int, data models.PostRequest, file []byte, users models.Tags) (models.PostResponse, error) {
	fmt.Println("userid", userID)
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(int(userID))
	fmt.Println("userid", userExist)
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
	post, tag, image, err := p.postRepository.CreatePost(userID, data.Caption, int(data.TypeId), url, users)
	if err != nil {
		return models.PostResponse{}, err
	}
	userData, err := p.authClient.UserData(int(userID))
	if err != nil {
		return models.PostResponse{}, err
	}
	return models.PostResponse{
		ID:        post.ID,
		Author:    userData,
		Caption:   post.Caption,
		Tag:       tag,
		ImageUrls: image,
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
	var image []models.Url
	post, tag, image, err := p.postRepository.GetPost(userID, postID)
	if err != nil {
		return models.PostResponse{}, err
	}
	userData, err := p.authClient.UserData(userID)
	if err != nil {
		return models.PostResponse{}, err
	}
	return models.PostResponse{
		ID:        post.ID,
		Author:    userData,
		Caption:   post.Caption,
		Tag:       tag,
		ImageUrls: image,
		Likes:     post.Likes,
		Comments:  post.Comments,
		CreatedAt: post.CreatedAt,
	}, nil
}

func (p *postUseCase) UpdatePost(userID int, data models.UpdatePostReq) (models.PostResponse, error) {
	userExist := p.authClient.CheckUserAvalilabilityWithUserID(userID)
	if !userExist {
		return models.PostResponse{}, errors.New("user doesn't exist")
	}
	ok := p.postRepository.CheckPostAvalilabilityWithID(int(data.PostID))
	if !ok {
		return models.PostResponse{}, errors.New("post doesn't exist")
	}
	if data.Caption != "" {
		p.postRepository.UpdateCaption(userID, int(data.PostID), data.Caption)
	}
	mediatype := p.postRepository.CheckMediaAvalilabilityWithID(int(data.TypeID))
	if !mediatype {
		return models.PostResponse{}, errors.New("type doesn't exist")
	}
	if data.TypeID != 0 {
		p.postRepository.UpdateTypeID(userID, int(data.PostID), int(data.TypeID))
	}
	p.postRepository.UpdateTags(userID, int(data.PostID), data.Tags)
	var image []models.Url
	post, tag, image, err := p.postRepository.PostDetails(int(data.TypeID), userID)
	if err != nil {
		return models.PostResponse{}, err
	}
	userData, err := p.authClient.UserData(userID)
	if err != nil {
		return models.PostResponse{}, err
	}
	return models.PostResponse{
		ID:        post.ID,
		Author:    userData,
		Caption:   post.Caption,
		Tag:       tag,
		ImageUrls: image,
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
