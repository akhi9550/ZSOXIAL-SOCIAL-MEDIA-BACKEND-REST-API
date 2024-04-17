package repository

import (
	"context"
	"fmt"
	"strconv"

	interfaces "github.com/akhi9550/chat-svc/pkg/repository/interface"
	"github.com/akhi9550/chat-svc/pkg/utils/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository struct {
	ChatCollection    *mongo.Collection
	MessageCollection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) interfaces.ChatRepository {
	return &ChatRepository{ChatCollection: db.Collection("chats"), MessageCollection: db.Collection("messages")}
}

func (c *ChatRepository) StoreFriendsChat(message models.MessageReq) error {
	fmt.Println("==repo", message)
	_, err := c.MessageCollection.InsertOne(context.TODO(), message)
	if err != nil {
		fmt.Println("data❤️", err)
		return err
	}
	return nil
}

func (c *ChatRepository) GetLastMessage(userID, friendID string) (*models.Message, error) {
	var res = models.Message{}
	filter := bson.M{"senderid": bson.M{"$in": bson.A{userID, friendID}}, "recipientid": bson.M{"$in": bson.A{friendID, userID}}}
	option := options.FindOne().SetSort(bson.D{{"timestamp", -1}})

	err := c.MessageCollection.FindOne(context.TODO(), filter, option).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *ChatRepository) GetMessageCount(userID, friendID string) (int, error) {

	filter := bson.M{"senderid": friendID, "recipientid": userID, "status": "pending"}
	count, err := c.MessageCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		fmt.Println("-", err)
		return 0, err
	}
	fmt.Println("count", count)
	return int(count), nil
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

// func (c *ChatRepository) CreateChatRoom(user1, user2 uint) error {
// 	newChat := domain.Chats{
// 		Users:     []uint{user1, user2},
// 		CreatedAt: time.Now(),
// 	}
// 	_, err := c.ChatCollection.InsertOne(context.TODO(), &newChat)
// 	return err
// }

// func (c *ChatRepository) IsChatExist(user1, user2 uint) (bool, error) {
// 	filter := bson.M{
// 		"users": bson.M{"$all": []uint{user1, user2}},
// 	}

// 	var chat domain.Chats
// 	err := c.ChatCollection.FindOne(context.TODO(), filter).Decode(&chat)

// 	if err == mongo.ErrNoDocuments {

// 		return false, nil
// 	} else if err != nil {

// 		return false, err
// 	}

// 	return true, nil
// }

// func (c *ChatRepository) IsValidChatId(chatId primitive.ObjectID) (bool, error) {
// 	filter := bson.M{
// 		"_id": chatId,
// 	}

// 	var chat domain.Chats
// 	err := c.ChatCollection.FindOne(context.TODO(), filter).Decode(&chat)

// 	if err == mongo.ErrNoDocuments {
// 		return false, nil
// 	} else if err != nil {
// 		return false, err
// 	}

// 	return true, nil

// }

// func (c *ChatRepository) GetAllChats(id uint) ([]models.Chat, error) {
// 	filter := bson.M{"users": bson.M{"$in": []uint{id}}}
// 	projection := bson.M{"_id": 1, "users": 1, "last_message": 1, "last_message_time": 1}

// 	cursor, err := c.ChatCollection.Find(context.TODO(), filter, options.Find().SetProjection(projection))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.TODO())

// 	var chats []models.Chat
// 	if err := cursor.All(context.TODO(), &chats); err != nil {
// 		return nil, err
// 	}
// 	return chats, nil
// }

// func (c *ChatRepository) GetMessages(id primitive.ObjectID) ([]domain.Messages, error) {
// 	ctx := context.TODO()
// 	filter := bson.M{"chat_id": id}
// 	cursor, err := c.MessageCollection.Find(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(ctx)
// 	var messages []domain.Messages
// 	err = cursor.All(ctx, &messages)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return messages, nil

// }

// func (c *ChatRepository) SaveMessage(message domain.Messages) (primitive.ObjectID, error) {
// 	id, err := c.MessageCollection.InsertOne(context.TODO(), message)
// 	if err != nil {
// 		return primitive.ObjectID{}, err
// 	}
// 	return id.InsertedID.(primitive.ObjectID), nil
// }

// func (c *ChatRepository) UpdateLastMessageAndTime(chatId primitive.ObjectID, lastMessage string, time time.Time) error {
// 	filter := bson.M{"_id": chatId}
// 	update := bson.M{"$set": bson.M{"last_message": lastMessage, "last_message_time": time}}
// 	_, err := c.ChatCollection.UpdateOne(context.TODO(), filter, update)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (c *ChatRepository) ReadMessage(chatId primitive.ObjectID, senderId uint) (int64, error) {

// 	filter := bson.M{"chat_id": chatId, "sender_id": senderId, "seen": false}

// 	update := bson.M{"$set": bson.M{"seen": true}}

// 	res, err := c.MessageCollection.UpdateMany(context.TODO(), filter, update)

// 	if err != nil {

// 		return 0, err
// 	}

// 	return res.UpsertedCount, nil

// }

// func (c *ChatRepository) FetchRecipient(chatId primitive.ObjectID, userId uint) (uint, error) {
// 	filter := bson.M{"_id": chatId}
// 	projection := bson.M{"_id": 0, "users": bson.M{"$elemMatch": bson.M{"$ne": userId}}, "created_at": 0}

// 	chat := domain.Chats{}
// 	c.ChatCollection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&chat)

// 	return chat.Users[0], nil

// }

// func (c *ChatRepository) DeleteChatsAndMessagesByUserID(userID uint) error {
// 	var chatIDs []primitive.ObjectID
// 	filter := bson.M{"users": userID}
// 	cursor, err := c.ChatCollection.Find(context.Background(), filter)
// 	if err != nil {
// 		return err
// 	}
// 	defer cursor.Close(context.Background())

// 	for cursor.Next(context.Background()) {
// 		var chat domain.Chats
// 		if err := cursor.Decode(&chat); err != nil {
// 			return err
// 		}
// 		chatIDs = append(chatIDs, chat.ID)
// 	}

// 	_, err = c.MessageCollection.DeleteMany(context.Background(), bson.M{"chat_id": bson.M{"$in": chatIDs}})
// 	if err != nil {
// 		return err
// 	}

// 	_, err = c.ChatCollection.DeleteMany(context.Background(), filter)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
