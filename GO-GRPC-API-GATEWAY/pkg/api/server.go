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

func NewServerHTTP(authHandler *handler.AuthHandler, postHandler *handler.PostHandler) *ServerHTTP {
	r := gin.New()

	r.Use(gin.Logger())
	r.POST("/login", authHandler.AdminLogin)

	// r.Use(middleware.AdminAuthMiddleware())
	r.GET("admin/users", middleware.AdminAuthMiddleware(), authHandler.ShowAllUsers)
	r.PUT("admin/user/block", middleware.AdminAuthMiddleware(), authHandler.BlockUser)
	r.PUT("admin/user/unblock", middleware.AdminAuthMiddleware(), authHandler.UnBlockUser)

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
		r.POST("/reportuser", authHandler.ReportUser)
		r.POST("/post", postHandler.CreatePost)
		r.GET("/getpost", postHandler.GetPost)
		r.PUT("/updatepost", postHandler.UpdatePost)
		r.DELETE("/deletepost", postHandler.DeletePost)
		r.GET("/getallpost", postHandler.GetAllPost)
		r.POST("/archivePost", postHandler.ArchivePost)
		r.POST("/unarchivePost", postHandler.UnArchivePost)
		r.GET("/getarchivePost", postHandler.GetAllArchivePost)
		r.PUT("/like", postHandler.LikePost)
		r.PUT("/unlike", postHandler.UnLinkPost)
		r.POST("/comment", postHandler.PostComment)
		r.POST("/replycomment", postHandler.ReplyComment)
		r.DELETE("/comment", postHandler.DeleteComment)
		r.GET("/allcomment", postHandler.GetAllPostComments)
		r.GET("/showcomments", postHandler.ShowAllPostComments)
		r.POST("/reportpost", postHandler.ReportPost)
		r.DELETE("/saved", postHandler.SavedPost)
		r.POST("/unsaved", postHandler.UnSavedPost)
		r.GET("/getsavedpost", postHandler.GetSavedPost)
		r.POST("/createstory", postHandler.CreateStory)
		r.GET("/getstories", postHandler.GetStory)
		r.DELETE("/deletestory", postHandler.DeleteStory)
		r.PUT("/likestory", postHandler.LikeStory)
		r.PUT("/unlikestory", postHandler.UnLikeStory)
		r.PUT("/followreq", authHandler.FollowREQ)
		r.GET("/showreq", authHandler.ShowFollowREQ)
		r.PUT("/followeraccept", authHandler.AcceptFollowREQ)
		r.GET("/following",authHandler.Following)
		r.GET("/follower",authHandler.Follower)
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
