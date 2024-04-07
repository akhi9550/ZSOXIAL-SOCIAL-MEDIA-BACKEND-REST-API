package client

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/config"
	pb "github.com/akhi9550/api-gateway/pkg/pb/auth"
	"github.com/akhi9550/api-gateway/pkg/utils/models"

	"google.golang.org/grpc"
)

type AuthClient struct {
	Client pb.AuthServiceClient
}

func NewAuthClient(cfg config.Config) interfaces.AuthClient {
	grpcConnection, err := grpc.Dial(cfg.AuthSvcUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect", err)
	}

	grpcClient := pb.NewAuthServiceClient(grpcConnection)

	return &AuthClient{
		Client: grpcClient,
	}
}

func (au *AuthClient) UserSignUp(user models.UserSignUpRequest) (*models.ReponseWithToken, error) {

	data, err := au.Client.UserSignUp(context.Background(), &pb.UserSignUpRequest{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Username:  user.Username,
		Phone:     user.Phone,
		Email:     user.Email,
		Password:  user.Password,
	})
	if err != nil {
		return nil, err
	}

	userData := models.UserResponse{
		Id:       uint(data.Reposnse.Info.Id),
		Username: data.Reposnse.Info.Username,
		Email:    data.Reposnse.Info.Email,
		Isadmin:  data.Reposnse.Info.Isadmin,
	}
	return &models.ReponseWithToken{
		Users:        userData,
		AccessToken:  data.Reposnse.Accesstoken,
		RefreshToken: data.Reposnse.Refreshtoken,
	}, nil
}

func (au *AuthClient) UserLogin(user models.UserLoginRequest) (*models.ReponseWithToken, error) {
	data, err := au.Client.UserLogin(context.Background(), &pb.UserLoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})
	if err != nil {
		return nil, err
	}
	userData := models.UserResponse{
		Id:       uint(data.Reposnse.Info.Id),
		Username: data.Reposnse.Info.Username,
		Email:    data.Reposnse.Info.Email,
		Imageurl: data.Reposnse.Info.ProfilePhoto,
		Isadmin:  data.Reposnse.Info.Isadmin,
	}
	return &models.ReponseWithToken{
		Users:        userData,
		AccessToken:  data.Reposnse.Accesstoken,
		RefreshToken: data.Reposnse.Refreshtoken,
	}, nil
}

func (au *AuthClient) ForgotPassword(phone string) error {
	_, err := au.Client.ForgotPassword(context.Background(), &pb.ForgotPasswordRequest{
		Phone: phone,
	})
	if err != nil {
		return err
	}
	return nil
}

func (au *AuthClient) ForgotPasswordVerifyAndChange(req models.ForgotVerify) error {
	forgotverify := &pb.ForgotVerify{Phone: req.Phone, Otp: req.Otp, Newpassword: req.NewPassword}
	_, err := au.Client.ForgotPasswordVerifyAndChange(context.Background(), &pb.ForgotPasswordVerifyAndChangeRequest{
		Verify: forgotverify,
	})
	if err != nil {
		return err
	}
	return nil
}

func (au *AuthClient) UserDetails(userID int) (models.UsersProfileDetails, error) {
	data, err := au.Client.UserDetails(context.Background(), &pb.UserDetailsRequest{
		Id: int64(userID),
	})
	if err != nil {
		return models.UsersProfileDetails{}, err
	}
	return models.UsersProfileDetails{
		Firstname: data.Responsedata.Firstname,
		Lastname:  data.Responsedata.Lastname,
		Username:  data.Responsedata.Username,
		Dob:       data.Responsedata.Dob,
		Gender:    data.Responsedata.Gender,
		Phone:     data.Responsedata.Phone,
		Email:     data.Responsedata.Email,
		Bio:       data.Responsedata.Bio,
		Imageurl:  data.Responsedata.ProfilePhoto,
	}, nil
}

func (au *AuthClient) UpdateUserDetails(userDetails models.UsersProfileDetail, file *multipart.FileHeader, userID int) (models.UsersProfileDetails, error) {
	f, err := file.Open()
	if err != nil {
		return models.UsersProfileDetails{}, err
	}
	defer f.Close()
	fileData, err := io.ReadAll(f)
	if err != nil {
		return models.UsersProfileDetails{}, err
	}
	files := &pb.UserProfile{ProfilePhoto: fileData}
	userData := &pb.UserDataUpdate{
		Firstname:    userDetails.Firstname,
		Lastname:     userDetails.Lastname,
		Username:     userDetails.Username,
		Dob:          userDetails.Dob,
		Gender:       userDetails.Gender,
		Phone:        userDetails.Phone,
		Email:        userDetails.Email,
		Bio:          userDetails.Bio,
		ProfilePhoto: files,
	}
	data, err := au.Client.UpdateUserDetails(context.Background(), &pb.UpdateUserDetailsRequest{
		UserDetails: userData,
		Id:          int64(userID),
	})
	if err != nil {
		return models.UsersProfileDetails{}, err
	}
	return models.UsersProfileDetails{
		Firstname: data.UserDetails.Firstname,
		Lastname:  data.UserDetails.Lastname,
		Username:  data.UserDetails.Username,
		Dob:       data.UserDetails.Dob,
		Gender:    data.UserDetails.Gender,
		Phone:     data.UserDetails.Phone,
		Email:     data.UserDetails.Email,
		Bio:       data.UserDetails.Bio,
		Imageurl:  data.UserDetails.ProfilePhoto,
	}, nil
}

