package config

import "github.com/spf13/viper"

type Config struct {
	Port       string `mapstructure:"PORT"`
	AuthSvcUrl string `mapstructure:"AUTH_SVC_URL"`
	PostSvcUrl string `mapstructure:"POST_SVC_URL"`
	ChatSvcUrl string `mapstructure:"CHAT_SVC_URL"`
	NotificationSvcUrl string `mapstructure:"NOTIFICATION_SVC_URL"`


	KafkaPort  string `mapstructure:"KAFKA_PORT"`
	KafkaTopic string `mapstructure:"KAFKA_TOPIC"`
	KafkaLikeTopic    string `mapstructure:"KAFKA__LIKE_TOPIC"`
	KafkaCommentTopic string `mapstructure:"KAFKA_COMMENT_TOPIC"`

	KEY       string `mapstructure:"KEY"`
	KEY_ADMIN string `mapstructure:"KEY_ADMIN"`
}

var envs = []string{
	"PORT", "AUTH_SVC_URL", "POST_SVC_URL", "CHAT_SVC_URL", "KEY", "KEY_ADMIN","NOTIFICATION_SVC_URL","KAFKA_PORT","KAFKA_TOPIC","KAFKA__LIKE_TOPIC","KAFKA_COMMENT_TOPIC",
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath("./")
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil

}
