run:
	go run cmd/main.go

proto:
	protoc --go_out=. --go-grpc_out=. ./pkg/pb/notification/notification.proto
	protoc --go_out=. --go-grpc_out=. ./pkg/pb/post/post.proto

