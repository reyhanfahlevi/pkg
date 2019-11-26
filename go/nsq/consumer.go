package nsq

import (
	"github.com/nsqio/go-nsq"
)

// Consumer instance
type Consumer struct {
	listenAddress []string
	prefix        string

	handlers     []ConsumerHandler
	nsqConsumers []*nsq.Consumer
}

// ConsumerConfig config for the consumer instance
type ConsumerConfig struct {
	ListenAddress []string
	Prefix        string
}

// ConsumerHandler handler for consumer
type ConsumerHandler struct {
	Topic       string
	Channel     string
	Concurrent  int
	MaxAttempts uint16
	MaxInFlight int
	Enable      bool

	Handler func(message IMessage) error
}

// NewConsumer will instantiate the nsq consumer
func NewConsumer(cfg ConsumerConfig) *Consumer {
	return &Consumer{
		listenAddress: cfg.ListenAddress,
		prefix:        cfg.Prefix,
	}
}

// RegisterHandler will register the consumer handlers
func (c *Consumer) RegisterHandler(handler ConsumerHandler) {
	if handler.Enable {
		c.handlers = append(c.handlers, handler)
	}
}

// Run will connecting all registered consumer handlers to the nsqlookupd address
func (c *Consumer) Run() error {
	for _, h := range c.handlers {
		cfg := nsq.NewConfig()
		cfg.MaxAttempts = h.MaxAttempts
		cfg.MaxInFlight = h.MaxInFlight
		q, err := nsq.NewConsumer(h.Topic, h.Channel, cfg)
		if err != nil {
			return err
		}

		if h.Concurrent != 0 {
			q.AddConcurrentHandlers(c.handle(h.Handler), h.Concurrent)
		} else {
			q.AddHandler(c.handle(h.Handler))
		}

		err = q.ConnectToNSQLookupds(c.listenAddress)
		if err != nil {
			err = q.ConnectToNSQDs(c.listenAddress)
		}

		if err != nil {
			return err
		}
		c.nsqConsumers = append(c.nsqConsumers, q)
	}

	return nil
}

// handle will convert the func(msg Message) error into nsq.HandlerFunc
func (c *Consumer) handle(fn func(IMessage) error) nsq.HandlerFunc {
	return func(message *nsq.Message) error {
		msg := &Message{
			message,
		}

		return fn(msg)
	}
}
