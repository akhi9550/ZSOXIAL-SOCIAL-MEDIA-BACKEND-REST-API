package main

import (
	"log"

	_ "github.com/akhi9550/api-gateway/cmd/docs"
	"github.com/akhi9550/api-gateway/pkg/config"
	"github.com/akhi9550/api-gateway/pkg/di"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Zsoxial_Microservice_CleanArchitecture
// @version 1.0.0
// @description Zsoxial is an Social media platform to interact with peoples
// @contact.name github.com/akhi9550
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @host localhost:8000
// @BasePath /
// @query.collection.format multi

func main() {
	u := gin.Default()
	u.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
