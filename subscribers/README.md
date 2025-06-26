# Subscribers Package

This package provides a clean and organized way to manage Redis subscribers in the microservice.

## Features

- **SubscriberManager**: Manages multiple Redis subscribers
- **Message Handling**: Type-safe message handling with custom handler functions
- **Graceful Shutdown**: Proper cleanup of all subscriptions
- **Concurrent Safe**: Thread-safe operations using mutexes
- **Game Result Handler**: Pre-built handler for game result messages

## Usage

### Basic Usage

```go
// Create a subscriber manager
manager := subscribers.NewSubscriberManager()
defer manager.UnsubscribeAll()

// Subscribe to a channel with a custom handler
err := manager.Subscribe("my-channel", func(channel, message string) error {
    fmt.Printf("Received: %s from %s\n", message, channel)
    return nil
})

if err != nil {
    log.Fatal(err)
}
```

### Game Results Subscriber

For game results, you can use the pre-built handler:

```go
manager := subscribers.NewSubscriberManager()
defer manager.UnsubscribeAll()

// Set up game result subscriber
err := subscribers.SetupGameResultSubscriber(manager)
if err != nil {
    log.Fatal(err)
}
```

### Custom Message Handlers

You can create custom message handlers for different types of messages:

```go
func MyCustomHandler(channel, message string) error {
    // Parse message
    var data MyDataType
    if err := json.Unmarshal([]byte(message), &data); err != nil {
        return err
    }
    
    // Process the data
    // ...
    
    return nil
}

// Use the custom handler
manager.Subscribe("my-custom-channel", MyCustomHandler)
```

## API Reference

### SubscriberManager

#### Methods

- `NewSubscriberManager() *SubscriberManager` - Creates a new subscriber manager
- `Subscribe(channel string, handler MessageHandler) error` - Subscribe to a channel
- `Unsubscribe(channel string) error` - Unsubscribe from a channel
- `UnsubscribeAll()` - Unsubscribe from all channels
- `GetSubscribedChannels() []string` - Get list of subscribed channels
- `IsSubscribed(channel string) bool` - Check if subscribed to a channel

### MessageHandler

```go
type MessageHandler func(channel, message string) error
```

A function type for handling messages. Should return an error if message processing fails.

### GameResult

```go
type GameResult struct {
    GameID    string `json:"game_id"`
    PlayerID  string `json:"player_id"`
    Result    string `json:"result"`
    Score     int    `json:"score"`
    Timestamp string `json:"timestamp"`
}
```

Structure for game result messages.

## Error Handling

The package uses the logger package for error reporting. Make sure the logger is properly initialized before using subscribers.

## Thread Safety

All operations are thread-safe and can be called from multiple goroutines safely.
