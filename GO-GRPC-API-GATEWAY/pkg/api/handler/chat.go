package handler

import (
	"fmt"
	"net/http"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/utils/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type client struct {
	ChatId primitive.ObjectID
	UserId uint
}

type VideoCall struct {
	CallId      string
	CallerId    uint
	RecipientId uint
	Status      string
}

var (
	connection = make(map[*websocket.Conn]*client)
	user       = make(map[uint]*websocket.Conn)
	// videoCall  = make(map[string]*VideoCall)
)

type ChatHandler struct {
	GRPC_Client interfaces.ChatClient
}

func NewChatHandler(chatClient interfaces.ChatClient) *ChatHandler {
	return &ChatHandler{
		GRPC_Client: chatClient,
	}
}

func (ch *ChatHandler) GetAllChats(c *gin.Context) {
	userIDInterface, _ := c.Get("user_id")
	userID, ok := userIDInterface.(int)
	if !ok {
		return
	}
	data, err := ch.GRPC_Client.GetAllChats(uint(userID))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	success := response.ClientResponse(http.StatusOK, "Successfully Get Chats", data, nil)
	c.JSON(http.StatusOK, success)
}

func (ch *ChatHandler) GetMessages(c *gin.Context) {
	chatID, err := primitive.ObjectIDFromHex(c.Param("chatId"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "chatID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	data, err := ch.GRPC_Client.GetMessages(chatID)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully showing messages in the given chatID", data, nil)
	c.JSON(http.StatusOK, success)
}

func (t *ChatHandler) ChatPage(c *gin.Context) {
	chatId := c.Param("chatId")
	c.HTML(http.StatusOK, "index.html", chatId)
}

func (ch *ChatHandler) MakeMessageRead(c *gin.Context) {
	userId, _ := c.Get("userId")
	chatId, err := primitive.ObjectIDFromHex(c.Param("chatId"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "chatID not in right format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	data, err := ch.GRPC_Client.ReadMessage(userId.(uint), chatId)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	success := response.ClientResponse(http.StatusOK, "Successfully make these messages to read", data, nil)
	c.JSON(http.StatusOK, success)
}

func (ch *ChatHandler) Chat(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	userID, _ := c.Get("userId")
	chatId, err := primitive.ObjectIDFromHex(c.Param("chatId"))
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "chatID not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	connection[conn] = &client{ChatId: chatId, UserId: userID.(uint)}
	user[userID.(uint)] = conn

	go func() {

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				break
			}
			userId := connection[conn].UserId
			chatID := connection[conn].ChatId
			_, err = ch.GRPC_Client.SaveMessage(chatID, userId, string(msg))
			if err != nil {
				fmt.Println("error in saving message")
				break
			}
			conn.WriteMessage(websocket.TextMessage, msg)
			recipient, err := ch.GRPC_Client.FetchRecipient(chatID, userId)
			if err != nil {
				fmt.Println("error in fetching recipient id")
				break
			}
			if value, ok := user[recipient]; ok {
				err = value.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					delete(connection, value)
					delete(user, recipient)
				}
			}
		}
	}()
}

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
