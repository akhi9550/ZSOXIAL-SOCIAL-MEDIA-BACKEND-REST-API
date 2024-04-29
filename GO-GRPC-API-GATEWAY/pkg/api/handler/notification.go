package handler

import (
	"net/http"

	"github.com/akhi9550/api-gateway/pkg/helper"
	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/akhi9550/api-gateway/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	GRPC_Client       interfaces.NotificationClient
	NotificationCachig *helper.RedisNotificationCaching
}

func NewNotificationHandler(notificationClient interfaces.NotificationClient, notificationCache *helper.RedisNotificationCaching) *NotificationHandler {
	return &NotificationHandler{
		GRPC_Client: notificationClient,
		NotificationCachig: notificationCache,
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
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	// defer cancel()

	result, err := n.NotificationCachig.GetNotification(UserID, notificationRequest)
	if err != nil {
		errs := response.ClientResponse(http.StatusBadRequest, "Failed to get notification details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}

	errs := response.ClientResponse(http.StatusOK, "Successfully retrieved Notifications", result, nil)
	c.JSON(http.StatusOK, errs)
}
