package helper

import "github.com/IBM/sarama"

func ConnectToLikeConsumer(brokerUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokerUrl, config)

	if err != nil {
		return nil, err
	}
	return conn, nil
}

func ConnectToCommentConsumer(brokerUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	conn, err := sarama.NewConsumer(brokerUrl, config)

	if err != nil {
		return nil, err
	}
	return conn, nil
}
