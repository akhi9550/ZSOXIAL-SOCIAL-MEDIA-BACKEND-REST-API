package main

import (
	"log"

	"github.com/akhi9550/post-svc/pkg/config"
	"github.com/akhi9550/post-svc/pkg/di"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	server, err := di.InitializeAPI(cfg)
	if err != nil {
		log.Fatal("cannot start server:", err)
	} else {
		server.Start()
	}
}
