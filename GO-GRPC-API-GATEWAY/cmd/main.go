package main

import (
	docs "github.com/akhi9550/api-gateway/cmd/docs"
	"github.com/akhi9550/api-gateway/pkg/config"
	"github.com/akhi9550/api-gateway/pkg/di"
	"github.com/akhi9550/api-gateway/pkg/logging"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

func main() {
	docs.SwaggerInfo.Title = "Zsoxial_Microservice_CleanArchitecture"
	docs.SwaggerInfo.Description = "Zsoxial is a Social media platform to interact with people"
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	logging.Init()
	logger := logging.GetLogger()

	logEntry := logger.WithField("context", "loading config")
	config, configErr := config.LoadConfig()
	if configErr != nil {
		logEntry.Fatal("cannot load config: ", configErr)
	}

	logEntry = logger.WithField("context", "initializing API")
	server, diErr := di.InitializeAPI(config)
	if diErr != nil {
		logEntry.WithError(diErr).Fatal("cannot start server")
	}
	server.Start()
}
