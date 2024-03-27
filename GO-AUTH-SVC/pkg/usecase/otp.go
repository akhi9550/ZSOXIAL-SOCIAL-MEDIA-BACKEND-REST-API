package usecase

import (
	"errors"

	"github.com/akhi9550/pkg/config"
	"github.com/akhi9550/pkg/helper"
	interfaces "github.com/akhi9550/pkg/repository/interface"
	services "github.com/akhi9550/pkg/usecase/interface"
	"github.com/akhi9550/pkg/utils/models"
	"github.com/jinzhu/copier"
)

type otpUseCase struct {
	otpRepository interfaces.OtpRepository
}

func NewOtpUseCase(repo interfaces.OtpRepository) services.OtpUseCase {
	return &otpUseCase{
		otpRepository: repo,
	}
}

func (op *otpUseCase) SendOtp(phone string) error {
	cfg, err := config.LoadConfig()

	if err != nil {
		return err
	}
	user, err := op.otpRepository.FindUserByPhoneNumber(phone)
	if err != nil {
		return errors.New("error with server")
	}
	if user == nil {
		return errors.New("user with this phone is not exists")
	}
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)
	_, err = helper.TwilioSendOTP(phone, cfg.SERVICESSID)

	if err != nil {
		return errors.New("error occured while generating otp")
	}
	return nil
}
func (op *otpUseCase) VerifyOTP(code models.VerifyData) (models.ReponseWithToken, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return models.ReponseWithToken{}, err
	}
	helper.TwilioSetup(cfg.ACCOUNTSID, cfg.AUTHTOKEN)
	err = helper.TwilioVerifyOTP(cfg.SERVICESSID, code.Code, code.User.PhoneNumber)
	if err != nil {
		return models.ReponseWithToken{}, errors.New("error while verifying")
	}
	userDetails, err := op.otpRepository.UserDetailsUsingPhone(code.User.PhoneNumber)
	if err != nil {
		return models.ReponseWithToken{}, err
	}
	accessToken, err := helper.GenerateAccessTokenUser(userDetails)
	if err != nil {
		return models.ReponseWithToken{}, errors.New("couldn't create token due to some internal error")
	}
	refreshToken, err := helper.GenerateRefreshTokenUser(userDetails)
	if err != nil {
		return models.ReponseWithToken{}, errors.New("couldn't create token due to some internal error")
	}
	var user models.UserResponse
	err = copier.Copy(&user, &userDetails)
	if err != nil {
		return models.ReponseWithToken{}, err
	}
	return models.ReponseWithToken{
		Users:        user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
