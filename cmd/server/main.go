// @title Todo API
// @version 1.0
// @description A Clean Architecture Todo API with Go, Fiber, and PostgreSQL
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := app.Start(); err != nil {
			log.Printf("Failed to start app: %v", err)
		}
	}()

	log.Println("Application started successfully")

	// Wait for interrupt signal
	<-quit
	log.Println("Received shutdown signal")

	// Graceful shutdown
	if err := app.Stop(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
} 