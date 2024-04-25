package helper

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/akhi9550/post-svc/pkg/config"
	"github.com/akhi9550/post-svc/pkg/utils/models"
)

func SendLikeNotification(data models.Notification, msg []byte) {
	data.Message = string(msg)
	err := KafkaLikeProducer(data)
	if err != nil {
		fmt.Println("error sending notification to Kafka:", err)
		return
	}

	fmt.Println("==sent like successfully to Kafka==")
}

func KafkaLikeProducer(message models.Notification) error {
	cfg, _ := config.LoadConfig()
	configs := sarama.NewConfig()
	configs.Producer.Return.Successes = true
	configs.Producer.RequiredAcks = sarama.WaitForAll
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

func SendCommentNotification(data models.Notification, msg []byte) {
	data.Message = string(msg)
	err := KafkaCommentProducer(data)
	if err != nil {
		fmt.Println("error sending notification to Kafka:", err)
		return
	}

	fmt.Println("==sent Comment successfully to Kafka==")
}

func KafkaCommentProducer(message models.Notification) error {
	cfg, _ := config.LoadConfig()
	configs := sarama.NewConfig()
	configs.Producer.Return.Successes = true
	configs.Producer.RequiredAcks = sarama.WaitForAll
	configs.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer([]string{cfg.KafkaPort}, configs)
	if err != nil {
		return err
	}

	result, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{Topic: cfg.KafkaCommentTopic, Key: sarama.StringEncoder("NotificationComment"), Value: sarama.StringEncoder(result)}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("err send message in kafka ", err)
	}
	log.Printf("[producer] partition id: %d; offset:%d, value: %v\n", partition, offset, msg)
	return nil
}
