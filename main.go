package main

import (
	"context"
	"go-microservice-template/config"
	"go-microservice-template/logger"
	"go-microservice-template/redis"
	"go-microservice-template/server"
	"go-microservice-template/subscribers"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	logger.Info("Starting microservice...")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Load configuration and initialize Redis
	config.LoadConfig()
	redis.InitRedis()

	// Create subscriber manager
	manager := subscribers.NewSubscriberManager()
	defer manager.UnsubscribeAll()

	// Set up game result subscriber
	if err := subscribers.SetupGameResultSubscriber(manager); err != nil {
		logger.Error("Failed to set up game result subscriber: " + err.Error())
		return
	}

	// Create and start HTTP server
	httpServer := server.NewServer()

	// Use WaitGroup to manage goroutines
	var wg sync.WaitGroup

	// Start HTTP server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(); err != nil {
			logger.Error("HTTP server error: " + err.Error())
		}
	}()

	logger.Info("Microservice is running. Press Ctrl+C to stop.")

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info("Shutting down microservice...")

	// Create a context with timeout for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server gracefully
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during HTTP server shutdown: " + err.Error())
	}

	// Wait for all goroutines to finish
	wg.Wait()
	logger.Info("Microservice shutdown complete")
}
