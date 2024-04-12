package interfaces

type NewChatClient interface {
	CreateChatRoom(user1, user2 int) error
}
