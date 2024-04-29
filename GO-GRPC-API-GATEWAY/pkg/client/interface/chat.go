package interfaces

import "github.com/akhi9550/api-gateway/pkg/utils/models"

type ChatClient interface {
	GetChat(userID string, req models.ChatRequest) ([]models.Message, error)
}
