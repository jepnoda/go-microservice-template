package subscribers

import (
	"encoding/json"
	"fmt"
	"go-microservice-template/logger"
)

// PlayerAction represents a player action message
type PlayerAction struct {
	PlayerID  string `json:"player_id"`
	Action    string `json:"action"`
	Details   string `json:"details"`
	Timestamp string `json:"timestamp"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	PlayerID  string `json:"player_id"`
	RoomID    string `json:"room_id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

// PlayerActionHandler handles messages from player action channels
func PlayerActionHandler(channel, message string) error {
	logger.Info(fmt.Sprintf("Received player action from channel %s: %s", channel, message))

	var action PlayerAction
	if err := json.Unmarshal([]byte(message), &action); err != nil {
		logger.Info(fmt.Sprintf("Raw player action message (not JSON): %s", message))
		return nil
	}

	logger.Info(fmt.Sprintf("Processing player action: PlayerID=%s, Action=%s, Details=%s",
		action.PlayerID, action.Action, action.Details))

	// Add your player action processing logic here
	return nil
}

// ChatMessageHandler handles messages from chat channels
func ChatMessageHandler(channel, message string) error {
	logger.Info(fmt.Sprintf("Received chat message from channel %s: %s", channel, message))

	var chat ChatMessage
	if err := json.Unmarshal([]byte(message), &chat); err != nil {
		logger.Info(fmt.Sprintf("Raw chat message (not JSON): %s", message))
		return nil
	}

	logger.Info(fmt.Sprintf("Processing chat message: PlayerID=%s, RoomID=%s, Message=%s",
		chat.PlayerID, chat.RoomID, chat.Message))

	// Add your chat message processing logic here
	return nil
}

// SetupPlayerActionSubscriber sets up the player action subscriber
func SetupPlayerActionSubscriber(manager *SubscriberManager) error {
	return manager.Subscribe("player-actions", PlayerActionHandler)
}

// SetupChatMessageSubscriber sets up the chat message subscriber
func SetupChatMessageSubscriber(manager *SubscriberManager) error {
	return manager.Subscribe("chat-messages", ChatMessageHandler)
}

// SetupAllSubscribers sets up all available subscribers
func SetupAllSubscribers(manager *SubscriberManager) error {
	// Game results
	if err := SetupGameResultSubscriber(manager); err != nil {
		return fmt.Errorf("failed to setup game result subscriber: %w", err)
	}

	// Player actions
	if err := SetupPlayerActionSubscriber(manager); err != nil {
		return fmt.Errorf("failed to setup player action subscriber: %w", err)
	}

	// Chat messages
	if err := SetupChatMessageSubscriber(manager); err != nil {
		return fmt.Errorf("failed to setup chat message subscriber: %w", err)
	}

	logger.Info("All subscribers have been set up successfully")
	return nil
}
