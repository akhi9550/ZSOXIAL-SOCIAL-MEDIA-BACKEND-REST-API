package usecase

import (
	"errors"

	postclientinterfaces "github.com/akhi9550/auth-svc/pkg/client/interface"
	"github.com/akhi9550/auth-svc/pkg/helper"
	interfaces "github.com/akhi9550/auth-svc/pkg/repository/interface"
	services "github.com/akhi9550/auth-svc/pkg/usecase/interface"
	"github.com/akhi9550/auth-svc/pkg/utils/models"
	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

type adminUseCase struct {
	adminRepository interfaces.AdminRepository
	postClient      postclientinterfaces.NewPostClient
}

func NewAdminUseCase(repository interfaces.AdminRepository, postclient postclientinterfaces.NewPostClient) services.AdminUseCase {
	return &adminUseCase{
		adminRepository: repository,
		postClient:      postclient,
	}
}

func (ad *adminUseCase) AdminLogin(admin models.AdminLoginRequest) (*models.AdminReponseWithToken, error) {
	email, err := ad.adminRepository.CheckAdminExistsByEmail(admin.Email)
	if err != nil {
		return &models.AdminReponseWithToken{}, errors.New("error with server")
	}
	if email == nil {
		return &models.AdminReponseWithToken{}, errors.New("email doesn't exist")
	}
	admindeatils, err := ad.adminRepository.FindAdminByEmail(admin)
	if err != nil {
		return &models.AdminReponseWithToken{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(admindeatils.Password), []byte(admin.Password))
	if err != nil {
		return &models.AdminReponseWithToken{}, errors.New("password not matching")
	}
	var adminDetails models.AdminResponse
	err = copier.Copy(&adminDetails, &admindeatils)
	if err != nil {
		return &models.AdminReponseWithToken{}, err
	}
	accessToken, err := helper.GenerateAccessTokenAdmin(adminDetails)
	if err != nil {
		return &models.AdminReponseWithToken{}, errors.New("couldn't create accesstoken due to internal error")
	}
	refreshToken, err := helper.GenerateRefreshTokenAdmin(adminDetails)
	if err != nil {
		return &models.AdminReponseWithToken{}, errors.New("counldn't create refreshtoken due to internal error")
	}

	return &models.AdminReponseWithToken{
		Users:        adminDetails,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (ad *adminUseCase) ShowAllUsers(page, count int) ([]models.UserDetailsAtAdmin, error) {
	users, err := ad.adminRepository.ShowAllUsers(page, count)
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}
	return users, nil
}

func (ad *adminUseCase) AdminBlockUser(userID uint) error {
	user, err := ad.adminRepository.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user.Blocked {
		return errors.New("already blocked")
	} else {
		user.Blocked = true
	}
	err = ad.adminRepository.AdminBlockUserByID(user)
	if err != nil {
		return err
	}
	return nil
}

func (ad *adminUseCase) AdminUnBlockUser(userID uint) error {
	user, err := ad.adminRepository.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user.Blocked {
		user.Blocked = false
	} else {
		return errors.New("already unblocked")
	}
	err = ad.adminRepository.AdminBlockUserByID(user)
	if err != nil {
		return err
	}
	return nil
}

func (ad *adminUseCase) ShowUserReports(page, count int) ([]models.UserReports, error) {
	reports, err := ad.adminRepository.ShowUserReports(page, count)
	if err != nil {
		return []models.UserReports{}, err
	}
	return reports, nil
}

func (ad *adminUseCase) ShowPostReports(page, count int) ([]models.PostReports, error) {
	reports, err := ad.postClient.ShowPostReports(page, count)
	if err != nil {
		return []models.PostReports{}, err
	}
	return reports, nil
}

func (ad *adminUseCase) GetAllPosts(page, count int) ([]models.PostResponse, error) {
	reports, err := ad.postClient.GetAllPosts(page, count)
	if err != nil {
		return []models.PostResponse{}, err
	}
	return reports, nil
}

func (ad *adminUseCase) RemovePost(postID int) error {
	postExist := ad.postClient.CheckPostIDByID(postID)
	if !postExist {
		return errors.New("postID doesn't exist")
	}
	err := ad.postClient.RemovePost(postID)
	if err != nil {
		return err
	}
	return nil
}
