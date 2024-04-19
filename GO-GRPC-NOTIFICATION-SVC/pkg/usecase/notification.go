package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	postclientinterfaces "github.com/akhi9550/notification-svc/pkg/client/interface"
	"github.com/akhi9550/notification-svc/pkg/config"
	"github.com/akhi9550/notification-svc/pkg/domain"
	"github.com/akhi9550/notification-svc/pkg/helper"
	interfaces "github.com/akhi9550/notification-svc/pkg/repository/interface"
	services "github.com/akhi9550/notification-svc/pkg/usecase/interface"
	"github.com/akhi9550/notification-svc/pkg/utils/models"
)

type NotificationUseCase struct {
	notificationRepository interfaces.NotificationRepository
	Postclient             postclientinterfaces.PostClient
}

func NewUserUserUseCase(repository interfaces.NotificationRepository, postclient postclientinterfaces.PostClient) services.NotificationUsecaseInterface {
	return &NotificationUseCase{
		notificationRepository: repository,
		Postclient:             postclient,
	}
}

var doneCh = make(chan struct{})

func (n *NotificationUseCase) AddLikeNotification(notification models.LikeNotification) (int64, error) {
	cfg,_:=config.LoadConfig()
	var userID int64
	// brokerUrl := []string{"localhost:9092"}
	// // Use the same Kafka broker URL as the producer

	// topic := "like_notification"

	worker, err := helper.ConnectToLikeConsumer([]string{cfg.KafkaPort})
	if err != nil {
		panic(err)
	}

	partitionConsumer, err := worker.ConsumePartition(cfg.KafkaLikeTopic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}
	fmt.Println("consumer started")

	signchan := make(chan os.Signal, 1)
	signal.Notify(signchan, syscall.SIGINT, syscall.SIGTERM)
	messagecount := 0

	go func() {
		for {
			select {
			case err := <-partitionConsumer.Errors():
				fmt.Println("Error:", err)
			case msg := <-partitionConsumer.Messages():

				var message models.LikeNotification

				if err := json.Unmarshal(msg.Value, &message); err != nil {
					fmt.Println("Error parsing message:", err)
					//continue
				}

				// Add the notification to the repository
				userID, err := n.notificationRepository.AddNotification(domain.Notification{
					UserID:  message.UserID,
					PostID:  message.PostID,
					Message: message.Message,
				})
				if err != nil {
					fmt.Println("Error adding notification to repository:", err)
					//continue
				}

				fmt.Printf("Notification added for UserID: %d\n", userID)

				messagecount++
				fmt.Printf("received message count: %d | Topic: %s | Message: %s\n", messagecount, string(msg.Topic), string(msg.Value))
			case <-signchan:
				fmt.Println("interruption detected")
				doneCh <- struct{}{}
			}
		}
	}()

	// Wait for the consumer to finish processing messages
	<-doneCh

	fmt.Println("processed", messagecount, "messages")
	if err := worker.Close(); err != nil {
		panic(err)
	}
	return userID, nil
}

//  consume kafka messages for like notifications

// func (u *NotifyUseCase) ConsumeMessage(user_id int64) (models.LikeNotification, error) {
// 	// Initialize a channel to pass extracted values back
// 	ch := make(chan models.LikeNotification)

// 	// Define Kafka consumer parameters
// 	brokerURL := []string{"localhost:9092"}
// 	topic := "like_notification"

// 	// Connect to Kafka broker
// 	worker, err := helper.ConnectToLikeConsumer(brokerURL)
// 	if err != nil {
// 		return models.LikeNotification{}, err
// 	}
// 	defer worker.Close() // Close Kafka worker at the end

// 	// Consume Kafka messages from the specified topic
// 	partitionConsumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
// 	if err != nil {
// 		return models.LikeNotification{}, err
// 	}
// 	defer partitionConsumer.Close() // Close Kafka partition consumer at the end

// 	// Goroutine to process Kafka messages
// 	go func() {
// 		defer close(ch) // Close the channel when done processing messages

// 		for {
// 			select {
// 			case err := <-partitionConsumer.Errors():
// 				fmt.Println("Error:", err)
// 			case msg := <-partitionConsumer.Messages():
// 				var likeNotification models.LikeNotification

// 				// Unmarshal Kafka message into LikeNotification struct
// 				if err := json.Unmarshal(msg.Value, &likeNotification); err != nil {
// 					fmt.Println("Error parsing message:", err)
// 					continue
// 				}
// 				// // Call external service to get user ID
// 				// userID, err := u.Postclient.GetUserId(likeNotification.PostID)
// 				// if err != nil {
// 				// 	fmt.Println("Error getting user ID:", err)
// 				// 	continue
// 				// }

