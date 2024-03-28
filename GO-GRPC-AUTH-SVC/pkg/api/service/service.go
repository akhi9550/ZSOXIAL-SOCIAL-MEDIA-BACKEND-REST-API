package service

import (
	"context"

	"github.com/akhi9550/auth-svc/pkg/pb"
	interfaces "github.com/akhi9550/auth-svc/pkg/usecase/interface"
	"github.com/akhi9550/auth-svc/pkg/utils/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthSever struct {
	userUseCase  interfaces.UserUseCase
	adminUsecase interfaces.AdminUseCase
	otpUsecase   interfaces.OtpUseCase
	pb.UnimplementedAuthServiceServer
}

func NewAuthServer(useCaseUser interfaces.UserUseCase, useCaseAdmin interfaces.AdminUseCase, useCaseOtp interfaces.OtpUseCase) pb.AuthServiceServer {
	return &AuthSever{
		userUseCase:  useCaseUser,
		adminUsecase: useCaseAdmin,
		otpUsecase:   useCaseOtp,
	}
}

func (au *AuthSever) UserSignUp(ctx context.Context, user *pb.UserSignUpRequest) (*pb.UserSignUpResponse, error) {
	signup := models.UserSignUpRequest{
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Username:  user.Username,
		Dob:       user.Dob,
		Gender:    user.Gender,
		Phone:     user.Phone,
		Email:     user.Email,
		Password:  user.Password,
		Bio:       user.Bio,
	}
	file := models.UserProfilePhoto{
		Imageurl: user.ProfilePhoto.ProfilePhoto,
	}
	File := file.Imageurl
	data, err := au.userUseCase.UserSignUp(signup, File)
	if err != nil {
		return &pb.UserSignUpResponse{}, err
	}
	userDetails := &pb.UserInfo{Id: int64(data.Users.Id), Username: data.Users.Username, ProfilePhoto: data.Users.Imageurl, Isadmin: data.Users.Isadmin}
	UserResponse := &pb.UserResponse{Info: userDetails, Accesstoken: data.AccessToken, Refreshtoken: data.RefreshToken}
	return &pb.UserSignUpResponse{
		Reposnse: UserResponse,
	}, nil
}

func (au *AuthSever) UserLogin(ctx context.Context, user *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	login := models.UserLoginRequest{
		Email:    user.Email,
		Password: user.Password,
	}
	data, err := au.userUseCase.UserLogin(login)
	if err != nil {
		return &pb.UserLoginResponse{}, err
	}
	userDetails := &pb.UserInfo{Id: int64(data.Users.Id), Username: data.Users.Username, ProfilePhoto: data.Users.Imageurl, Isadmin: data.Users.Isadmin}
	UserResponse := &pb.UserResponse{Info: userDetails, Accesstoken: data.AccessToken, Refreshtoken: data.RefreshToken}
	return &pb.UserLoginResponse{
		Reposnse: UserResponse,
	}, nil
}

func (au *AuthSever) SendOtp(ctx context.Context, req *pb.SendOtpRequest) (*pb.SendOtpResponse, error) {
	phone := req.Phone
	err := au.otpUsecase.SendOtp(phone)
	if err != nil {
		return &pb.SendOtpResponse{}, err
	}
	return &pb.SendOtpResponse{}, nil
}

func (au *AuthSever) VerifyOtp(ctx context.Context, req *pb.VerifyOtpRequest) (*pb.VerifyOtpResponse, error) {
	otpData := &models.OTPData{
		PhoneNumber: req.Otpdata.User.Phone,
	}

	verifyData := models.VerifyData{
		User: otpData,
		Code: req.Otpdata.Code,
	}
	data, err := au.otpUsecase.VerifyOTP(verifyData)
	if err != nil {
		return &pb.VerifyOtpResponse{}, err
	}
	userDetails := &pb.UserInfo{Id: int64(data.Users.Id), Username: data.Users.Username, ProfilePhoto: data.Users.Imageurl, Isadmin: data.Users.Isadmin}
	UserResponse := &pb.UserResponse{Info: userDetails, Accesstoken: data.AccessToken, Refreshtoken: data.RefreshToken}
	return &pb.VerifyOtpResponse{
		Reposnse: UserResponse,
	}, nil
}

func (au *AuthSever) ForgotPassword(ctx context.Context, req *pb.ForgotPasswordRequest) (*pb.ForgotPasswordResponse, error) {
	phone := req.Phone
	err := au.userUseCase.ForgotPassword(phone)
	if err != nil {
		return &pb.ForgotPasswordResponse{}, nil
	}
	return &pb.ForgotPasswordResponse{}, nil
}

func (au *AuthSever) ForgotPasswordVerifyAndChange(ctx context.Context, req *pb.ForgotPasswordVerifyAndChangeRequest) (*pb.ForgotPasswordVerifyAndChangeResponse, error) {
	verify := models.ForgotVerify{
		Phone:       req.Verify.Phone,
		Otp:         req.Verify.Otp,
		NewPassword: req.Verify.Newpassword,
	}
	err := au.userUseCase.ForgotPasswordVerifyAndChange(verify)
	if err != nil {
		return &pb.ForgotPasswordVerifyAndChangeResponse{}, nil
	}
	return &pb.ForgotPasswordVerifyAndChangeResponse{}, nil
}

