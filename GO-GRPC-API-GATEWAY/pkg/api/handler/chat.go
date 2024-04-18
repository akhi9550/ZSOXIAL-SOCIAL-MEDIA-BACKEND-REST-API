package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
}

func NewChatHandler(chatClient pb.ChatServiceClient, helper *helper.Helper) *ChatHandler {
	return &ChatHandler{
		GRPC_Client: chatClient,
		helper:      helper,
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


// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
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
//
//	connection = make(map[*websocket.Conn]*client)
//	user       = make(map[uint]*websocket.Conn)
//	// videoCall  = make(map[string]*VideoCall)
//
// )

// func (ch *ChatHandler) Chat(c *gin.Context) {
// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

// 	if err != nil {
// 		errorRes := response.ClientResponse(http.StatusBadRequest, "data is not in required format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errorRes)
// 		return
// 	}

// 	userid, ok := c.Get("user_id")

// 	if !ok || userid == nil {
// 		c.JSON(401, gin.H{"error": "userId not found in context or is nil"})
// 		return
// 	}

// 	userIdInt64 := userid.(int64)
// 	chatID := c.Query("chatId")
// 	if chatID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "chatID parameter is missing"})
// 		return
// 	}

// 	objectChatId, err := primitive.ObjectIDFromHex(chatID)
// 	if err != nil {
// 		fmt.Println("Error converting chatID to ObjectID:", err)
// 		return
// 	}

// 	connection[conn] = &client{ChatId: objectChatId, UserId: uint(userIdInt64)}
// 	user[uint(userIdInt64)] = conn
// 	go func() {

// 		for {
// 			_, msg, err := conn.ReadMessage()
// 			if err != nil {
// 				break
// 			}
// 			userId := connection[conn].UserId
// 			chatID := connection[conn].ChatId

// 			_, err = ch.GRPC_Client.SaveMessage(chatID, int(userId), string(msg))
// 			if err != nil {
// 				fmt.Println("Error saving message:", err)
// 				break
// 			}

// 			conn.WriteMessage(websocket.TextMessage, msg)

// 			recipient, err := ch.GRPC_Client.FetchRecipient(chatID, int(userId))
// 			if err != nil {
// 				fmt.Println("Error fetching recipient:", err)
// 				break
// 			}
// 			recipientID := recipient

// 			if value, ok := user[recipientID]; ok {
// 				err = value.WriteMessage(websocket.TextMessage, msg)
// 				if err != nil {
// 					delete(connection, value)
// 					delete(user, recipientID)
// 				}
// 			}
// 		}
// 	}()
// }

// func (ch *ChatHandler) GetAllChats(c *gin.Context) {
// 	userIDInterface, _ := c.Get("user_id")
// 	userID, ok := userIDInterface.(int)
// 	if !ok {
// 		return
// 	}
// 	data, err := ch.GRPC_Client.GetAllChats(uint(userID))
// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}

// 	success := response.ClientResponse(http.StatusOK, "Successfully Get Chats", data, nil)
// 	c.JSON(http.StatusOK, success)
// }

// func (ch *ChatHandler) GetMessages(c *gin.Context) {
// 	chatID, err := primitive.ObjectIDFromHex(c.Param("chatId"))
// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "chatID not in right format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}
// 	data, err := ch.GRPC_Client.GetMessages(chatID)
// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}
// 	success := response.ClientResponse(http.StatusOK, "Successfully showing messages in the given chatID", data, nil)
// 	c.JSON(http.StatusOK, success)
// }

// ///////////

//////////////////////////////
// func (ch *ChatHandler) MakeMessageRead(c *gin.Context) {
// 	userIdInterface, exists := c.Get("user_id")
// 	if !exists {
// 		errs := response.ClientResponse(http.StatusBadRequest, "User ID not found in context", nil, "")
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}

// 	userId, ok := userIdInterface.(int)
// 	if !ok {
// 		errs := response.ClientResponse(http.StatusBadRequest, "User ID in context is not of type int", nil, "")
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}

// 	chatId, err := primitive.ObjectIDFromHex(c.Param("chatId"))
// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "chatID not in right format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}

// 	data, err := ch.GRPC_Client.ReadMessage(userId, chatId)
// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}
// 	success := response.ClientResponse(http.StatusOK, "Successfully make these messages to read", data, nil)
// 	c.JSON(http.StatusOK, success)
// }

// func (t *ChatHandler) ChatPage(c *gin.Context) {
// 	chatId := c.Query("chatId")
// 	if chatId == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "chatId is required"})
// 		return
// 	}
// 	// Render the HTML template with the chatId
// 	c.HTML(http.StatusOK, "index.html", gin.H{
// 		"chatId": chatId,
// 	})
// }

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
