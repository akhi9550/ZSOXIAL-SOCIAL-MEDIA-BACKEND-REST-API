package helper

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/akhi9550/auth-svc/pkg/config"
	"github.com/akhi9550/auth-svc/pkg/utils/models"
)

func SendFollowReqNotification(data models.Notification, msg []byte) {
	data.Message = string(msg)
	err := KafkaFollowReqProducer(data)
	if err != nil {
		fmt.Println("error sending notification to Kafka:", err)
		return
	}

	fmt.Println("==sent followreq successfully to Kafka==")
}

func KafkaFollowReqProducer(message models.Notification) error {
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

	msg := &sarama.ProducerMessage{Topic: cfg.KafkaFollowReqTopic, Key: sarama.StringEncoder("NotificationFollowReq"), Value: sarama.StringEncoder(result)}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("err send message in kafka ", err)
	}
	log.Printf("[producer] partition id: %d; offset:%d, value: %v\n", partition, offset, msg)
	return nil
}

func SendAcceptReqNotification(data models.Notification, msg []byte) {
	data.Message = string(msg)
	err := KafkaAcceptReqProducer(data)
	if err != nil {
		fmt.Println("error sending notification to Kafka:", err)
		return
	}

	fmt.Println("==sent acceptreq successfully to Kafka==")
}

func KafkaAcceptReqProducer(message models.Notification) error {
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

	msg := &sarama.ProducerMessage{Topic: cfg.KafkaAcceptReqTopic, Key: sarama.StringEncoder("NotificationAcceptReq"), Value: sarama.StringEncoder(result)}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("err send message in kafka ", err)
	}
	log.Printf("[producer] partition id: %d; offset:%d, value: %v\n", partition, offset, msg)
	return nil
}
