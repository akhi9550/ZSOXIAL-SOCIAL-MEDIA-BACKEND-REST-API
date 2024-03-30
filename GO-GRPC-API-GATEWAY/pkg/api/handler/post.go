package handler

import (
	"mime/multipart"
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
	user1 := c.PostForm("1")
	user2 := c.PostForm("2")
	user3 := c.PostForm("3")
	user4 := c.PostForm("4")
	user5 := c.PostForm("5")

	posttype, _ := strconv.Atoi(typeid)
	one, _ := strconv.Atoi(user1)
	two, _ := strconv.Atoi(user2)
	three, _ := strconv.Atoi(user3)
	fore, _ := strconv.Atoi(user4)
	five, _ := strconv.Atoi(user5)

	req := models.PostRequest{
		Caption: caption,
		TypeId:  uint(posttype),
	}

	users := models.Tags{
		User1: uint(one),
		User2: uint(two),
		User3: uint(three),
		User4: uint(fore),
		User5: uint(five),
	}
	form, err := c.MultipartForm()
	if err != nil {
		errorRes := response.ClientResponse(http.StatusBadRequest, "Retrieving images from form error", nil, err.Error())
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	files := form.File["photo"]
	if len(files) == 0 {
		errorRes := response.ClientResponse(http.StatusBadRequest, "No files provided", nil, nil)
		c.JSON(http.StatusBadRequest, errorRes)
		return
	}
	var fileHeaders []*multipart.FileHeader
	for _, file := range files {
		fileHeaders = append(fileHeaders, file)
	}
	data, err := p.GRPC_Client.CreatePost(userID.(int), req, fileHeaders, users)
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
	var req models.UpdatePostReq
	if err := c.BindJSON(&req); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "fields provided are in wrong format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	data, err := p.GRPC_Client.UpdatePost(userID.(int), req)
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
	err = p.GRPC_Client.DeletePost(userID.(int), int(PostID))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully Deleted Post", nil, nil)
	c.JSON(http.StatusOK, success)
}
