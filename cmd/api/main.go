package main

import (
	"context"
	"log"
	"transaction-server/app"
	"transaction-server/app/boot"
	"transaction-server/internal/routes"
)

func main() {
	// This is the entry point for the application.
	// It should start the server and listen for incoming
	// requests.
	// It should also initialize the application context
	ctx := context.Background()
	if err := boot.Initialize(ctx); err != nil {
		log.Fatalf("failed to initialize the application: %v", err)
	}
	router := routes.RegisterRoutes(ctx)
	err := router.Run(app.Context().Config().App.Port)
	if err != nil {
		log.Fatal("Unable to start server. Error: ", err)
	}
}
