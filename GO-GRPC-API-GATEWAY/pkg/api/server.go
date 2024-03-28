package server

import (
	"log"

	"github.com/akhi9550/api-gateway/pkg/api/handler"
	"github.com/akhi9550/api-gateway/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(authHandler *handler.AuthHandler) *ServerHTTP {
	r := gin.New()

	r.Use(gin.Logger())
	r.POST("/login", authHandler.AdminLogin)

	// r.Use(middleware.AdminAuthMiddleware())
	// {
	// 	r.GET("admin/users", authHandler.ShowAllUsers)
	// 	r.PUT("admin/user/block", authHandler.BlockUser)
	// 	r.PUT("admin/user/unblock", authHandler.UnBlockUser)
	// }

	r.POST("user/signup", authHandler.UserSignup)
	r.POST("user/login", authHandler.Userlogin)

	r.POST("user/send-otp", authHandler.SendOtp)
	r.POST("user/verify-otp", authHandler.VerifyOtp)

	r.POST("user/forgot-password", authHandler.ForgotPassword)
	r.POST("user/forgot-password-verify", authHandler.ForgotPasswordVerifyAndChange)

	r.Use(middleware.UserAuthMiddleware())
	{
		r.GET("user/users", authHandler.UserDetails)
		r.PUT("user/update", authHandler.UpdateUserDetails)
		r.PUT("user/changepassword", authHandler.ChangePassword)
	}
	return &ServerHTTP{engine: r}
}
func (s *ServerHTTP) Start() {
	log.Printf("Starting Server on 8000")
	err := s.engine.Run(":8000")
	if err != nil {
		log.Printf("error while starting the server")
	}
}
