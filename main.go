package main

import (
	"log"

	"github.com/Satishcg12/multicommers/internal"
	"github.com/Satishcg12/multicommers/utils/dotenv"
)

func main() {
	// load config
	dotenv.LoadConfig()
	err := internal.NewServer(internal.ServerConfig{
		Host: dotenv.GetEnvOrDefault("HOST", "localhost"),
		Port: dotenv.GetEnvOrDefault("PORT", "8080"),
	}).Start()
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
