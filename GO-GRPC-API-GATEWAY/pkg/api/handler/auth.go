package handler

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/akhi9550/api-gateway/pkg/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	GRPC_Client interfaces.AuthClient
}

func NewAuthHandler(authClient interfaces.AuthClient) *AuthHandler {
	return &AuthHandler{
		GRPC_Client: authClient,
	}
}

func (au *AuthHandler) UserSignup(c *gin.Context) {
	var UserSignupDetail models.UserSignUpRequest

	if err := c.ShouldBindJSON(&UserSignupDetail); err != nil {
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
		errs := response.ClientResponse(http.StatusBadRequest, "Constraints not statisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	user, err := au.GRPC_Client.UserSignUp(UserSignupDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	fmt.Println("", user)
	success := response.ClientResponse(http.StatusCreated, "User successfully signed up", user, nil)
	c.JSON(http.StatusCreated, success)
}

func (au *AuthHandler) Userlogin(c *gin.Context) {
	var UserLoginDetail models.UserLoginRequest
	if err := c.ShouldBindJSON(&UserLoginDetail); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
	}
	err := validator.New().Struct(UserLoginDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Constraints not statisfied", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	user, err := au.GRPC_Client.UserLogin(UserLoginDetail)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusCreated, "User successfully logged in with password", user, nil)
	c.JSON(http.StatusCreated, success)
}

func (au *AuthHandler) SendOtp(c *gin.Context) {
	var phone models.OTPData
	if err := c.ShouldBindJSON(&phone); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err := au.GRPC_Client.SendOtp(phone.PhoneNumber)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "user with this phone is not exists", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "OTP sent successfully", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (au *AuthHandler) VerifyOtp(c *gin.Context) {
	var code models.VerifyData
	if err := c.ShouldBindJSON(&code); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	user, err := au.GRPC_Client.VerifyOTP(code)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "Could not verify OTP", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	sucess := response.ClientResponse(http.StatusOK, "Successfully verified OTP", user, nil)
	c.JSON(http.StatusOK, sucess)
}

func (au *AuthHandler) ForgotPassword(c *gin.Context) {
	var model models.ForgotPasswordSend
	if err := c.BindJSON(&model); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err := au.GRPC_Client.ForgotPassword(model.Phone)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Could not send OTP", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	success := response.ClientResponse(http.StatusOK, "OTP sent successfully", nil, nil)
	c.JSON(http.StatusOK, success)

}

func (au *AuthHandler) ForgotPasswordVerifyAndChange(c *gin.Context) {
	var model models.ForgotVerify
	if err := c.BindJSON(&model); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	err := au.GRPC_Client.ForgotPasswordVerifyAndChange(model)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Could not verify OTP", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	success := response.ClientResponse(http.StatusOK, "Successfully Changed the password", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (au *AuthHandler) UserDetails(c *gin.Context) {
	userID, _ := c.Get("user_id")
	UserDetails, err := au.GRPC_Client.UserDetails(userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "failed to retrieve details", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "User Details", UserDetails, nil)
	c.JSON(http.StatusOK, success)
}

func (au *AuthHandler) UpdateUserDetails(c *gin.Context) {
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
		errs := response.ClientResponse(http.StatusInternalServerError, "failed update user", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Updated User Details", updateDetails, nil)
	c.JSON(http.StatusOK, success)
}

func (au *AuthHandler) ChangePassword(c *gin.Context) {
	user_id, _ := c.Get("user_id")
	var changePassword models.ChangePassword
	if err := c.BindJSON(&changePassword); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	if err := au.GRPC_Client.ChangePassword(user_id.(int), changePassword); err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Could not change the password", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	success := response.ClientResponse(http.StatusOK, "password changed Successfully ", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (au *AuthHandler) AdminLogin(c *gin.Context) {
	var AdminLoginDetail models.AdminLoginRequest
	if err := c.ShouldBindJSON(&AdminLoginDetail); err != nil {
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
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusCreated, "Admin successfully logged in with password", admin, nil)
	c.JSON(http.StatusCreated, success)
}

func (au *AuthHandler) ShowAllUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "page number not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	countStr := c.DefaultQuery("count", "10")
	pageSize, err := strconv.Atoi(countStr)
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "user count in a page not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	users, err := au.GRPC_Client.ShowAllUsers(page, pageSize)
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "couldn't retrieve users", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Retrieved all Users", users, nil)
	c.JSON(http.StatusOK, success)
}

func (au *AuthHandler) BlockUser(c *gin.Context) {
	id := c.Query("id")
	userID, _ := strconv.Atoi(id)
	err := au.GRPC_Client.AdminBlockUser(int(userID))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be blocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	sucess := response.ClientResponse(http.StatusOK, "Successfully blocked the user", nil, nil)
	c.JSON(http.StatusOK, sucess)
}

func (au *AuthHandler) UnBlockUser(c *gin.Context) {
	id := c.Query("id")
	userID, _ := strconv.Atoi(id)
	err := au.GRPC_Client.AdminUnblockUser(int(userID))
	if err != nil {
		errs := response.ClientResponse(http.StatusInternalServerError, "user could not be unblocked", nil, err.Error())
		c.JSON(http.StatusInternalServerError, errs)
		return
	}
	sucess := response.ClientResponse(http.StatusOK, "Successfully unblocked the user", nil, nil)
	c.JSON(http.StatusOK, sucess)
}

func (au *AuthHandler) ReportUser(c *gin.Context) {
	ReportedID, _ := c.Get("user_id")
	var req models.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
	}
	err := au.GRPC_Client.ReportUser(ReportedID.(int), req)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Reported", nil, nil)
	c.JSON(http.StatusOK, success)
}
