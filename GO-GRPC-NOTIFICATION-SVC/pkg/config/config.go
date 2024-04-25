package config

import "github.com/spf13/viper"

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	Port       string `mapstructure:"PORT"`

	AuthSvcUrl string `mapstructure:"AUTH_SVC_URL"`
	PostSvcUrl string `mapstructure:"POST_SVC_URL"`

	KafkaPort         string `mapstructure:"KAFKA_PORT"`
	KafkaLikeTopic    string `mapstructure:"KAFKA_LIKE_TOPIC"`
	KafkaCommentTopic string `mapstructure:"KAFKA_COMMENT_TOPIC"`
}

var envs = []string{
	"DB_URL", "DB_NAME", "PORT", "KAFKA_PORT", "KAFKA_LIKE_TOPIC", "KAFKA_COMMENT_TOPIC", "POST_SVC_URL",
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
