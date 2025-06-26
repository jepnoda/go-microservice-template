package main

import (
	"context"
	"go-microservice-template/config"
	"go-microservice-template/logger"
	"go-microservice-template/redis"
	"go-microservice-template/subscribers"
	"os/signal"
	"syscall"
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

	logger.Info("Microservice is running. Press Ctrl+C to stop.")

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info("Shutting down microservice...")
}
