package main

import (
	"log"
	"stark8/internal/api"
	"stark8/internal/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Error loading config", err)
	}
	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal("Error creating server", err)
	}
	server.Start()
}
