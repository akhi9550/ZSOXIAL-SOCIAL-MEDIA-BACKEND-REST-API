package handler

import (
	"fmt"
	"net/http"
	"strconv"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/helper"
	"github.com/akhi9550/api-gateway/pkg/logging"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/akhi9550/api-gateway/pkg/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var User = make(map[string]*websocket.Conn)

type ChatHandler struct {
	GRPC_Client interfaces.ChatClient
	helper      *helper.Helper
}

func NewChatHandler(chatClient interfaces.ChatClient, helper *helper.Helper) *ChatHandler {
	return &ChatHandler{
		GRPC_Client: chatClient,
		helper:      helper,
	}
}

// @Summary		    WebSocket Chat
// @Description		Establish a WebSocket connection for real-time chat messaging. This endpoint allows users to send and receive messages in real time.
// @Tags			Chat
// @Router			/zsoxial.zhooze.shop/chat   [GET]
func (ch *ChatHandler) FriendMessage(c *gin.Context) {
	fmt.Println("message called")
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Websocket Connection Issue", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		errs := response.ClientResponse(http.StatusBadRequest, "User ID in JWT claims is not a string", nil, "")
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	defer delete(User, strconv.Itoa(userID.(int)))
	defer conn.Close()
	user := strconv.Itoa(userID.(int))
	User[user] = conn

	for {
		fmt.Println("loop starts", userID, User)
		_, msg, err := conn.ReadMessage()
		if err != nil {
			errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
			c.JSON(http.StatusBadRequest, errs)
			return
		}
		ch.helper.SendMessageToUser(User, msg, user)
	}
}

// @Summary			Get Users Chats
// @Description		Retrieve UsersChats
// @Tags			Chat
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param			chatRequest  	body		models.ChatRequest	true	"GetChat details"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/chat/message   [GET]
func (ch *ChatHandler) GetChat(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "GetChat")
	logEntry.Info("Processing GetChat")
	var chatRequest models.ChatRequest
	if err := c.ShouldBindJSON(&chatRequest); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		errs := response.ClientResponse(http.StatusBadRequest, "User ID not found in JWT claims", nil, "")
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	userID := strconv.Itoa(userIDInterface.(int))
	result, err := ch.GRPC_Client.GetChat(userID, chatRequest)
	if err != nil {
		logEntry.WithError(err).Error("Error during GetChat rpc call")
		errs := response.ClientResponse(http.StatusBadRequest, "Failed to get chat details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	logEntry.Info("Successfully retrieved chat details")
	errs := response.ClientResponse(http.StatusOK, "Successfully retrieved chat details", result, nil)
	c.JSON(http.StatusOK, errs)
}

//Group Chat
func (ch *ChatHandler) GroupMessage(c *gin.Context) {
	userID, _ := c.Get("user_id")
	groupID := c.Param("groupID")
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Websocket Connection Issue", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	userID, ok := c.Get("user_id")
	if !ok {
		errs := response.ClientResponse(http.StatusBadRequest, "User ID in JWT claims is not a string", nil, "")
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	defer func() {
		groupKey := groupID + "_" + strconv.Itoa(userID.(int))
		delete(User, groupKey)
		conn.Close()
	}()

	user := strconv.Itoa(userID.(int))
	groupKey := groupID + "_" + user
	User[groupKey] = conn
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
			c.JSON(http.StatusBadRequest, errs)
			return
		}
		ch.helper.SendMessageToGroup(User, msg, groupID, user)

	}
}

//videocall backend
// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// type client struct {
// 	ChatId primitive.ObjectID
// 	UserId uint
// }

// // type VideoCall struct {
// // 	CallId      string
// // 	CallerId    uint
// // 	RecipientId uint
// // 	Status      string
// // }

// var (
// videoCall  = make(map[string]*VideoCall)
// )

// func (t *ChatHandler) VideoCall(c *gin.Context) {
// 	peerConnectionConfig := webrtc.Configuration{
// 		ICEServers: []webrtc.ICEServer{
// 			{
// 				URLs: []string{"stun:stun.l.google.com:19302"},
// 			},
// 		},
// 	}
// 	peerConnection, err := webrtc.NewPeerConnection(peerConnectionConfig)
// 	if err != nil {
// 		errRes := response.MakeResponse(http.StatusInternalServerError, "internal server error", nil, err.Error())
// 		c.JSON(http.StatusInternalServerError, errRes)
// 		return
// 	}
// 	peerConnection.OnICEConnectionStateChange(func(is webrtc.ICEConnectionState) { fmt.Printf("connection state has changed %s /n", is.String()) })

// 	offer := webrtc.SessionDescription{}

// 	peerConnection.SetRemoteDescription(offer)

// 	answer, err := peerConnection.CreateAnswer(nil)
// 	if err != nil {
// 		errRes := response.MakeResponse(http.StatusInternalServerError, "internal server error", nil, err.Error())
// 		c.JSON(http.StatusInternalServerError, errRes)
// 		return
// 	}
// 	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

// 	peerConnection.SetLocalDescription(answer)
// 	<-gatherComplete
// }
