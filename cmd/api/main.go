package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/orchid-craft/orchid-auth-starter/internal/bootstrap/container"
	restfulServer "github.com/orchid-craft/orchid-auth-starter/internal/bootstrap/server/restful_server"
)

func main() {
	// 1. Initialize bootstrap container (Config, DI, Router, Logger)
	appContainer, err := container.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize application container: %v", err)
	}
	defer appContainer.Close()

	// 2. Initialize and start RESTful HTTP Server
	srv := restfulServer.NewServer(appContainer)
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
