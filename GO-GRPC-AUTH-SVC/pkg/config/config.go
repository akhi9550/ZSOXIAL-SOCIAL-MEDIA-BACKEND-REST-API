package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBName     string `mapstructure:"DB_NAME"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	Port       string `mapstructure:"PORT"`

	PostSvcUrl string `mapstructure:"POST_SVC_URL"`
	ChatSvcUrl string `mapstructure:"CHAT_SVC_URL"`

	KEY       string `mapstructure:"KEY"`
	KEY_ADMIN string `mapstructure:"KEY_ADMIN"`

	AUTHTOKEN   string `mapstructure:"TWILIO_AUTHTOKEN"`
	ACCOUNTSID  string `mapstructure:"TWILIO_ACCOUNTSID"`
	SERVICESSID string `mapstructure:"TWILIO_SERVICESID"`

	AWS_REGION            string `mapstructure:"AWS_REGION"`
	AWS_ACCESS_KEY_ID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AWS_SECRET_ACCESS_KEY string `mapstructure:"AWS_SECRET_ACCESS_KEY"`

	KafkaPort  string `mapstructure:"KAFKA_PORT"`
	KafkaTopic string `mapstructure:"KAFKA_TOPIC"`
}

var envs = []string{
	"DB_HOST", "DB_NAME", "DB_USER", "DB_PORT", "DB_PASSWORD", "PORT", "POST_SVC_URL", "CHAT_SVC_URL", "KEY", "KEY_ADMIN", "TWILIO_AUTHTOKEN", "TWILIO_ACCOUNTSID", "TWILIO_SERVICESID", "AWS_REGION", "AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY", "KAFKA_PORT", "KAFKA_TOPIC",
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
