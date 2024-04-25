package handler

import (
	"context"
	"net/http"
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
	var notificationRequest models.NotificationPagination
	if err := c.ShouldBindJSON(&notificationRequest); err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Details not in correct format", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		errs := response.ClientResponse(http.StatusBadRequest, "User ID not found in JWT claims", nil, "")
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	UserID, _ := userID.(int)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	result, err := n.GRPC_Client.GetNotification(ctx, &pb.GetNotificationRequest{
		UserID: int64(UserID),
		Offset: int64(notificationRequest.Offset),
		Limit:  int64(notificationRequest.Limit),
	})
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Failed to get notification details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	errs := response.ClientResponse(http.StatusOK, "Successfully retrieved Notifications", result, nil)
	c.JSON(http.StatusOK, errs)
}
