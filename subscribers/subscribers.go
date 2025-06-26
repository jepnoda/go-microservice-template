package subscribers

import (
	"context"
	"fmt"
	"go-microservice-template/logger"
	"go-microservice-template/redis"
	"sync"

	goredis "github.com/redis/go-redis/v9"
)

// MessageHandler defines the function signature for handling messages
type MessageHandler func(channel, message string) error

// Subscriber represents a Redis subscriber
type Subscriber struct {
	channel string
	handler MessageHandler
	pubsub  *goredis.PubSub
	ctx     context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

// SubscriberManager manages multiple subscribers
type SubscriberManager struct {
	subscribers map[string]*Subscriber
	mu          sync.RWMutex
}

// NewSubscriberManager creates a new subscriber manager
func NewSubscriberManager() *SubscriberManager {
	return &SubscriberManager{
		subscribers: make(map[string]*Subscriber),
	}
}

// Subscribe creates a new subscriber for the given channel
func (sm *SubscriberManager) Subscribe(channel string, handler MessageHandler) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Check if already subscribed to this channel
	if _, exists := sm.subscribers[channel]; exists {
		return fmt.Errorf("already subscribed to channel: %s", channel)
	}

	// Create Redis subscription
	pubsub := redis.Subscribe(channel)
	if pubsub == nil {
		return fmt.Errorf("failed to subscribe to channel: %s", channel)
	}

	// Create context for this subscriber
	ctx, cancel := context.WithCancel(context.Background())

	subscriber := &Subscriber{
		channel: channel,
		handler: handler,
		pubsub:  pubsub,
		ctx:     ctx,
		cancel:  cancel,
	}

	sm.subscribers[channel] = subscriber

	// Start listening for messages
	subscriber.wg.Add(1)
	go subscriber.listen()

	logger.Info(fmt.Sprintf("Successfully subscribed to channel: %s", channel))
	return nil
}

// Unsubscribe removes a subscriber for the given channel
func (sm *SubscriberManager) Unsubscribe(channel string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	subscriber, exists := sm.subscribers[channel]
	if !exists {
		return fmt.Errorf("not subscribed to channel: %s", channel)
	}

	// Cancel the context and wait for goroutine to finish
	subscriber.cancel()
	subscriber.wg.Wait()

	// Close the pubsub connection
	if err := subscriber.pubsub.Close(); err != nil {
		logger.Error(fmt.Sprintf("Error closing pubsub for channel %s: %v", channel, err))
	}

	delete(sm.subscribers, channel)
	logger.Info(fmt.Sprintf("Successfully unsubscribed from channel: %s", channel))
	return nil
}

// UnsubscribeAll removes all subscribers
func (sm *SubscriberManager) UnsubscribeAll() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for channel := range sm.subscribers {
		subscriber := sm.subscribers[channel]
		subscriber.cancel()
		subscriber.wg.Wait()

		if err := subscriber.pubsub.Close(); err != nil {
			logger.Error(fmt.Sprintf("Error closing pubsub for channel %s: %v", channel, err))
		}
	}

	sm.subscribers = make(map[string]*Subscriber)
	logger.Info("All subscribers have been unsubscribed")
}

// GetSubscribedChannels returns a list of currently subscribed channels
func (sm *SubscriberManager) GetSubscribedChannels() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	channels := make([]string, 0, len(sm.subscribers))
	for channel := range sm.subscribers {
		channels = append(channels, channel)
	}
	return channels
}

// IsSubscribed checks if subscribed to a specific channel
func (sm *SubscriberManager) IsSubscribed(channel string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	_, exists := sm.subscribers[channel]
	return exists
}

// listen listens for messages on the subscriber's channel
func (s *Subscriber) listen() {
	defer s.wg.Done()

	ch := s.pubsub.Channel()

	for {
		select {
		case <-s.ctx.Done():
			logger.Info(fmt.Sprintf("Stopping subscriber for channel: %s", s.channel))
			return
		case msg, ok := <-ch:
			if !ok {
				logger.Error(fmt.Sprintf("Channel closed for subscriber: %s", s.channel))
				return
			}

			// Handle the message
			if err := s.handler(msg.Channel, msg.Payload); err != nil {
				logger.Error(fmt.Sprintf("Error handling message from channel %s: %v", msg.Channel, err))
			}
		}
	}
}
