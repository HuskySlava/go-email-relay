package main

import (
	"fmt"
	"github.com/HuskySlava/go-email-relay/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}
	fmt.Println(cfg)
}
