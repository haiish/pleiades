package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"pleiades/server/config/grpc"
	postgre "pleiades/server/config/postgre"
	"pleiades/server/config/redis"

	http "pleiades/server/api"
	// "pleiades/server/grpc"
	middlewares "pleiades/server/middlewares"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func main() {
	// Initialize DB and Redis
	postgre.InitDB()
	redis.InitRedis()

	// Set up Fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173, https://example2.com, http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = grpc.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to gRPC: %v", err)
	}

	// Set environment and logger if in development mode
	env := os.Getenv("ENV")
	if env == "development" {
		middlewares.SetupLogger(app)
		app.Use(logger.New())
	}

	// Load gRPC config
	cfg := grpc.LoadConfig()

	// Create a channel to listen for termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start the gRPC server in a goroutine
	go func() {
		if err := grpc.Start(cfg); err != nil {
			log.Fatalf("❌ Failed to start gRPC server: %v", err)
		}
	}()

	// Set up Swagger documentation for Fiber
	app.Static("/swagger", "./docs")
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger/v1.yaml",
	}))

	// Set up Fiber routes
	http.ApiV1Routes(app)

	// Start the Fiber server in a goroutine
	go func() {
		log.Fatal(app.Listen(":8081"))
	}()

	// Wait for a termination signal
	<-sigChan

	// Gracefully shut down Fiber
	log.Println("Shutting down Fiber...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Fiber shutdown error: %v", err)
	}

	<-ctx.Done()
	log.Println("Graceful shutdown completed")
}
