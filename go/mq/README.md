# Message Queue (MQ) Package

A lightweight and flexible message queue implementation in Go that provides a unified interface for working with different message queue adapters.

## Features

- Adapter-based architecture supporting multiple message queue implementations
- Built-in RabbitMQ adapter with robust features
- Configurable consumer settings
- Automatic reconnection handling
- Exponential backoff with jitter for failed messages
- Thread-safe operations
- Concurrent message processing
- Message retry mechanism with configurable attempts

## Installation

```go
go get github.com/reyhanfahlevi/pkg/go/mq
```

## Usage
### Creating a Message Queue Client

```go
import (
    "github.com/reyhanfahlevi/pkg/go/mq"
    "github.com/reyhanfahlevi/pkg/go/mq/adapter/rmqa"
)

// Initialize with RabbitMQ adapter
consumer := rmqa.NewManager()
mqClient := mq.New(consumer)
```

### Configuring and Registering a Consumer

```go
config := mq.ConsumerConfig{
    Topic:       "my-topic",          // Queue name
    Channel:     "my-channel",        // Consumer identifier
    Concurrent:  5,                   // Number of concurrent workers
    MaxAttempts: 3,                   // Maximum retry attempts
    MaxInFlight: 10,                  // Maximum unacknowledged messages
    Enable:      true,
    URL:         "amqp://localhost:5672",
    ExtraConfig: map[string]interface{}{
        "durable":     true,          // Queue survives broker restart
        "auto_delete": false,         // Queue remains when consumer disconnects
        "exclusive":   false,         // Queue can be used by other consumers
        "no_wait":     false,         // Wait for server confirmation
    },
}

// Define message handler
handler := func(ctx context.Context, msg adapter.IMessage) error {
    // Process your message here
    data := msg.GetBody()
    
    // Your message processing logic
    
    // Acknowledge successful processing
    msg.Finish()
    return nil
}

// Register the consumer
err := mqClient.RegisterConsumerHandler(config, handler)
if err != nil {
    log.Fatal(err)
}
```

### Starting the Consumer
```go
err := mqClient.RunConsumer()
if err != nil {
    log.Fatal(err)
}
```

## Message Handling
The package provides several methods for handling messages:

### Message Interface
```go
type IMessage interface {
    Finish()                                    // Acknowledge message
    RequeueWithoutBackoff(delay time.Duration) // Requeue without exponential backoff
    Requeue(delay time.Duration)               // Requeue with exponential backoff
    GetAttempts() int32                        // Get processing attempts count
    GetBody() []byte                           // Get message payload
}
```

### Error Handling and Retries
The package implements a sophisticated retry mechanism with exponential backoff and jitter:

- Failed messages are automatically requeued
- Exponential backoff prevents overwhelming the system
- Random jitter helps prevent thundering herd problems
- Maximum retry attempts are configurable
- Maximum backoff delay is capped at 10 minutes

## Features in Detail
### Adapter-Based Architecture
The package uses an adapter interface that allows implementing different message queue backends:
```go
type IConsumerAdapter interface {
    RegisterConsumerHandler(consumer ConsumerHandler) error
    Run() error
}
```

### RabbitMQ Adapter Features
The built-in RabbitMQ adapter provides:

- Automatic connection management
- Channel pooling
- Queue declaration with configurable parameters
- Message acknowledgment handling
- Dead letter exchange support
- Consumer concurrency control
- Quality of Service (QoS) settings
### Thread Safety
All operations are thread-safe and can be used in concurrent environments:

- Mutex-protected connection handling
- Thread-safe message processing
- Concurrent worker support
- Safe shutdown handling
## Contributing
Contributions are welcome! Please feel free to submit a Pull Request