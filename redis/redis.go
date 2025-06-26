package redis

import (
	"context"
	"fmt"
	"go-microservice-template/config"
	"log"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

func InitRedis() {
	config := config.GetRedisConfig()

	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	}
	Client = redis.NewClient(opts)
	if err := Client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
		return
	}
	log.Println("Redis client initialized successfully.")
}

func Publish(channel string, message any) error {
	res, err := Client.Publish(ctx, channel, message).Result()
	if err != nil {
		log.Printf("Failed to publish message to channel %s: %v", channel, err)
		return err
	}
	if res == 0 {
		log.Printf("No subscribers for channel %s", channel)
		return nil
	}
	log.Printf("Message published to channel %s: %v", channel, message)
	return nil
}

func Subscribe(channel string) *redis.PubSub {
	pubsub := Client.Subscribe(ctx, channel)
	// Subscribe does not return an error directly, but we can check for errors by receiving a message
	_, err := pubsub.Receive(ctx)
	if err != nil {
		log.Printf("Failed to subscribe to channel %s: %v", channel, err)
		return nil
	}
	log.Printf("Subscribed to channel %s", channel)
	return pubsub
}
