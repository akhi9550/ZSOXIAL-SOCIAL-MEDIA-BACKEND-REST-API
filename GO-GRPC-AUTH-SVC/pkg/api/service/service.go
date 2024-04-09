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
		Phone:     user.Phone,
		Email:     user.Email,
		Password:  user.Password,
	}
	data, err := au.userUseCase.UserSignUp(signup)
	if err != nil {
		return &pb.UserSignUpResponse{}, err
	}
	userDetails := &pb.UserInfo{Id: int64(data.Users.Id), Username: data.Users.Username, Email: data.Users.Email, Isadmin: data.Users.Isadmin}
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
	userDetails := &pb.UserInfo{Id: int64(data.Users.Id), Username: data.Users.Username, ProfilePhoto: data.Users.Imageurl, Email: data.Users.Email, Isadmin: data.Users.Isadmin}
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
	userDetails := &pb.UserInfo{Id: int64(data.Users.Id), Username: data.Users.Username, Email: data.Users.Email, ProfilePhoto: data.Users.Imageurl, Isadmin: data.Users.Isadmin}
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
	data, err := au.userUseCase.UpdateUserDetails(userData, File, int(userID))
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

func (au *AuthSever) CheckUserAvalilabilityWithUserID(ctx context.Context, req *pb.CheckUserAvalilabilityWithUserIDRequest) (*pb.CheckUserAvalilabilityWithUserIDResponse, error) {
	userId := req.Id
	ok, err := au.userUseCase.CheckUserAvalilabilityWithUserID(int(userId))
	if !ok {
		return &pb.CheckUserAvalilabilityWithUserIDResponse{}, err
	}
	if err != nil {
		return &pb.CheckUserAvalilabilityWithUserIDResponse{}, err
	}
	return &pb.CheckUserAvalilabilityWithUserIDResponse{
		Valid: ok,
	}, nil
}

func (au *AuthSever) UserData(ctx context.Context, req *pb.UserDataRequest) (*pb.UserDataResponse, error) {
	userId := req.Id
	data, err := au.userUseCase.UserData(int(userId))
	if err != nil {
		return &pb.UserDataResponse{}, err
	}
	return &pb.UserDataResponse{
		Id:           int64(data.UserId),
		Username:     data.Username,
		ProfilePhoto: data.Profile,
	}, err
}

func (au *AuthSever) CheckUserAvalilabilityWithTagUserID(ctx context.Context, req *pb.CheckUserAvalilabilityWithTagUserIDRequest) (*pb.CheckUserAvalilabilityWithTagUserIDResponse, error) {
	var users []models.Tag
	for _, user := range req.Tag.User {
		tag := models.Tag{User: user}
		users = append(users, tag)
	}

	ok, err := au.userUseCase.CheckUserAvalilabilityWithTagUserID(users)
	if !ok {
		return &pb.CheckUserAvalilabilityWithTagUserIDResponse{}, err
	}
	if err != nil {
		return &pb.CheckUserAvalilabilityWithTagUserIDResponse{}, err
	}
	return &pb.CheckUserAvalilabilityWithTagUserIDResponse{
		Valid: ok,
	}, nil
}

func (au *AuthSever) GetUserNameWithTagUserID(ctx context.Context, req *pb.GetUserNameWithTagUserIDRequest) (*pb.GetUserNameWithTagUserIDResponse, error) {
	var users []models.Tag
	for _, user := range req.Tag.User {
		tag := models.Tag{User: user}
		users = append(users, tag)
	}
	data, err := au.userUseCase.GetUserNameWithTagUserID(users)
	if err != nil {
		return nil, err
	}

	var tagUsers []*pb.TagUsernames
	for _, user := range data {
		tagUsers = append(tagUsers, &pb.TagUsernames{
			Username: user.Username,
		})
	}
	return &pb.GetUserNameWithTagUserIDResponse{
		Name: tagUsers,
	}, nil
}

func (au *AuthSever) ReportUser(ctx context.Context, req *pb.ReportUserRequest) (*pb.ReportUserResponse, error) {
	ReportUser := req.RepostedUserid
	reportReq := models.ReportRequest{
		UserID: uint(req.Userid),
		Report: req.Report,
	}
	err := au.userUseCase.ReportUser(int(ReportUser), reportReq)
	if err != nil {
		return &pb.ReportUserResponse{}, err
	}
	return &pb.ReportUserResponse{}, nil
}

func (au *AuthSever) FollowREQ(ctx context.Context, req *pb.FollowREQRequest) (*pb.FollowREQResponse, error) {
	userID, FollowingUserID := req.Userid, req.FollowingUser
	err := au.userUseCase.FollowREQ(int(userID), int(FollowingUserID))
	if err != nil {
		return &pb.FollowREQResponse{}, err
	}
	return &pb.FollowREQResponse{}, nil
}

func (au *AuthSever) ShowFollowREQ(ctx context.Context, req *pb.ShowREQRequest) (*pb.ShowREQResponse, error) {
	userID := req.Userid
	data, err := au.userUseCase.ShowFollowREQ(int(userID))
	if err != nil {
		return &pb.ShowREQResponse{}, err
	}
	var response []*pb.REQResponse
	for _, req := range data {
		requests := &pb.REQResponse{
			FollowingUserID: int64(req.FollowingUserID),
			FollowingUser:   req.FollowingUser,
			ProfilePhoto:    req.Profile,
			CreatedAt:       timestamppb.New(req.CreatedAt),
		}
		response = append(response, requests)
	}
	return &pb.ShowREQResponse{
		Response: response,
	}, nil
}

func (au *AuthSever) AcceptFollowREQ(ctx context.Context, req *pb.AcceptFollowREQRequest) (*pb.AcceptFollowREQResponse, error) {
	userID, FollowingUserID := req.Userid, req.FollowingUser
	err := au.userUseCase.AcceptFollowREQ(int(userID), int(FollowingUserID))
	if err != nil {
		return &pb.AcceptFollowREQResponse{}, err
	}
	return &pb.AcceptFollowREQResponse{}, nil
}

func (au *AuthSever) Following(ctx context.Context, req *pb.FollowingRequest) (*pb.FollowingResponse, error) {
	userID := req.Userid
	data, err := au.userUseCase.Following(int(userID))
	if err != nil {
		return &pb.FollowingResponse{}, err
	}
	var response []*pb.FollowResponse
	for _, req := range data {
		requests := &pb.FollowResponse{
			Username:    req.FollowingUser,
			UserProfile: req.Profile,
		}
		response = append(response, requests)
	}
	return &pb.FollowingResponse{
		Users: response,
	}, nil
}

func (au *AuthSever) Follower(ctx context.Context, req *pb.FollowerRequest) (*pb.FollowerResponse, error) {
	userID := req.Userid
	data, err := au.userUseCase.Follower(int(userID))
	if err != nil {
		return &pb.FollowerResponse{}, err
	}
	var response []*pb.FollowResponse
	for _, req := range data {
		requests := &pb.FollowResponse{
			Username:    req.FollowingUser,
			UserProfile: req.Profile,
		}
		response = append(response, requests)
	}
	return &pb.FollowerResponse{
		Users: response,
	}, nil
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

func (au *AuthSever)ShowUserReports(ctx context.Context,req *pb.showuser)(*pb.s)
ShowUserReports(page, count int) ([]models.UserReports, error)
	ShowPostReports(page, count int) ([]models.PostReports, error)
	GetAllPosts(page, count int) ([]models.PostResponse, error)
	RemovePost(postID int) error