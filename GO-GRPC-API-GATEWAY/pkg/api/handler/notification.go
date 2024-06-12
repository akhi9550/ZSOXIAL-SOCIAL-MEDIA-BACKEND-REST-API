package handler

import (
	"net/http"

	interfaces "github.com/akhi9550/api-gateway/pkg/client/interface"
	"github.com/akhi9550/api-gateway/pkg/logging"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
	"github.com/akhi9550/api-gateway/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	GRPC_Client interfaces.NotificationClient
}

func NewNotificationHandler(notificationClient interfaces.NotificationClient) *NotificationHandler {
	return &NotificationHandler{
		GRPC_Client: notificationClient,
	}
}

// @Summary			Show All Notifications
// @Description		Retrieve  User All Notifications
// @Tags			Notifications
// @Accept			json
// @Produce		    json
// @Security		Bearer
// @Param			notificationRequest  	body		models.NotificationPagination	true	"Notification details"
// @Success		200		{object}	response.Response{}
// @Failure		500		{object}	response.Response{}
// @Router			/notification   [GET]
func (n *NotificationHandler) GetNotification(c *gin.Context) {
	logEntry := logging.GetLogger().WithField("context", "GetNotification")
	logEntry.Info("Processing GetNotification")
	var notificationRequest models.NotificationPagination
	if err := c.ShouldBindJSON(&notificationRequest); err != nil {
		logEntry.WithError(err).Error("Error binding request body")
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

	result, err := n.GRPC_Client.GetNotification(UserID, notificationRequest)
	if err != nil {
		logEntry.WithError(err).Error("Error during GetNotification rpc call")
		errs := response.ClientResponse(http.StatusBadRequest, "Failed to get notification details", nil, err.Error())
		c.JSON(http.StatusBadRequest, errs)
		return
	}
	logEntry.Info("Successfully retrieved Notifications")
	errs := response.ClientResponse(http.StatusOK, "Successfully retrieved Notifications", result, nil)
	c.JSON(http.StatusOK, errs)
}
