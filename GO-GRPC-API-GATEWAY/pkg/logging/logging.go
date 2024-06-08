package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Init() {
	log.SetOutput(os.Stdout)

	log.SetLevel(logrus.DebugLevel)

	log.SetFormatter(&logrus.JSONFormatter{})
}

func GetLogger() *logrus.Logger {
	return log
}