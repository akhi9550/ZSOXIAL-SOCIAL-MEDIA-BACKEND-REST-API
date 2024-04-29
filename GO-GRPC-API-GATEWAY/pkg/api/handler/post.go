package handler

import (
	"fmt"
	"net/http"
	"strconv"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/helper"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/akhi9550/api-gateway/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	GRPC_Client interfaces.PostClient
	PostCachig  *helper.RedisPostCaching
}

func NewPostHandler(postClient interfaces.PostClient, postCaching *helper.RedisPostCaching) *PostHandler {
	return &PostHandler{
		GRPC_Client: postClient,
		PostCachig:  postCaching,
	}
}

func (p *PostHandler) CreatePost(c *gin.Context) {
	userID, _ := c.Get("user_id")

	caption := c.PostForm("caption")
	typeid := c.PostForm("posttype")
	user := c.PostFormArray("user")

	for _, i := range user {
		if i == "1" {
			return
		}
	}

	posttype, _ := strconv.Atoi(typeid)

	req := models.PostRequest{
		Caption: caption,
		TypeId:  uint(posttype),
	}

	file, err := c.FormFile("photo")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "No file provided", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	data, err := p.GRPC_Client.CreatePost(userID.(int), req, file, user)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusCreated, "User successfully Posted", data, nil)
	c.JSON(http.StatusCreated, success)
}

func (p *PostHandler) GetUserPost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	data, err := p.PostCachig.GetUserPost(userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Get a Post", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) GetPost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Query("post_id")
	PostID, err := strconv.Atoi(postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "PostID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	data, err := p.PostCachig.GetPost(userID.(int), int(PostID))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Get a Post", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) UpdatePost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postid := c.PostForm("post_id")
	caption := c.PostForm("caption")
	typeid := c.PostForm("posttype")
	user := c.PostFormArray("user")

	for _, i := range user {
		if i == "1" {
			return
		}
	}

	posttype, _ := strconv.Atoi(typeid)
	postID, _ := strconv.Atoi(postid)

	req := models.UpdatePostReq{
		PostID:  uint(postID),
		Caption: caption,
		TypeID:  uint(posttype),
	}

	data, err := p.GRPC_Client.UpdatePost(userID.(int), req, user)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusCreated, "User successfully Updated", data, nil)
	c.JSON(http.StatusCreated, success)
}

func (p *PostHandler) DeletePost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Query("post_id")
	PostID, err := strconv.Atoi(postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Post_id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = p.GRPC_Client.DeletePost(userID.(int), PostID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Deleted Post", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) GetAllPost(c *gin.Context) {
	userID := c.Query("user_id")
	UserID, err := strconv.Atoi(userID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "UserID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	data, err := p.PostCachig.GetAllPost(UserID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Get a Post", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) ArchivePost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Query("post_id")
	PostID, err := strconv.Atoi(postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Post_id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = p.GRPC_Client.ArchivePost(userID.(int), PostID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Archive Post", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) UnArchivePost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Query("post_id")
	PostID, err := strconv.Atoi(postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Post_id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = p.GRPC_Client.UnArchivePost(userID.(int), PostID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully UnArchive Post", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) GetAllArchivePost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	data, err := p.PostCachig.GetAllArchivePost(userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Get a Post", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) LikePost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Query("post_id")
	PostID, err := strconv.Atoi(postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Post_id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	data, err := p.GRPC_Client.LikePost(userID.(int), PostID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	success := response.ClientResponse(http.StatusOK, "Successfully Liked Post", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) UnLinkPost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Query("post_id")
	PostID, err := strconv.Atoi(postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Post_id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = p.GRPC_Client.UnLinkPost(userID.(int), PostID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Successfully Liked Post", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully UnLiked Post", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) PostComment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req models.PostCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
	}
	data, err := p.GRPC_Client.PostComment(userID.(int), req)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Commented Post", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) GetAllPostComments(c *gin.Context) {
	postID := c.Query("post_id")
	PostID, err := strconv.Atoi(postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Post_id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	data, err := p.PostCachig.GetAllPostComments(PostID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Received All Comments in the Post", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) DeleteComment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	comment := c.Query("comment_id")
	commentID, err := strconv.Atoi(comment)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "CommentID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = p.GRPC_Client.DeleteComment(userID.(int), commentID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Deleted Comment in the Post", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) ReplyComment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req models.ReplyCommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
	}
	data, err := p.GRPC_Client.ReplyComment(userID.(int), req)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Received Comment in the Post", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) SavedPost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	PostID := c.Query("post_id")
	postID, err := strconv.Atoi(PostID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Post_id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = p.GRPC_Client.SavedPost(userID.(int), postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Saved Post", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) UnSavedPost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	PostID := c.Query("post_id")
	postID, err := strconv.Atoi(PostID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Post_id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = p.GRPC_Client.UnSavedPost(userID.(int), postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully UnSaved Post", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) GetSavedPost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	data, err := p.PostCachig.GetSavedPost(userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	fmt.Println("darta", data)
	success := response.ClientResponse(http.StatusOK, "Successfully Received Saved Post", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) CreateStory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	file, err := c.FormFile("photo")
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "No file provided", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	data, err := p.GRPC_Client.CreateStory(userID.(int), file)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusCreated, "User successfully Posted Stroy", data, nil)
	c.JSON(http.StatusCreated, success)
}

func (p *PostHandler) GetStory(c *gin.Context) {
	viewUser, _ := c.Get("user_id")
	UserID := c.Query("user_id")
	userID, err := strconv.Atoi(UserID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Post_id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	data, err := p.PostCachig.GetStory(userID, viewUser.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "User successfully Recived All Stroy", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) DeleteStory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	StoryID := c.Query("story_id")
	storyID, err := strconv.Atoi(StoryID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "StoryID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = p.GRPC_Client.DeleteStory(userID.(int), storyID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "User successfully Deleted Stroy", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) LikeStory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	StoryID := c.Query("story_id")
	storyID, err := strconv.Atoi(StoryID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "StoryID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = p.GRPC_Client.LikeStory(userID.(int), storyID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "User successfully Liked Stroy", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) UnLikeStory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	StoryID := c.Query("story_id")
	storyID, err := strconv.Atoi(StoryID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "StoryID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	err = p.GRPC_Client.UnLikeStory(userID.(int), storyID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully UnLiked Stroy", nil, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) StoryDetails(c *gin.Context) {
	userID, _ := c.Get("user_id")
	StoryID := c.Query("story_id")
	storyID, err := strconv.Atoi(StoryID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "StoryID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	data, err := p.PostCachig.StoryDetails(userID.(int), storyID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Retrive Stroy Details", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) ShowAllPostComments(c *gin.Context) {
	PostID := c.Query("post_id")
	postID, err := strconv.Atoi(PostID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "PostID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	data, err := p.PostCachig.ShowAllPostComments(postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "User successfully Recived All Comments", data, nil)
	c.JSON(http.StatusOK, success)
}

func (p *PostHandler) ReportPost(c *gin.Context) {
	ReportedID, _ := c.Get("user_id")
	var req models.ReportPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
	}
	err := p.GRPC_Client.ReportPost(ReportedID.(int), req)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Reported", nil, nil)
	c.JSON(http.StatusOK, success)
}
