package repository

import (
	"context"
	"strconv"

	interfaces "github.com/akhi9550/chat-svc/pkg/repository/interface"
	"github.com/akhi9550/chat-svc/pkg/utils/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository struct {
	MessageCollection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) interfaces.ChatRepository {
	return &ChatRepository{MessageCollection: db.Collection("messages")}
}

func (c *ChatRepository) StoreFriendsChat(message models.MessageReq) error {
	_, err := c.MessageCollection.InsertOne(context.TODO(), message)
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatRepository) GetFriendChat(userID, friendID string, pagination models.Pagination) ([]models.Message, error) {
	var messages []models.Message
	filter := bson.M{"senderid": bson.M{"$in": bson.A{userID, friendID}}, "recipientid": bson.M{"$in": bson.A{friendID, userID}}}
	limit, _ := strconv.Atoi(pagination.Limit)
	offset, _ := strconv.Atoi(pagination.OffSet)

	option := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := c.MessageCollection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{"timestamp", -1}}), option)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func (c *ChatRepository) UpdateReadAsMessage(userID, friendID string) error {
	_, err := c.MessageCollection.UpdateMany(context.TODO(), bson.M{"senderid": bson.M{"$in": bson.A{friendID}}, "recipientid": bson.M{"$in": bson.A{userID}}}, bson.D{{"$set", bson.D{{"status", "send"}}}})
	if err != nil {
		return err
	}
	return nil
}

func (c *ChatRepository) GetGroupMessages(groupID string, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	filter := bson.M{"groupid": groupID}

	option := options.Find().SetLimit(int64(limit)).SetSkip(int64(offset))
	cursor, err := c.MessageCollection.Find(context.TODO(), filter, options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}), option)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}