package main

import (
	"github.com/HuskySlava/go-email-relay/internal/config"
	"github.com/HuskySlava/go-email-relay/internal/server"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	if err := server.StartServer(&cfg.HTTPConfig); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