func (au *AuthClient) ChangePassword(userID int, change models.ChangePassword) error {
	changepassword := &pb.ChangePassword{
		Oldpassword: change.Oldpassword,
		Password:    change.Password,
		Repassword:  change.Repassword,
	}
	_, err := au.Client.ChangePassword(context.Background(), &pb.ChangePasswordRequest{
		Id:       int64(userID),
		Password: changepassword,
	})
	if err != nil {
		return err
	}
	return nil
}

func (au *AuthClient) SendOtp(phone string) error {
	_, err := au.Client.SendOtp(context.Background(), &pb.SendOtpRequest{
		Phone: phone,
	})
	if err != nil {
		return err
	}
	return nil
}

func (au *AuthClient) VerifyOTP(code models.VerifyData) (models.ReponseWithToken, error) {
	otp := &pb.OTPData{
		Phone: code.User.PhoneNumber,
	}
	otpverify := &pb.VerifyOtp{
		User: otp,
		Code: code.Code,
	}
	data, err := au.Client.VerifyOtp(context.Background(), &pb.VerifyOtpRequest{
		Otpdata: otpverify,
	})
	if err != nil {
		return models.ReponseWithToken{}, err
	}
	userData := models.UserResponse{
		Id:       uint(data.Reposnse.Info.Id),
		Username: data.Reposnse.Info.Username,
		Email:    data.Reposnse.Info.Email,
		Imageurl: string(data.Reposnse.Info.ProfilePhoto),
		Isadmin:  data.Reposnse.Info.Isadmin,
	}
	return models.ReponseWithToken{
		Users:        userData,
		AccessToken:  data.Reposnse.Accesstoken,
		RefreshToken: data.Reposnse.Refreshtoken,
	}, nil
}

func (au *AuthClient) AdminLogin(admin models.AdminLoginRequest) (*models.AdminReponseWithToken, error) {
	data, err := au.Client.AdminLogin(context.Background(), &pb.AdminLoginRequest{
		Email:    admin.Email,
		Password: admin.Password,
	})
	if err != nil {
		return &models.AdminReponseWithToken{}, err
	}
	adminData := models.AdminResponse{
		Id:       uint(data.Reposnse.Info.Id),
		Email:    data.Reposnse.Info.Email,
		Imageurl: string(data.Reposnse.Info.ProfilePhoto),
		Isadmin:  data.Reposnse.Info.Isadmin,
	}
	return &models.AdminReponseWithToken{
		Users:        adminData,
		AccessToken:  data.Reposnse.Accesstoken,
		RefreshToken: data.Reposnse.Refreshtoken,
	}, nil
}

func (au *AuthClient) ShowAllUsers(page, count int) ([]models.UserDetailsAtAdmin, error) {
	data, err := au.Client.ShowAllUsers(context.Background(), &pb.ShowAllUsersRequest{
		Page:  int64(page),
		Count: int64(count),
	})
	if err != nil {
		return []models.UserDetailsAtAdmin{}, err
	}
	var result []models.UserDetailsAtAdmin

	for _, userData := range data.UsersData {
		userDetails := models.UserDetailsAtAdmin{
			Id:        uint(userData.Id),
			Firstname: userData.Firstname,
			Lastname:  userData.Lastname,
			Username:  userData.Username,
			Dob:       userData.Dob,
			Gender:    userData.Gender,
			Phone:     userData.Phone,
			Email:     userData.Email,
			Imageurl:  string(userData.ProfilePhoto),
			CreatedAt: time.Unix(userData.CreatedAt.Seconds, 0),
			Blocked:   userData.Blocked,
		}
		result = append(result, userDetails)
	}

	return result, nil
}

func (au *AuthClient) AdminBlockUser(userID int) error {
	_, err := au.Client.AdminBlockUser(context.Background(), &pb.AdminBlockUserRequest{
		Id: int64(userID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (au *AuthClient) AdminUnblockUser(userID int) error {
	_, err := au.Client.AdminUnblockUser(context.Background(), &pb.AdminUnblockUserRequest{
		Id: int64(userID),
	})
	if err != nil {
		return err
	}
	return nil
}

func (au *AuthClient) ReportUser(userID int, req models.ReportRequest) error {
	_, err := au.Client.ReportUser(context.Background(), &pb.ReportUserRequest{
		RepostedUserid: int64(userID),
		Userid:         int64(req.UserID),
		Report:         req.Report,
	})
	if err != nil {
		return err
	}
	return nil
}