// 				// // Check if the user IDs match
// 				// if userID != likeNotification.UserID {
// 				// 	fmt.Println("User IDs don't match, skipping message")
// 				// 	continue
// 				// }

// 				// Send extracted LikeNotification through the channel
// 				ch <- likeNotification
// 			}
// 		}
// 	}()

// 	// Wait for Kafka messages to be processed and return the extracted LikeNotification
// 	likeNotification := <-ch

// 	return likeNotification, nil
// }

func (n *NotificationUseCase) ConsumeCommentMessage(user_id int64) (models.CommentNotification, error) {
	// Initialize a channel to pass extracted values back
	cfg,_:=config.LoadConfig()
	ch := make(chan models.CommentNotification, 1) // Buffered channel to avoid goroutine leak

	// Define Kafka consumer parameters
	// brokerURL := []string{"localhost:9092"}
	// topic := "comment_notifications"

	// Connect to Kafka broker
	worker, err := helper.ConnectToCommentConsumer([]string{cfg.KafkaPort})

	if err != nil {
		return models.CommentNotification{}, err
	}

	// Consume Kafka messages from the specified topic
	partitionConsumer, err := worker.ConsumePartition(cfg.KafkaCommentTopic, 0, sarama.OffsetOldest)
	if err != nil {
		worker.Close() // Close Kafka worker in case of error
		return models.CommentNotification{}, err
	}

	// Goroutine to process Kafka messages
	go func() {
		defer close(ch)                 
		defer partitionConsumer.Close()

		for {
			select {
			case err := <-partitionConsumer.Errors():
				fmt.Println("Error:", err)
			case msg := <-partitionConsumer.Messages():
				var commentNotification models.CommentNotification

				// Unmarshal Kafka message into CommentNotification struct
				if err := json.Unmarshal(msg.Value, &commentNotification); err != nil {
					fmt.Println("Error parsing message:", err)
					continue
				}

				fmt.Println("conetnt is ", commentNotification.Content)

				// Send extracted CommentNotification through the channel
				ch <- commentNotification
			}
		}
	}()

	// Wait for Kafka messages to be processed and return the extracted CommentNotification
	select {
	case commentsContent, ok := <-ch:
		if !ok {
			return models.CommentNotification{}, errors.New("channel closed before receiving data")
		}
		return commentsContent, nil
	case <-time.After(time.Second * 5): // Timeout after 5 seconds
		return models.CommentNotification{}, errors.New("timeout waiting for comment notification")
	}
}

func (u *NotificationUseCase) ConsumeMessage(user_id int64) (models.LikeNotification, error) {
	// Initialize a channel to pass extracted values back
	cfg,_:=config.LoadConfig()
	ch := make(chan models.LikeNotification, 1)

	// Define Kafka consumer parameters
	// brokerURL := []string{"localhost:9092"}
	// topic := "like_notification"

	// Connect to Kafka broker
	worker, err := helper.ConnectToLikeConsumer([]string{cfg.KafkaPort})
	if err != nil {
		return models.LikeNotification{}, err
	}

	// Consume Kafka messages from the specified topic
	partitionConsumer, err := worker.ConsumePartition(cfg.KafkaLikeTopic, 0, sarama.OffsetOldest)
	if err != nil {
		defer worker.Close() // Close Kafka worker at the end

		return models.LikeNotification{}, err

	}

	// Goroutine to process Kafka messages
	go func() {
		defer close(ch)                 // Close the channel when done processing messages
		defer partitionConsumer.Close() // Close Kafka partition consumer at the end

		for {
			select {
			case err := <-partitionConsumer.Errors():
				fmt.Println("Error:", err)
			case msg := <-partitionConsumer.Messages():
				var likeNotification models.LikeNotification

				// Unmarshal Kafka message into LikeNotification struct
				if err := json.Unmarshal(msg.Value, &likeNotification); err != nil {
					fmt.Println("Error parsing message:", err)
					continue
				}
				fmt.Println("content is ", likeNotification.Content)
				ch <- likeNotification
			}
		}
	}()

	select {
	case likesContent, ok := <-ch:
		if !ok {
			return models.LikeNotification{}, errors.New("channel closed before receiving data")
		}
		return likesContent, nil
	case <-time.After(time.Second * 5): // Timeout after 5 seconds
		return models.LikeNotification{}, errors.New("timeout waiting for comment notification")
	}
}
