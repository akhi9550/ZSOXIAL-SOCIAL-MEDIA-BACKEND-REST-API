package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type VideoCallHandler struct{}

// type WebRTCHandler struct{}
func NewVideoCallHandler() *VideoCallHandler {
	return &VideoCallHandler{}
}

func (v *VideoCallHandler) RequestToRoom(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (v *VideoCallHandler) ConnectedPage(c *gin.Context) {
	c.HTML(http.StatusOK, "lobby.html", nil)
}

func (v *VideoCallHandler) IndexedPage(c *gin.Context) {
	room := c.DefaultQuery("room", "")
	c.HTML(http.StatusOK, "index.html", gin.H{"room": room})
}

func (v *VideoCallHandler) GetStoredOffer(c *gin.Context) {
	c.String(http.StatusOK, GetStoredOffer())
}

var signalingMsg struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

var StoredOffer string

func GetStoredOffer() string {
	return StoredOffer
}

var answerSDP1 string

func (v *VideoCallHandler) SetupRoutes(router *gin.RouterGroup) {
	router.GET("/ws", v.handleWebSocket)
}

func (v *VideoCallHandler) handleWebSocket(c *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("error upgrade", err)
		return
	}
	defer conn.Close()
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		return
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("error readmsg", err)
			break
		}
		if err := json.Unmarshal(message, &signalingMsg); err != nil {
			fmt.Println("error json:", err)
			continue
		}
		switch signalingMsg.Type {
		case "offer":
			fmt.Println("case offer")
			err := handleOffer(conn, signalingMsg.Data, peerConnection)
			if err != nil {
				fmt.Println("error case1", err.Error())
				break
			}
		case "answer":
			fmt.Println("case answer")
			err := handleAnswer(conn, signalingMsg.Data, peerConnection)
			if err != nil {
				fmt.Println("error case2", err.Error())
				break
			}
		case "candidate":
			fmt.Println("case candidate")
			err := handleICECandidate(conn, signalingMsg.Data, peerConnection)
			if err != nil {
				fmt.Println("error case3", err.Error())
				break
			}
		default:
			fmt.Println("error", "unknown type")
		}
	}

}

func handleOffer(conn *websocket.Conn, offerSDPstring string, peerConnection *webrtc.PeerConnection) error {
	offerSDP := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offerSDPstring,
	}

	StoredOffer = offerSDPstring
	err := peerConnection.SetLocalDescription(offerSDP)
	if err != nil {
		return errors.Join(errors.New("setRemote"), err)
	}
	return nil
}

func handleAnswer(conn *websocket.Conn, answerSDP string, peerConnection *webrtc.PeerConnection) error {

	answer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answerSDP1,
	}

	err := peerConnection.SetLocalDescription(answer)
	if err != nil {
		return errors.Join(errors.New("SetRemoteDsc"), err)
	}

	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  StoredOffer,
	}

	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		return errors.Join(errors.New("setRemote"), err)
	}

	return nil
}

func handleICECandidate(conn *websocket.Conn, candidate string, peerConnection *webrtc.PeerConnection) error {

	iceCandidate := webrtc.ICECandidateInit{
		Candidate: candidate,
	}

	err := peerConnection.AddICECandidate(iceCandidate)
	if err != nil {
		return errors.Join(errors.New("AddICE"), err)
	}

	return nil
}
