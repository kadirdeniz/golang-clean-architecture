package app

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kadirdeniz/golang-clean-architecture/internal/infrastructure/datasource"
)

type App interface {
	Start() error
	Stop() error
}

type app struct {
	server *fiber.App
	db     datasource.Database
}

func NewApp(server *fiber.App, db datasource.Database) App {
	return &app{
		server: server,
		db:     db,
	}
}

func (a *app) Start() error {
	// Open database connection
	if err := a.db.Open(context.Background()); err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return err
	}
	log.Println("Database connection established")

	// Start server
	log.Printf("Starting server on :8080")
	if err := a.server.Listen(":8080"); err != nil {
		log.Printf("Server error: %v", err)
		return err
	}

	return nil
}

func (a *app) Stop() error {
	log.Println("Initiating graceful shutdown...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := a.server.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("HTTP server stopped gracefully")
	}

	// Close database connection
	if err := a.db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	} else {
		log.Println("Database connection closed")
	}

	log.Println("Application shutdown complete")
	return nil
} 