package rmqa

import (
	"context"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/reyhanfahlevi/pkg/go/log"
	"github.com/reyhanfahlevi/pkg/go/mq/adapter"
	"golang.org/x/exp/rand"
)

type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	mu           sync.Mutex
	connected    bool
	notifyClose  chan *amqp.Error
	shutdown     chan struct{}
	reconnecting bool

	// Single handler configuration
	handler       adapter.ConsumerHandler
	handlerConfig HandlerConfig
	isConfigured  bool
	baseDelay     time.Duration
	maxDelay      time.Duration
}

type HandlerConfig struct {
	Durable    bool `json:"durable,omitempty"`
	AutoDelete bool `json:"auto_delete,omitempty"`
	Exclusive  bool `json:"exclusive,omitempty"`
	NoWait     bool `json:"no_wait,omitempty"`
	AutoAck    bool `json:"auto_ack,omitempty"`
	NoLocal    bool `json:"no_local,omitempty"`
}

type ConsumerOptions func(*Consumer)

func NewConsumer(handler adapter.ConsumerHandler, opt ...ConsumerOptions) *Consumer {
	c := &Consumer{
		shutdown:  make(chan struct{}),
		handler:   handler,
		baseDelay: time.Second,      // Base delay of 1 second
		maxDelay:  time.Minute * 10, // Maximum delay of 10 minutes
	}

	for _, opt := range opt {
		opt(c)
	}

	return c
}

func (r *Consumer) connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.connected {
		return nil
	}

	conn, err := amqp.Dial(r.handler.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	r.conn = conn
	r.connected = true
	r.notifyClose = make(chan *amqp.Error)
	r.conn.NotifyClose(r.notifyClose)

	return nil
}

func (r *Consumer) setupConsumer() error {
	r.mu.Lock()

	// Configure the consumer handler here instead of during initialization
	if r.isConfigured {
		r.mu.Unlock()
		return fmt.Errorf("consumer handler already configured")
	}

	if r.handler.Concurrent < 1 {
		r.handler.Concurrent = 1
	}

	if r.handler.MaxAttempts < 1 {
		r.handler.MaxAttempts = 1
	}

	if r.handler.MaxInFlight < 1 {
		r.handler.MaxInFlight = 1
	}

	handlerConfig := HandlerConfig{}
	_ = r.handler.ParseExtraConfig(&handlerConfig)

	r.handlerConfig = handlerConfig
	r.isConfigured = true
	r.mu.Unlock()

	ch, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	r.channel = ch

	// Declare queue
	_, err = ch.QueueDeclare(
		r.handler.Topic,
		r.handlerConfig.Durable,    // durable
		r.handlerConfig.AutoDelete, // auto-delete
		r.handlerConfig.Exclusive,  // exclusive
		r.handlerConfig.NoWait,     // no-wait
		nil,                        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Set QoS
	err = ch.Qos(
		r.handler.MaxInFlight,
		0,
		false,
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Start consuming
	deliveries, err := ch.Consume(
		r.handler.Topic,
		r.handler.Channel,
		r.handlerConfig.AutoAck,   // auto-ack
		r.handlerConfig.Exclusive, // exclusive
		r.handlerConfig.NoLocal,   // no-local
		r.handlerConfig.NoWait,    // no-wait
		nil,                       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	for i := 0; i < r.handler.Concurrent; i++ {
		go func(workerID int) {
			for delivery := range deliveries {
				ctx := context.Background()

				msg := &Message{
					Delivery:    delivery,
					maxAttempts: r.handler.MaxAttempts,
					topic:       r.handler.Topic,
					ch:          r.channel,
				}

				msg.increaseAttempts()
				err := r.handler.Handler(ctx, msg)
				if err != nil {
					log.Error(err)

					if msg.requeued {
						continue
					}

					// Calculate exponential backoff delay
					attempts := msg.GetAttempts()
					delay := r.calculateBackoff(attempts)
					msg.Requeue(delay)
					continue
				}
				msg.Ack(false)
			}
		}(i)
	}

	return nil
}

// Add this new method for calculating backoff
func (r *Consumer) calculateBackoff(attempts int32) time.Duration {
	// Calculate exponential delay: baseDelay * 2^(attempts-1)
	delay := r.baseDelay * time.Duration(1<<uint(attempts-1))

	// Add some jitter to prevent thundering herd
	jitter := time.Duration(rand.Int63n(int64(delay) / 2))
	delay = delay + jitter

	// Ensure we don't exceed maximum delay
	if delay > r.maxDelay {
		delay = r.maxDelay
	}

	return delay
}

func (r *Consumer) Run() error {
	if r.handler.Handler == nil {
		return fmt.Errorf("no consumer handler specified")
	}

	if err := r.connect(); err != nil {
		return err
	}

	go r.reconnectLoop()

	return r.setupConsumer()
}

func (r *Consumer) reconnectLoop() {
	for {
		select {
		case <-r.shutdown:
			return
		case err := <-r.notifyClose:
			if err != nil {
				r.reconnect()
			}
		}
	}
}

func (r *Consumer) reconnect() {
	r.mu.Lock()
	if r.reconnecting {
		r.mu.Unlock()
		return
	}
	r.reconnecting = true
	r.connected = false
	r.isConfigured = false // Reset the configuration flag
	r.mu.Unlock()

	backoff := time.Second
	maxBackoff := 30 * time.Second

	for {
		select {
		case <-r.shutdown:
			return
		default:
			time.Sleep(backoff)

			if err := r.connect(); err != nil {
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
				continue
			}

			if err := r.setupConsumer(); err != nil {
				r.mu.Lock()
				r.connected = false
				r.mu.Unlock()
				continue
			}

			r.mu.Lock()
			r.reconnecting = false
			r.mu.Unlock()
			return
		}
	}
}

// Close gracefully shuts down the consumer
func (r *Consumer) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	close(r.shutdown)

	if r.channel != nil {
		r.channel.Close()
	}

	if r.conn != nil {
		return r.conn.Close()
	}

	return nil
}
