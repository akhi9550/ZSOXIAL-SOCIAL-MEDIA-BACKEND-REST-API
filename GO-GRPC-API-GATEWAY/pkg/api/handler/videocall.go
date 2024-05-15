package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type VideoCallHandler struct{}

type Room struct {}

func NewVideoCallHandler() *VideoCallHandler {
	return &VideoCallHandler{}
}

func (v *VideoCallHandler) ExitPage(c *gin.Context) {
	c.HTML(http.StatusOK, "exit.html", nil)
}

func (v *VideoCallHandler) ErrorPage(c *gin.Context) {
	c.HTML(http.StatusOK, "error.html", nil)
}

func (v *VideoCallHandler) IndexedPage(c *gin.Context) {
	room := c.DefaultQuery("room", "")
	c.HTML(http.StatusOK, "index.html", gin.H{"room": room})
}

// func (v *VideoCallHandler) SetupRoutes(group *gin.RouterGroup) {
// 	group.GET("/ws", v.handleWebSocket)
// }
// 
// func (v *VideoCallHandler) handleWebSocket(c *gin.Context) {
// 	upgrader := websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		return
// 	}
// 	defer conn.Close()

// 	roomId := c.DefaultQuery("room", "")
// 	v.mu.Lock()
// 	room, exists := v.rooms[roomId]
// 	v.mu.Unlock()
// 	if !exists {
// 		return
// 	}

// 	room.mu.Lock()
// 	room.Participants[conn] = true
// 	room.mu.Unlock()

// 	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
// 	if err != nil {
// 		return
// 	}
// 	defer peerConnection.Close()
// }
