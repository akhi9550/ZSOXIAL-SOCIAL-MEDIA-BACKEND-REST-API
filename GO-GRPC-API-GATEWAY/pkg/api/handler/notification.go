package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	pb "github.com/akhi9550/api-gateway/pkg/pb/notification"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/akhi9550/api-gateway/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	GRPC_Client pb.NotificationServiceClient
}

func NewNotificationHandler(notificationClient pb.NotificationServiceClient) *NotificationHandler {
	return &NotificationHandler{
		GRPC_Client: notificationClient,
	}
}

func (n *NotificationHandler) GetNotification(c *gin.Context) {
	var notificationRequest models.NotificationReq
	if err := c.ShouldBindJSON(&notificationRequest); err != nil {
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
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := n.GRPC_Client.GetNotification(ctx, &pb.GetNotificationRequest{
		UserID: userID,
		Offset: notificationRequest.Offset,
		Limit:  notificationRequest.Limit,
	})
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Failed to get notification details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	errs := response.ClientResponse(http.StatusOK, "Successfully retrieved Notifications", result, nil)
	c.JSON(http.StatusOK, errs)
}


// func (n *NotificationHandler) SendCommentedNotification(c *gin.Context) {

// }

// func (n *NotificationHandler) SendLikeNotification(c *gin.Context) {
// 	userID, ok := c.Get("user_id")
// 	if !ok || userID == nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, "")
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}

// 	PostId := c.Query("postId")
// 	postID, err := strconv.Atoi(PostId)
// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}
// 	data, err := n.GRPC_Client.SendLikeNotification(userID.(int), postID)

// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "Error in connecting with  notification service", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}
// 	success := response.ClientResponse(http.StatusOK, "Successfully Get All Notifications", data, nil)
// 	c.JSON(http.StatusOK, success)

// }

// func ConsumeKafkaMessages(c *gin.Context, p pb.NotificationServiceClient) {
// 	// Call the gRPC method to start streaming Kafka messages
// 	stream, err := p.ConsumeKafkaMessages(c, &pb.Empty{})
// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "Error while receiving Kafka message from stream", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}

// 	// Stream Kafka messages to the client
// 	c.Stream(func(w io.Writer) bool {
// 		// Receive next Kafka message from the stream
// 		message, err := stream.Recv()
// 		if err == io.EOF {
// 			return false
// 		}
// 		if err != nil {
// 			errs := response.ClientResponse(http.StatusBadRequest, "Error while receiving Kafka message from stream", nil, err.Error())
// 			c.JSON(http.StatusBadRequest, errs)
// 			return false
// 		}

// 		// Process the Kafka message as needed (e.g., log, store in database, etc.)
// 		fmt.Printf("Received Kafka message: %+v\n", message)

// 		// Send the Kafka message to the client
// 		c.SSEvent("message", message)

// 		// Continue streaming
// 		return true
// 	})
// 	stream.CloseSend()
// }
// func (n *NotificationHandler) ConsumeKafkaCommentMessages(c *gin.Context) {
// 	userID, ok := c.Get("user_id")
// 	if !ok || userID == nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "Invalid userID", nil, "")
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}

// 	data, err := n.GRPC_Client.ConsumeKafkaCommentMessages(userID.(int))

// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "Error connecting notification service", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}
// 	success := response.ClientResponse(http.StatusOK, "Successfully New message", data, nil)
// 	c.JSON(http.StatusOK, success)

// }
// func (n *NotificationHandler) ConsumeKafkaLikeMessages(c *gin.Context) {
// 	userID, ok := c.Get("user_id")
// 	if !ok || userID == nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "Invalid userID", nil, "")
// 		c.JSON(http.StatusBadRequest, errs)
// 		return
// 	}
// 	data, err := n.GRPC_Client.ConsumeKafkaLikeMessages(userID.(int))
// 	if err != nil {
// 		errs := response.ClientResponse(http.StatusBadRequest, "Error connecting notification service", nil, err.Error())
// 		c.JSON(http.StatusBadRequest, errs)
// 		return

// 	}
// 	success := response.ClientResponse(http.StatusOK, "Successfully New like", data, nil)
// 	c.JSON(http.StatusOK, success)

// }

// ////////////
