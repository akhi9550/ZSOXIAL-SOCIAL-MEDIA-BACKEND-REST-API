package main

import (
	"log"

	"github.com/akhi9550/notification-svc/pkg/config"
	"github.com/akhi9550/notification-svc/pkg/di"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	server, err := di.InitializeAPI(config)
	if err != nil {
		log.Fatal("cannot start server:", err)
	} else {
		server.Start()
	}
}
