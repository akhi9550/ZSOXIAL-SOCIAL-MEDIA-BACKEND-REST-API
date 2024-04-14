package service

import (
	"context"

	pb "github.com/akhi9550/chat-svc/pkg/pb/chat"
	interfaces "github.com/akhi9550/chat-svc/pkg/usecase/interface"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatServer struct {
	chatUseCase interfaces.ChatUseCase
	pb.UnimplementedChatServiceServer
}

func NewChatServer(UseCaseChat interfaces.ChatUseCase) pb.ChatServiceServer {
	return &ChatServer{
		chatUseCase: UseCaseChat,
	}
}

func (c *ChatServer) CreateChatRoom(ctx context.Context, req *pb.CreateChatRoomRequest) (*pb.CreateChatRoomResponse, error) {
	err := c.chatUseCase.CreateChatRoom(req.Userid, req.Followingid)
	if err != nil {
		return &pb.CreateChatRoomResponse{}, err
	}
	return &pb.CreateChatRoomResponse{}, nil
}

func (c *ChatServer) SaveMessage(ctx context.Context, req *pb.SaveMessageRequest) (*pb.SaveMessageResponse, error) {
	chatID := req.Chatid
	senderID := req.Senderid
	messageContent := req.Message
	chatObjectID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return nil, err
	}
	savedMessageID, err := c.chatUseCase.SaveMessage(chatObjectID, uint(senderID), messageContent)
	if err != nil {
		return nil, err
	}
	response := &pb.SaveMessageResponse{
		Id: savedMessageID.Hex(),
	}
	return response, nil
}

func (c *ChatServer) FetchRecipient(ctx context.Context, req *pb.FetchRecipientRequest) (*pb.FetchRecipientResponse, error) {
	chatIdToHex, err := primitive.ObjectIDFromHex(req.Chatid)
	if err != nil {
		return &pb.FetchRecipientResponse{}, err
	}

	res, err := c.chatUseCase.FetchRecipient(chatIdToHex, uint(req.Userid))

	if err != nil {
		return &pb.FetchRecipientResponse{}, err
	}

	return &pb.FetchRecipientResponse{
		Id: int64(res),
	}, nil

}
func (c *ChatServer) GetAllChats(ctx context.Context, req *pb.GetAllChatsRequest) (*pb.GetAllChatsResponse, error) {
	userID := req.Userid
	data, err := c.chatUseCase.GetAllChats(uint(userID))
	if err != nil {
		return nil, err
	}

	var pbChatResponses []*pb.ChatResponse

	for _, chatResponse := range data {
		chat := chatResponse.Chat
		user := chatResponse.User

		var users []uint64
		for _, u := range chat.Users {
			users = append(users, uint64(u))
		}

		pbChat := &pb.Chat{
			Id:              chat.ID.Hex(),
			Users:           users,
			Lastmessage:     chat.LastMessage,
			Lastmessagetime: &timestamp.Timestamp{Seconds: chat.LastMessageTime.Unix(), Nanos: int32(chat.LastMessageTime.Nanosecond())},
		}

		pbUser := &pb.UserData{
			Userid:   int64(user.UserId),
			Username: user.Username,
			Profile:  user.Profile,
		}

		pbChatResponse := &pb.ChatResponse{
			Chat: pbChat,
			User: pbUser,
		}

		pbChatResponses = append(pbChatResponses, pbChatResponse)
	}

	response := &pb.GetAllChatsResponse{
		Response: pbChatResponses,
	}

	return response, nil
}

func (c *ChatServer) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	chatIDString := req.Chatid
	chatID, err := primitive.ObjectIDFromHex(chatIDString)
	if err != nil {
		return nil, err
	}
	data, err := c.chatUseCase.GetMessages(chatID)
	if err != nil {
		return nil, err
	}

	var response []*pb.Response

	for _, message := range data {
		pbMessage := &pb.Response{
			Id:             message.ID.Hex(),
			SenderId:       uint32(message.SenderID),
			ChatId:         message.ChatID.Hex(),
			Seen:           message.Seen,
			Image:          message.Image,
			MessageContent: message.MessageContent,
			Timestamp:      timestamppb.New(message.Timestamp),
		}
		response = append(response, pbMessage)
	}
	return &pb.GetMessagesResponse{
		Response: response,
	}, nil
}

