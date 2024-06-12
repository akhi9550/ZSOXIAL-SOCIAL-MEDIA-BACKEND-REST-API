package server

import (
	"log"

	"github.com/akhi9550/api-gateway/pkg/api/handler"
	"github.com/akhi9550/api-gateway/pkg/api/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(authHandler *handler.AuthHandler, postHandler *handler.PostHandler, chatHandler *handler.ChatHandler, notificationHandler *handler.NotificationHandler, videocallHandler *handler.VideoCallHandler) *ServerHTTP {
	r := gin.New()

	r.Use(gin.Logger())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Static("/static", "./static")
	r.LoadHTMLGlob("template/*")

	r.GET("/exit", videocallHandler.ExitPage)
	r.GET("/error", videocallHandler.ErrorPage)
	r.GET("/index", videocallHandler.IndexedPage)

	r.POST("admin/login", authHandler.AdminLogin)

	r.GET("admin/users", middleware.AdminAuthMiddleware(), authHandler.ShowAllUsers)
	r.PUT("admin/user/block", middleware.AdminAuthMiddleware(), authHandler.BlockUser)
	r.PUT("admin/user/unblock", middleware.AdminAuthMiddleware(), authHandler.UnBlockUser)
	r.GET("/admin/report/user", middleware.AdminAuthMiddleware(), authHandler.ShowUserReports)
	r.GET("/admin/report/post", middleware.AdminAuthMiddleware(), authHandler.ShowPostReports)
	r.GET("/admin/posts", middleware.AdminAuthMiddleware(), authHandler.GetAllPosts)
	r.DELETE("/admin/post", middleware.AdminAuthMiddleware(), authHandler.RemovePost)
	r.POST("/admin/post/type", middleware.AdminAuthMiddleware(), postHandler.CreatePostType)
	r.GET("/admin/post/type", middleware.AdminAuthMiddleware(), postHandler.ShowPostType)
	r.DELETE("/admin/post/type", middleware.AdminAuthMiddleware(), postHandler.DeletePostType)

	r.POST("user/signup", authHandler.UserSignup)
	r.POST("user/login", authHandler.Userlogin)

	r.POST("user/send-otp", authHandler.SendOtp)
	r.POST("user/verify-otp", authHandler.VerifyOtp)

	r.POST("user/forgot-password", authHandler.ForgotPassword)
	r.POST("user/forgot-password-verify", authHandler.ForgotPasswordVerifyAndChange)

	r.Use(middleware.UserAuthMiddleware())
	{
		user := r.Group("/user")
		{
			user.GET("", authHandler.UserDetails)
			user.PUT("", authHandler.UpdateUserDetails)
			user.GET("/users", authHandler.SpecificUserDetails)
			user.PUT("/changepassword", authHandler.ChangePassword)
			user.GET("/search", authHandler.SearchUser)
		}

		report := r.Group("/report")
		{
			report.POST("/user", authHandler.ReportUser)
			report.POST("/post", postHandler.ReportPost)
		}

		follow := r.Group("/follow")
		{
			follow.POST("/request", authHandler.FollowREQ)
			follow.GET("/requests", authHandler.ShowFollowREQ)
			follow.POST("/accept", authHandler.AcceptFollowREQ)
			follow.POST("/unfollow", authHandler.UnFollow)
			follow.GET("/following", authHandler.Following)
			follow.GET("/followers", authHandler.Follower)
		}

		post := r.Group("/post")
		{
			post.POST("", postHandler.CreatePost)
			post.GET("", postHandler.GetUserPost)
			post.PUT("", postHandler.UpdatePost)
			post.DELETE("", postHandler.DeletePost)
			post.GET("/type", postHandler.ShowPostTypeUser)
			post.GET("/posts", postHandler.GetPost)
			post.GET("/getposts", postHandler.GetAllPost)
			post.GET("/home", postHandler.Home)
		}

		savepost := r.Group("/saved")
		{
			savepost.POST("", postHandler.SavedPost)
			savepost.GET("", postHandler.GetSavedPost)
			savepost.POST("/unsaved", postHandler.UnSavedPost)
		}

		archive := r.Group("/archive")
		{
			archive.POST("", postHandler.ArchivePost)
			archive.GET("", postHandler.GetAllArchivePost)
			archive.POST("/unarchive", postHandler.UnArchivePost)

		}

		like := r.Group("/like")
		{
			like.PUT("", postHandler.LikePost)
			like.PUT("/unlike", postHandler.UnLikePost)
		}

		comment := r.Group("/comment")
		{
			comment.POST("", postHandler.PostComment)
			comment.DELETE("", postHandler.DeleteComment)
			comment.POST("/reply", postHandler.ReplyComment)
			comment.GET("", postHandler.GetAllPostComments)
			comment.GET("/showcomments", postHandler.ShowAllPostComments)
		}

		story := r.Group("/story")
		{
			story.POST("", postHandler.CreateStory)
			story.GET("", postHandler.GetStory)
			story.DELETE("", postHandler.DeleteStory)
			story.PUT("/like", postHandler.LikeStory)
			story.PUT("/unlike", postHandler.UnLikeStory)
			story.GET("/details", postHandler.StoryDetails)
		}

		chat := r.Group("/chat")
		{
			chat.GET("", chatHandler.FriendMessage)
			chat.GET("/message", chatHandler.GetChat)
		}
		notification := r.Group("/notification")
		{
			notification.GET("", notificationHandler.GetNotification)
		}
		videoCall := r.Group("/videocall")
		{
			videoCall.GET("/key", authHandler.VideoCallKey)
		}
		group:=r.Group("/group")
		{
			group.POST("",authHandler.CreateGroup)
			group.DELETE("",authHandler.ExitFormGroup)
			group.GET("",authHandler.ShowGroups)
			group.GET("/members",authHandler.ShowGroupMembers)
		}
	}

	return &ServerHTTP{engine: r}
}

func (s *ServerHTTP) Start() {
	log.Printf("Starting Server on 8080")
	err := s.engine.Run(":8080")
	if err != nil {
		log.Printf("error while starting the server")
	}
}
