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

type Message struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Chat Chat
	User UserData
}

type Messages struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SenderID       uint               `json:"sender_id" bson:"sender_id"`
	ChatID         primitive.ObjectID `json:"chat_id" bson:"chat_id"`
	Seen           bool               `json:"seen" bson:"seen"`
	Image          string             `json:"image" bson:"image"`
	MessageContent string             `json:"message_content" bson:"message_content"`
	Timestamp      time.Time          `json:"timestamp" bson:"timestamp"`
}
