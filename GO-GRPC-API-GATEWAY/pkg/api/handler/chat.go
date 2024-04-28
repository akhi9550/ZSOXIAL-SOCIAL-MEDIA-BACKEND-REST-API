package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/akhi9550/api-gateway/pkg/helper"
	pb "github.com/akhi9550/api-gateway/pkg/pb/chat"
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
	GRPC_Client pb.ChatServiceClient
	helper      *helper.Helper
	ChatCachig  *helper.RedisChatCaching
}

func NewChatHandler(chatClient pb.ChatServiceClient, helper *helper.Helper, chatCache *helper.RedisChatCaching) *ChatHandler {
	return &ChatHandler{
		GRPC_Client: chatClient,
		helper:      helper,
		ChatCachig:  chatCache,
	}
}

func (ch *ChatHandler) FriendMessage(c *gin.Context) {
	fmt.Println("message called")
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
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

func (ch *ChatHandler) GetChat(c *gin.Context) {
	var chatRequest models.ChatRequest
	if err := c.ShouldBindJSON(&chatRequest); err != nil {
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
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// defer cancel()

	result, err := ch.ChatCachig.GetChat(userID, chatRequest)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Failed to get chat details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	errs := response.ClientResponse(http.StatusOK, "Successfully retrieved chat details", result, nil)
	c.JSON(http.StatusOK, errs)
}

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
