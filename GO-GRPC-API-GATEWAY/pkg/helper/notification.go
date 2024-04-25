package helper

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/akhi9550/api-gateway/pkg/config"
	"github.com/akhi9550/api-gateway/pkg/utils/models"
)

func SendLikeNotification(message string, data models.LikePostResponse, postID int) {
	notification := models.Notification{
		LikedUser:   data.LikedUser,
		UserProfile: data.Profile,
		PostID:      postID,
		Content:     message,
		Timestamp:   data.CreatedAt,
	}

	fmt.Println("message:", message)
	err := KafkaLikeProducer(notification)
	if err != nil {
		fmt.Println("error sending notification to Kafka:", err)
		return
	}

	fmt.Println("==sent successfully to Kafka==")
}

func KafkaLikeProducer(message models.Notification) error {
	cfg, _ := config.LoadConfig()
	configs := sarama.NewConfig()
	configs.Producer.Return.Successes = true
	configs.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer([]string{cfg.KafkaPort}, configs)
	if err != nil {
		return err
	}

	result, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{Topic: cfg.KafkaLikeTopic, Key: sarama.StringEncoder("Notification"), Value: sarama.StringEncoder(result)}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("err send message in kafka ", err)
	}
	log.Printf("[producer] partition id: %d; offset:%d, value: %v\n", partition, offset, msg)
	return nil
}

func SendCommentNotification(message string, data models.LikePostResponse, postID int) {
	notification := models.Notification{
		LikedUser:   data.LikedUser,
		UserProfile: data.Profile,
		PostID:      postID,
		Content:     message,
		Timestamp:   data.CreatedAt,
	}

	fmt.Println("message:", message)

	err := KafkaLikeProducer(notification)
	if err != nil {
		fmt.Println("error sending notification to Kafka:", err)
		return
	}

	fmt.Println("==sent successfully to Kafka==")
}

func KafkaCommentProducer(message models.Message) error {
	fmt.Println("from kafka ", message)
	cfg, _ := config.LoadConfig()
	configs := sarama.NewConfig()
	configs.Producer.Return.Successes = true
	configs.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer([]string{cfg.KafkaPort}, configs)
	if err != nil {
		return err
	}

	result, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{Topic: cfg.KafkaCommentTopic, Key: sarama.StringEncoder("message"), Value: sarama.StringEncoder(result)}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("err send message in kafka ", err)
	}
	log.Printf("[producer] partition id: %d; offset:%d, value: %v\n", partition, offset, msg)
	return nil
}
