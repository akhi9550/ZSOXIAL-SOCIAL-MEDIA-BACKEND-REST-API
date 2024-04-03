package handler

import (
	"net/http"
	"strconv"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/akhi9550/api-gateway/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	GRPC_Client interfaces.PostClient
}

func NewPostHandler(postClient interfaces.PostClient) *PostHandler {
	return &PostHandler{
		GRPC_Client: postClient,
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

func (p *PostHandler) GetPost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Query("post_id")
	PostID, err := strconv.Atoi(postID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Post_id not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	data, err := p.GRPC_Client.GetPost(userID.(int), int(PostID))
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
	userID, _ := c.Get("user_id")
	data, err := p.GRPC_Client.GetAllPost(userID.(int))
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
	data, err := p.GRPC_Client.GetAllArchivePost(userID.(int))
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
	data,err := p.GRPC_Client.GetSavedPost(userID.(int))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Received Saved Post", data, nil)
	c.JSON(http.StatusOK, success)

}
