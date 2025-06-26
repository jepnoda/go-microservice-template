package subscribers

import (
	"encoding/json"
	"fmt"
	"go-microservice-template/logger"
)

// GameResult represents the structure of a game result message
type GameResult struct {
	GameID    string `json:"game_id"`
	PlayerID  string `json:"player_id"`
	Result    string `json:"result"`
	Score     int    `json:"score"`
	Timestamp string `json:"timestamp"`
}

// GameResultHandler handles messages from the game-results channel
func GameResultHandler(channel, message string) error {
	logger.Info(fmt.Sprintf("Received message from channel %s: %s", channel, message))

	// Try to parse the message as JSON
	var gameResult GameResult
	if err := json.Unmarshal([]byte(message), &gameResult); err != nil {
		// If it's not valid JSON, just log the raw message
		logger.Info(fmt.Sprintf("Raw message (not JSON): %s", message))
		return nil
	}

	// Process the game result
	logger.Info(fmt.Sprintf("Processing game result: GameID=%s, PlayerID=%s, Result=%s, Score=%d",
		gameResult.GameID, gameResult.PlayerID, gameResult.Result, gameResult.Score))

	// Here you can add your business logic for processing game results
	// For example:
	// - Store in database
	// - Update player statistics
	// - Trigger notifications
	// - etc.

	return nil
}

// SetupGameResultSubscriber is a convenience function to set up the game result subscriber
func SetupGameResultSubscriber(manager *SubscriberManager) error {
	return manager.Subscribe("game-results", GameResultHandler)
}
