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

func example() {
	logger.Info("Starting microservice...")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Load configuration and initialize Redis
	config.LoadConfig()
	redis.InitRedis()

	// Create subscriber manager
	manager := subscribers.NewSubscriberManager()
	defer manager.UnsubscribeAll()

	// Option 1: Set up only game result subscriber
	if err := subscribers.SetupGameResultSubscriber(manager); err != nil {
		logger.Error("Failed to set up game result subscriber: " + err.Error())
		return
	}

	// Option 2: Uncomment the following line to set up all subscribers
	// if err := subscribers.SetupAllSubscribers(manager); err != nil {
	//     logger.Error("Failed to set up subscribers: " + err.Error())
	//     return
	// }

	// Option 3: Set up individual subscribers as needed
	// if err := subscribers.SetupPlayerActionSubscriber(manager); err != nil {
	//     logger.Error("Failed to set up player action subscriber: " + err.Error())
	//     return
	// }
	// if err := subscribers.SetupChatMessageSubscriber(manager); err != nil {
	//     logger.Error("Failed to set up chat message subscriber: " + err.Error())
	//     return
	// }

	logger.Info("Microservice is running. Press Ctrl+C to stop.")

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info("Shutting down microservice...")
}
