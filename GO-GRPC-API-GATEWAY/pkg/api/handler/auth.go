package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/helper"
	"github.com/akhi9550/api-gateway/pkg/logging"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/akhi9550/api-gateway/pkg/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	GRPC_Client interfaces.AuthClient
	AuthCachig  *helper.RedisAuthCaching
}

func NewAuthHandler(authClient interfaces.AuthClient, authCaching *helper.RedisAuthCaching) *AuthHandler {
	return &AuthHandler{
		GRPC_Client: authClient,
		AuthCachig:  authCaching,
	}
}

// @Summary		User Signup
// @Description	user can signup by giving their details
// @Tags			User
// @Accept			json
// @Produce		    json
// @Param			UserSignupDetail   body  models.UserSignUpRequest  true  "User Signup"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/signup    [POST]
func (au *AuthHandler) UserSignup(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "UserSignupHandler")
	logEntry.Info("Processing user Signup request")
	var UserSignupDetail models.UserSignUpRequest
	if err := c.ShouldBindJSON(&UserSignupDetail); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
	}

	pattern := `^\d{10}$`
	regex := regexp.MustCompile(pattern)
	value := regex.MatchString(UserSignupDetail.Phone)
	if !value {
		fmt.Printf("%s phone number is not valid", UserSignupDetail.Phone)
		return
	}

	err := validator.New().Struct(UserSignupDetail)
	if err != nil {
		logEntry.WithError(err).Error("Error Constraints not statisfied")
		errs := response.ClientResponse(http.StatusBadRequest, "Constraints not statisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	user, err := au.GRPC_Client.UserSignUp(UserSignupDetail)
	if err != nil {
		logEntry.WithError(err).Error("Error during User Signup rpc call")
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	logEntry.Info("user signup succesfull for user")
	success := response.ClientResponse(http.StatusCreated, "User successfully signed up", user, nil)
	c.JSON(http.StatusCreated, success)
}

// @Summary		User Login
// @Description	user can log in by giving their details
// @Tags			User
// @Accept			json
// @Produce		    json
// @Param			UserLoginDetail  body  models.UserLoginRequest  true	"User Login"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/login     [POST]
func (au *AuthHandler) Userlogin(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "LogginHandler")
	logEntry.Info("Processing Loggin request")
	var UserLoginDetail models.UserLoginRequest
	if err := c.ShouldBindJSON(&UserLoginDetail); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
	}
	err := validator.New().Struct(UserLoginDetail)
	if err != nil {
		logEntry.WithError(err).Error("Error Constraints not statisfied")
		errs := response.ClientResponse(http.StatusBadRequest, "Constraints not statisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	user, err := au.GRPC_Client.UserLogin(UserLoginDetail)
	if err != nil {
		logEntry.WithError(err).Error("Error During UserLogin RPC call")
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	logEntry.Info("Login successful for User")
	success := response.ClientResponse(http.StatusOK, "User successfully logged in with password", user, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary 	OTP login
// @Description Send OTP to Authenticate user
// @Tags 		User OTP Login
// @Accept 		json
// @Produce 	json
// @Param 		phone 	body models.OTPData 	true 	"phone number details"
// @Success 	200 	{object} response.Response{}
// @Failure 	500 	{object} response.Response{}
// @Router 		/user/send-otp   [POST]
func (au *AuthHandler) SendOtp(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "SendOtp")
	logEntry.Info("Processing Sendotp Login request")
	var phone models.OTPData
	if err := c.ShouldBindJSON(&phone); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err := au.GRPC_Client.SendOtp(phone.PhoneNumber)
	if err != nil {
		logEntry.WithError(err).Error("Error During Sendotp RPC call")
		errs := response.ClientResponse(http.StatusBadRequest, "user with this phone is not exists", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	logEntry.Info("OTP sent successfully")
	success := response.ClientResponse(http.StatusOK, "OTP sent successfully", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary 	Verify OTP
// @Description Verify OTP by passing the OTP in order to authenticate user
// @Tags 		User OTP Login
// @Accept 		json
// @Produce 	json
// @Param 		code 	body models.VerifyData 	true 	"Verify OTP Details"
// @Success 	200 	{object} response.Response{}
// @Failure 	500 	{object} response.Response{}
// @Router 		/user/verify-otp      [POST]
func (au *AuthHandler) VerifyOtp(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "VerifyOtp")
	logEntry.Info("Processing VerifyOtp Login request")
	var code models.VerifyData
	if err := c.ShouldBindJSON(&code); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	user, err := au.GRPC_Client.VerifyOTP(code)
	if err != nil {
		logEntry.WithError(err).Error("Error During VerifyOTP RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "Could not verify OTP", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully verified OTP")
	sucess := response.ClientResponse(http.StatusOK, "Successfully verified OTP", user, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary			Forgot password Send OTP
// @Description		User can change their password if user forgot the password and login
// @Tags			User OTP Login
// @Accept			json
// @Produce		    json
// @Param			model  body  models.ForgotPasswordSend  true	"forgot-send"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/forgot-password   [POST]
func (au *AuthHandler) ForgotPassword(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "ForgotPassword")
	logEntry.Info("Processing ForgotPassword")
	var model models.ForgotPasswordSend
	if err := c.BindJSON(&model); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err := au.GRPC_Client.ForgotPassword(model.Phone)
	if err != nil {
		logEntry.WithError(err).Error("Error During ForgotPassword RPC call")
		errs := response.ClientResponse(http.StatusBadRequest, "Could not send OTP", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	logEntry.Info("OTP sent successfully")
	success := response.ClientResponse(http.StatusOK, "OTP sent successfully", nil, nil)
	c.JSON(http.StatusOK, success)

}

// @Summary		Forgot password Verfy and Change
// @Description	user can change their password if user forgot the password and login
// @Tags			User OTP Login
// @Accept			json
// @Produce		    json
// @Param			model  body  models.ForgotVerify  true	"forgot-verify"
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/user/forgot-password-verify      [POST]
func (au *AuthHandler) ForgotPasswordVerifyAndChange(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "ForgotPasswordVerifyAndChange")
	logEntry.Info("Processing ForgotPasswordVerifyAndChange")
	var model models.ForgotVerify
	if err := c.BindJSON(&model); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	err := au.GRPC_Client.ForgotPasswordVerifyAndChange(model)
	if err != nil {
		logEntry.WithError(err).Error("Error During ForgotPasswordVerifyAndChange RPC call")
		errs := response.ClientResponse(http.StatusBadRequest, "Could not verify OTP", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	logEntry.Info("Successfully Changed the password")
	success := response.ClientResponse(http.StatusOK, "Successfully Changed the password", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary 		User Profile Details
// @Description 	User Details from User Profile
// @Tags 			User Profile
// @Accept 			json
// @Produce 		json
// @Security 		Bearer
// @Success 	200 {object} response.Response{}
// @Failure 	500 {object} response.Response{}
// @Router 		/user   [GET]
func (au *AuthHandler) UserDetails(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "UserDetails")
	logEntry.Info("Processing UserDetails")
	userID, _ := c.Get("user_id")
	UserDetails, err := au.AuthCachig.GetUserDetails(userID.(int))
	if err != nil {
		logEntry.WithError(err).Error("Error During UserDetails RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Get User Details")
	success := response.ClientResponse(http.StatusOK, "Successfully Get User Details", UserDetails, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary 	Specific User Details
// @Description User Details from User Profile
// @Tags 		User Profile
// @Accept 		json
// @Produce 	json
// @Security 	Bearer
// @Param 		user_id 	query 	string true 	"user id"
// @Success 	200 {object} response.Response{}
// @Failure 	500 {object} response.Response{}
// @Router 		/user/users   [GET]
func (au *AuthHandler) SpecificUserDetails(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "SpecificUserDetails")
	logEntry.Info("Processing SpecificUserDetails")
	userID := c.Query("user_id")
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Post_id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	UserDetails, err := au.AuthCachig.GetSpecificUserDetails(UserID)
	if err != nil {
		logEntry.WithError(err).Error("Error During SpecificUserDetails RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Get User Details")
	success := response.ClientResponse(http.StatusOK, "Successfully Get User Details", UserDetails, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary 	Update User Details
// @Description Update User Details by sending in user id
// @Tags 		User Profile
// @Accept 		json
// @Produce	 	json
// @Security	Bearer
// @Param 		firstname formData string true "First Name of the user"
// @Param 		lastname formData string true "Last Name of the user"
// @Param 		username formData string true "Username of the user"
// @Param 		dob formData string true "Date of Birth of the user (YYYY-MM-DD)"
// @Param 		gender formData string true "Gender of the user"
// @Param 		phone formData string false "Phone number of the user"
// @Param 		email formData string false "Email address of the user"
// @Param 		bio formData string false "Biography of the user"
// @Param 		photo formData file true "Photo of the user"
// @Success 	200 {object} response.Response{}
// @Failure 	500 {object} response.Response{}
// @Router 		/user [PUT]
func (au *AuthHandler) UpdateUserDetails(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "UpdateUserDetails")
	logEntry.Info("Processing UpdateUserDetails")
	firstname := c.PostForm("firstname")
	lastname := c.PostForm("lastname")
	username := c.PostForm("username")
	dob := c.PostForm("dob")
	gender := c.PostForm("gender")
	phone := c.PostForm("phone")
	email := c.PostForm("email")
	bio := c.PostForm("bio")
	if phone != "" {
		pattern := `^\d{10}$`
		regex := regexp.MustCompile(pattern)
		value := regex.MatchString(phone)
		if !value {
			fmt.Printf("%s This phone number is not valid", phone)
			return
		}
	}
	if email != "" {
		pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		regexpPattern := regexp.MustCompile(pattern)
		value := regexpPattern.MatchString(email)
		if !value {
			fmt.Printf("%s This email is not valid", phone)
			return
		}
	}
	user := models.UsersProfileDetail{
		Firstname: firstname,
		Lastname:  lastname,
		Username:  username,
		Dob:       dob,
		Gender:    gender,
		Phone:     phone,
		Email:     email,
		Bio:       bio,
	}
	err := validator.New().Struct(user)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Constraints not statisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "No file provided", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}

	userID, _ := c.Get("user_id")

	updateDetails, err := au.GRPC_Client.UpdateUserDetails(user, file, userID.(int))
	if err != nil {
		logEntry.WithError(err).Error("Error During UpdateUserDetails RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "failed update user", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Updated User Details")
	success := response.ClientResponse(http.StatusOK, "Successfully Updated User Details", updateDetails, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary 	Change User Password
// @Description Change User Password
// @Tags 		User Profile
// @Accept 		json
// @Produce 	json
// @Security 	Bearer
// @Param 		changePassword 	body   models.ChangePassword 	true 	"User Password Change"
// @Success 	200 {object} response.Response{}
// @Failure 	500 {object} response.Response{}
// @Router 		/user/changepassword     [PUT]
func (au *AuthHandler) ChangePassword(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "ChangePassword")
	logEntry.Info("Processing ChangePassword")
	user_id, _ := c.Get("user_id")
	var changePassword models.ChangePassword
	if err := c.BindJSON(&changePassword); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	if err := au.GRPC_Client.ChangePassword(user_id.(int), changePassword); err != nil {
		logEntry.WithError(err).Error("Error During ChangePassword RPC call")
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not change the password", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	logEntry.Info("Password changed Successfully")
	success := response.ClientResponse(http.StatusOK, "Password changed Successfully", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary 	Search Users
// @Description Search Users
// @Tags 		User Profile
// @Accept 		json
// @Produce 	json
// @Security 	Bearer
// @Param 		req 	body   models.SearchUser 	true 	"Search Details"
// @Success 	200 {object} response.Response{}
// @Failure 	500 {object} response.Response{}
// @Router 		/user/search     [GET]
func (au *AuthHandler) SearchUser(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "SearchUser")
	logEntry.Info("Processing SearchUser")
	var req models.SearchUser
	if err := c.ShouldBindJSON(&req); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
	}
	data, err := au.GRPC_Client.SearchUser(req)
	if err != nil {
		logEntry.WithError(err).Error("Error During SearchUser RPC call")
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	logEntry.Info("Successfully Searched User")
	success := response.ClientResponse(http.StatusOK, "Successfully Searched User", data, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary 	Report User
// @Description Report User to User
// @Tags 		Reports
// @Accept 		json
// @Produce 	json
// @Security 	Bearer
// @Param 		req 	body models.ReportRequest	 true 	"User Report"
// @Success 	200 {object} response.Response{}
// @Failure 	500 {object} response.Response{}
// @Router 		/report/user     [POST]
func (au *AuthHandler) ReportUser(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "ReportUser")
	logEntry.Info("Processing ReportUser")
	ReportedID, _ := c.Get("user_id")
	var req models.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
	}
	err := au.GRPC_Client.ReportUser(ReportedID.(int), req)
	if err != nil {
		logEntry.WithError(err).Error("Error During ReportUser RPC call")
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	logEntry.Info("Successfully Reported")
	success := response.ClientResponse(http.StatusOK, "Successfully Reported", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary 		Send Follow Request
// @Description 	Send Follow Request
// @Tags 			Follow
// @Accept 			json
// @Produce 		json
// @Security 		Bearer
// @Param 			user_id 	query 	string true 	"user id"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router 		/follow/request      [POST]
func (au *AuthHandler) FollowREQ(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "FollowREQ")
	logEntry.Info("Processing FollowREQ")
	userID, _ := c.Get("user_id")
	id := c.Query("user_id")
	FollowUserID, err := strconv.Atoi(id)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "UserID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = au.GRPC_Client.FollowREQ(userID.(int), FollowUserID)
	if err != nil {
		logEntry.WithError(err).Error("Error During FollowREQ RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "Details is incorrect", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully followed")
	sucess := response.ClientResponse(http.StatusOK, "Successfully followed", nil, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary 		Show Follow Request
// @Description 	Show Follow Request
// @Tags 			Follow
// @Accept 			json
// @Produce 		json
// @Security 		Bearer
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router 		/follow/requests      [GET]
func (au *AuthHandler) ShowFollowREQ(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "ShowFollowREQ")
	logEntry.Info("Processing ShowFollowREQ")
	userID, _ := c.Get("user_id")
	data, err := au.GRPC_Client.ShowFollowREQ(userID.(int))
	if err != nil {
		logEntry.WithError(err).Error("Error During ShowFollowREQ RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Show All Followed Request")
	sucess := response.ClientResponse(http.StatusOK, "Successfully Show All Followed Request", data, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary 		Accept Follow Request
// @Description 	Accept Follow Request
// @Tags 			Follow
// @Accept 			json
// @Produce 		json
// @Security 		Bearer
// @Param 		    user_id 	query 	string true 	"user id"
// @Success 	200 {object} response.Response{}
// @Failure 	500 {object} response.Response{}
// @Router		 /follow/accept      [POST]
func (au *AuthHandler) AcceptFollowREQ(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "AcceptFollowREQ")
	logEntry.Info("Processing AcceptFollowREQ")
	userID, _ := c.Get("user_id")
	id := c.Query("user_id")
	FollowingID, err := strconv.Atoi(id)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "UserID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = au.GRPC_Client.AcceptFollowREQ(userID.(int), FollowingID)
	if err != nil {
		logEntry.WithError(err).Error("Error During AcceptFollowREQ RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Accepted Following")
	sucess := response.ClientResponse(http.StatusOK, "Successfully Accepted Following", nil, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary 		User Unfollow
// @Description 	User Unfollow
// @Tags 			Follow
// @Accept 			json
// @Produce 		json
// @Security 		Bearer
// @Param 		 	user_id  	query 	string true 	"user id"
// @Success 	200 {object} response.Response{}
// @Failure		500 {object} response.Response{}
// @Router 		/follow/unfollow      [POST]
func (au *AuthHandler) UnFollow(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "UnFollow")
	logEntry.Info("Processing UnFollow")
	userID, _ := c.Get("user_id")
	id := c.Query("user_id")
	UnFollowUserID, err := strconv.Atoi(id)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "UserID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = au.GRPC_Client.UnFollow(userID.(int), UnFollowUserID)
	if err != nil {
		logEntry.WithError(err).Error("Error During UnFollow RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully UnFollowed")
	sucess := response.ClientResponse(http.StatusOK, "Successfully UnFollowed", nil, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary 		Show Followings
// @Description 	Show Followings
// @Tags 			Follow
// @Accept 			json
// @Produce 		json
// @Security 		Bearer
// @Success 	200 {object} response.Response{}
// @Failure 	500 {object} response.Response{}
// @Router 		/follow/following      [GET]
func (au *AuthHandler) Following(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "Following")
	logEntry.Info("Processing Get Following")
	userID, _ := c.Get("user_id")
	data, err := au.GRPC_Client.Following(userID.(int))
	if err != nil {
		logEntry.WithError(err).Error("Error During Following RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Retrive Followings")
	sucess := response.ClientResponse(http.StatusOK, "Successfully Retrive Followings", data, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary 		Show Followers
// @Description 	Show Followers
// @Tags 			Follow
// @Accept 			json
// @Produce 		json
// @Security 		Bearer
// @Success 	200 {object} response.Response{}
// @Failure 	500 {object} response.Response{}
// @Router 		/follow/followers      [GET]
func (au *AuthHandler) Follower(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "Follower")
	logEntry.Info("Processing Get Follower")
	userID, _ := c.Get("user_id")
	data, err := au.GRPC_Client.Follower(userID.(int))
	if err != nil {
		logEntry.WithError(err).Error("Error During Follower RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Retrive Followers")
	sucess := response.ClientResponse(http.StatusOK, "Successfully Retrive Followers", data, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary		Admin Login
// @Description	Login handler for Zsoxial admins
// @Tags			Admin
// @Accept 			json
// @Produce 		json
// @Security 		Bearer
// @Param			AdminLoginDetail	body		models.AdminLoginRequest	true	"Admin login details"
// @Success			200		{object}	response.Response{}
// @Failure			500		{object}	response.Response{}
// @Router			/admin/login  [POST]
func (au *AuthHandler) AdminLogin(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "AdminLogin")
	logEntry.Info("Processing Loggin request")
	var AdminLoginDetail models.AdminLoginRequest
	if err := c.ShouldBindJSON(&AdminLoginDetail); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
	}
	err := validator.New().Struct(AdminLoginDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Constraints not statisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	admin, err := au.GRPC_Client.AdminLogin(AdminLoginDetail)
	if err != nil {
		logEntry.WithError(err).Error("Error During AdminLogin RPC call")
		errs := response.ClientResponse(http.StatusBadRequest, "Cannot Authenthicate Admin", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	logEntry.Info("Admin successfully logged in with password")
	success := response.ClientResponse(http.StatusOK, "Admin successfully logged in with password", admin, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Show All Users
// @Description	Retrieve users with pagination
// @Tags			Admin User Management
// @Accept 			json
// @Produce 		json
// @Security 		Bearer
// @Param 	page 	query 	string	 false	 "Page number"
// @Param 	count 	query 	string 	false	 "Page size"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/users   [GET]
func (au *AuthHandler) ShowAllUsers(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "ShowAllUsers")
	logEntry.Info("Processing ShowAllUsers")
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	countStr := c.DefaultQuery("count", "100")
	pageSize, err := strconv.Atoi(countStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user count in a page not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	users, err := au.AuthCachig.ShowAllUsers(page, pageSize)
	if err != nil {
		logEntry.WithError(err).Error("Error During ShowAllUsers RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "couldn't retrieve users", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Retrieved all Users")
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all Users", users, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Block an existing User
// @Description	Using this handler admins can block an user
// @Tags			Admin User Management
// @Accept			json
// @Produce			json
// @Security		Bearer
// @Param			id	query	string	true	"user id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/user/block   [PUT]
func (au *AuthHandler) BlockUser(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "BlockUser")
	logEntry.Info("Processing BlockUser")
	id := c.Query("id")
	userID, _ := strconv.Atoi(id)
	err := au.GRPC_Client.AdminBlockUser(userID)
	if err != nil {
		logEntry.WithError(err).Error("Error During BlockUser RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be blocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully blocked the user")
	sucess := response.ClientResponse(http.StatusOK, "Successfully blocked the user", nil, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary		 UnBlock an existing user
// @Description	 Using this handler admins can block an user
// @Tags			Admin User Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param			id	query		string	true	"user id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/user/unblock    [PUT]
func (au *AuthHandler) UnBlockUser(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "UnBlockUser")
	logEntry.Info("Processing UnBlockUser")
	id := c.Query("id")
	userID, _ := strconv.Atoi(id)
	err := au.GRPC_Client.AdminUnblockUser(userID)
	if err != nil {
		logEntry.WithError(err).Error("Error During UnBlockUser RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully unblocked the user")
	sucess := response.ClientResponse(http.StatusOK, "Successfully unblocked the user", nil, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary			Get User Reports
// @Description		Retrieve UserReports with pagination
// @Tags			Admin Reports Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param 	page 	query 	string 	false 	"Page number"
// @Param 	count 	query 	string 	false 	"Page size"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/report/user   [GET]
func (au *AuthHandler) ShowUserReports(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "ShowUserReports")
	logEntry.Info("Processing ShowUserReports")
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	countStr := c.DefaultQuery("count", "100")
	pageSize, err := strconv.Atoi(countStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user count in a page not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	users, err := au.GRPC_Client.ShowUserReports(page, pageSize)
	if err != nil {
		logEntry.WithError(err).Error("Error During ShowUserReports RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "couldn't retrieve users", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Retrieved all UserReports")
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all UserReports", users, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary			Get Post Reports
// @Description		Retrieve UserReports with pagination
// @Tags			Admin Reports Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param 	page 	query 	string 	false 	"Page number"
// @Param 	count 	query 	string 	false 	"Page size"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/report/post   [GET]
func (au *AuthHandler) ShowPostReports(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "ShowPostReports")
	logEntry.Info("Processing ShowPostReports")
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	countStr := c.DefaultQuery("count", "100")
	pageSize, err := strconv.Atoi(countStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user count in a page not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	users, err := au.GRPC_Client.ShowPostReports(page, pageSize)
	if err != nil {
		logEntry.WithError(err).Error("Error During ShowPostReports RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "couldn't retrieve users", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Retrieved all PostReports")
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all PostReports", users, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary			Get Posts
// @Description		Retrieve posts with pagination
// @Tags			Admin Post Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param 		page query string false "Page number"
// @Param 		count query string false "Page size"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/admin/posts   [GET]
func (au *AuthHandler) GetAllPosts(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "GetAllPosts")
	logEntry.Info("Processing GetAllPosts")
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	countStr := c.DefaultQuery("count", "100")
	pageSize, err := strconv.Atoi(countStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "post count in a page not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	users, err := au.GRPC_Client.GetAllPosts(page, pageSize)
	if err != nil {
		logEntry.WithError(err).Error("Error During GetAllPosts RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "couldn't retrieve posts", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Retrieved all Posts")
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all Posts", users, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Remove Post
// @Description	 Admin can delete a post
// @Tags			Admin Post Management
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param			post_id	query	string	true	"post id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/admin/post     [DELETE]
func (au *AuthHandler) RemovePost(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "RemovePost")
	logEntry.Info("Processing RemovePost")
	postID := c.Query("post_id")
	PostID, err := strconv.Atoi(postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "PostID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = au.GRPC_Client.RemovePost(PostID)
	if err != nil {
		logEntry.WithError(err).Error("Error During RemovePost RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't Remove Post", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Removed Post")
	sucess := response.ClientResponse(http.StatusOK, "Successfully Removed Post", nil, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary		Generate Key For VideoCall
// @Description	Generate Key For VideoCall
// @Tags			VideoCall
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param			user	query	string	true	"user"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/videocall/key     [GET]
func (au *AuthHandler) VideoCallKey(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "VideoCallKey")
	logEntry.Info("Processing VideoCallKey")
	userID, _ := c.Get("user_id")
	UserID := c.Query("user")
	oppositeUser, err := strconv.Atoi(UserID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "PostID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	key, err := au.GRPC_Client.VideoCallKey(userID.(int), oppositeUser)
	if err != nil {
		logEntry.WithError(err).Error("Error During VideoCallKey RPC call")
		errs := response.ClientResponse(http.StatusInternalServerError, "Couldn't not reterive link", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	logEntry.Info("Successfully Get a VideoCallKey And Private Link")
	url := fmt.Sprintf("https://zsoxial.zhooze.shop/index?room=%s", key)
	sucess := response.ClientResponse(http.StatusOK, "Successfully Get a VideoCallKey And Private Link", url, nil)
	c.JSON(http.StatusOK, sucess)
}

// @Summary		Create Groups
// @Description	Multiple Users Can together
// @Tags			Groups
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param name formData string true "Name of the Group"
// @Param description formData string true "Description of the Group"
// @Param           photo formData file true "Photo of the post"
// @Param user formData array true "Users associated with the Group. Provide multiple user IDs"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/group     [post]
func (au *AuthHandler) CreateGroup(c *gin.Context) {
	userID, _ := c.Get("user_id")
	groupName := c.PostForm("name")
	description := c.PostForm("description")
	user := c.PostFormArray("user")
	fmt.Println("dat", user)
	users, err := helper.ConvertStringToArray(user)
	fmt.Println("dadddt", users)
	if err != nil {
		return
	}
	req := models.GroupReq{
		Name:        groupName,
		Description: description,
	}
	file, err := c.FormFile("photo")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "No file provided", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	err = au.GRPC_Client.CreateGroup(userID.(int), req, users, file)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusCreated, "Successfully Created Group", nil, nil)
	c.JSON(http.StatusCreated, success)
}

// @Summary		Exit From Groups
// @Description	Exit From Groups
// @Tags			Groups
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param			group_id	query	string	true	"group id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/group     [DELETE]
func (au *AuthHandler) ExitFormGroup(c *gin.Context) {
	userID, _ := c.Get("user_id")
	groupID := c.Query("group_id")
	GroupID, err := strconv.Atoi(groupID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "GroupID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = au.GRPC_Client.ExitFormGroup(userID.(int), GroupID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Exit From Group", nil, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Show All Groups
// @Description	Show All Groups
// @Tags			Groups
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/group     [GET]
func (au *AuthHandler) ShowGroups(c *gin.Context) {
	userID, _ := c.Get("user_id")
	data, err := au.GRPC_Client.ShowGroups(userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Get Group", data, nil)
	c.JSON(http.StatusOK, success)
}

// @Summary		Show All Members From Group
// @Description	Show All Members From Group
// @Tags			Groups
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param			group_id	query	string	true	"group id"
// @Success		200	{object}	response.Response{}
// @Failure		500	{object}	response.Response{}
// @Router			/group/members     [GET]
func (au *AuthHandler) ShowGroupMembers(c *gin.Context) {
	userID, _ := c.Get("user_id")
	groupID := c.Query("group_id")
	GroupID, err := strconv.Atoi(groupID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "GroupID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	data, err := au.GRPC_Client.ShowGroupMembers(userID.(int), GroupID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Get Members From Group", data, nil)
	c.JSON(http.StatusOK, success)
}
