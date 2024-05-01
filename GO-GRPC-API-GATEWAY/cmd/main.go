package main

import (
	"log"

	"github.com/akhi9550/api-gateway/cmd/docs"
	_ "github.com/akhi9550/api-gateway/cmd/docs"
	"github.com/akhi9550/api-gateway/pkg/config"
	"github.com/akhi9550/api-gateway/pkg/di"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

func main() {
	docs.SwaggerInfo.Title = "Zsoxial_Microservice_CleanArchitecture"
	docs.SwaggerInfo.Description = "Zsoxial is an Social media platform to interact with peoples"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8000"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	config, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	server, diErr := di.InitializeAPI(config)

	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}
}
