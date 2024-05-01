run:
	go run cmd/main.go

swag: 
	swag init -g cmd/main.go -o ./cmd/docs

proto:
	protoc --go_out=. --go-grpc_out=. ./pkg/pb/auth/auth.proto
	protoc --go_out=. --go-grpc_out=. ./pkg/pb/post/post.proto
	protoc --go_out=. --go-grpc_out=. ./pkg/pb/chat/chat.proto
	protoc --go_out=. --go-grpc_out=. ./pkg/pb/notification/notification.proto