func (au *AuthSever) UserDetails(ctx context.Context, req *pb.UserDetailsRequest) (*pb.UserDetailsResponse, error) {
	userID := req.Id
	data, err := au.userUseCase.UserDetails(int(userID))
	if err != nil {
		return &pb.UserDetailsResponse{}, err
	}
	userData := &pb.UserData{Firstname: data.Firstname, Lastname: data.Lastname, Username: data.Username, Dob: data.Dob, Gender: data.Gender, Phone: data.Phone, Email: data.Email, Bio: data.Bio, ProfilePhoto: data.Imageurl}
	return &pb.UserDetailsResponse{
		Responsedata: userData,
	}, nil
}

func (au *AuthSever) UpdateUserDetails(ctx context.Context, req *pb.UpdateUserDetailsRequest) (*pb.UpdateUserDetailsResponse, error) {
	userData := models.UsersProfileDetail{
		Firstname: req.UserDetails.Firstname,
		Lastname:  req.UserDetails.Lastname,
		Username:  req.UserDetails.Username,
		Dob:       req.UserDetails.Dob,
		Gender:    req.UserDetails.Gender,
		Phone:     req.UserDetails.Phone,
		Email:     req.UserDetails.Email,
		Bio:       req.UserDetails.Bio,
	}
	file := models.UserProfilePhoto{
		Imageurl: req.UserDetails.ProfilePhoto.ProfilePhoto,
	}
	File := file.Imageurl
	userID := req.Id
	data, err := au.userUseCase.UpdateUserDetails(userData,File, int(userID))
	if err != nil {
		return &pb.UpdateUserDetailsResponse{}, err
	}
	UserData := &pb.UserData{Firstname: data.Firstname, Lastname: data.Lastname, Username: data.Username, Dob: data.Dob, Gender: data.Gender, Phone: data.Phone, Email: data.Email, Bio: data.Bio, ProfilePhoto: data.Imageurl}
	return &pb.UpdateUserDetailsResponse{
		UserDetails: UserData,
	}, nil
}

func (au *AuthSever) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	userID := req.Id
	PasswordDetails := models.ChangePassword{
		Oldpassword: req.Password.Oldpassword,
		Password:    req.Password.Password,
		Repassword:  req.Password.Repassword,
	}
	err := au.userUseCase.ChangePassword(int(userID), PasswordDetails)
	if err != nil {
		return &pb.ChangePasswordResponse{}, err
	}
	return &pb.ChangePasswordResponse{}, nil
}

func (au *AuthSever) AdminLogin(ctx context.Context, req *pb.AdminLoginRequest) (*pb.AdminLoginResponse, error) {
	login := models.AdminLoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
	data, err := au.adminUsecase.AdminLogin(login)
	if err != nil {
		return &pb.AdminLoginResponse{}, err
	}
	adminDetails := &pb.AdminInfo{Id: int64(data.Users.Id), Email: data.Users.Email, ProfilePhoto: []byte(data.Users.Imageurl), Isadmin: data.Users.Isadmin}
	AdminResponse := &pb.AdminResponse{Info: adminDetails, Accesstoken: data.AccessToken, Refreshtoken: data.RefreshToken}
	return &pb.AdminLoginResponse{
		Reposnse: AdminResponse,
	}, nil
}

func (au *AuthSever) ShowAllUsers(ctx context.Context, req *pb.ShowAllUsersRequest) (*pb.ShowAllUsersResponse, error) {
	page := req.Page
	count := req.Count
	data, err := au.adminUsecase.ShowAllUsers(int(page), int(count))
	if err != nil {
		return &pb.ShowAllUsersResponse{
			UsersData: []*pb.Users{},
		}, err
	}
	var result pb.ShowAllUsersResponse
	for _, v := range data {
		createdAtTimestamp := timestamppb.New(v.CreatedAt)
		result.UsersData = append(result.UsersData, &pb.Users{
			Id:           int64(v.Id),
			Firstname:    v.Firstname,
			Lastname:     v.Lastname,
			Username:     v.Username,
			Dob:          v.Dob,
			Gender:       v.Gender,
			Phone:        v.Phone,
			Email:        v.Email,
			ProfilePhoto: v.Imageurl,
			CreatedAt:    createdAtTimestamp,
			Blocked:      v.Blocked,
		})
	}
	return &result, nil
}

func (au *AuthSever) AdminBlockUser(ctx context.Context, req *pb.AdminBlockUserRequest) (*pb.AdminBlockUserResponse, error) {
	userID := req.Id
	err := au.adminUsecase.AdminBlockUser(uint(userID))
	if err != nil {
		return &pb.AdminBlockUserResponse{}, err
	}
	return &pb.AdminBlockUserResponse{}, nil
}

func (au *AuthSever) AdminUnblockUser(ctx context.Context, req *pb.AdminUnblockUserRequest) (*pb.AdminUnblockUserResponse, error) {
	userID := req.Id
	err := au.adminUsecase.AdminUnBlockUser(uint(userID))
	if err != nil {
		return &pb.AdminUnblockUserResponse{}, err
	}
	return &pb.AdminUnblockUserResponse{}, nil
}
