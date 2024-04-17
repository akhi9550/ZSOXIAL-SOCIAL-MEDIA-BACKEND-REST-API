package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Users           []uint             `json:"users" bson:"users"`
	LastMessage     string             `json:"last_message" bson:"last_message"`
	LastMessageTime time.Time          `json:"last_message_time" bson:"last_message_time"`
}

type Messages struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Chat Chat
	User UserData
}

type Message struct {
	ID          string    `bson:"_id"`
	SenderID    string    `bson:"senderid"`
	RecipientID string    `bson:"recipientid"`
	Content     string    `bson:"content"`
	Timestamp   time.Time `bson:"timestamp"`
}

type MessageReq struct {
	SenderID    string    `bson:"senderid"`
	RecipientID string    `bson:"recipientid"`
	Content     string    `bson:"content"`
	Timestamp   time.Time `bson:"timestamp"`
}

type Pagination struct {
	Limit  string
	OffSet string
}
