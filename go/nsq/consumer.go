package nsq

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nsqio/go-nsq"
)

// Consumer instance
type Consumer struct {
	listenAddress []string
	prefix        string

	handlers     []ConsumerHandler
	nsqConsumers []*nsq.Consumer
	stopTimeout  time.Duration
}

// ConsumerConfig config for the consumer instance
type ConsumerConfig struct {
	ListenAddress []string
	Prefix        string
	StopTimeout   time.Duration
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

// RunDirect will connecting all registered consumer handlers directly to the nsqd address
func (c *Consumer) RunDirect() error {
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

		err = q.ConnectToNSQDs(c.listenAddress)
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

// Wait waits for the stop/restart signal and shutdown the NSQ consumers
// gracefully
func (c *Consumer) Wait() {
	<-WaitTermSig(c.stop)
}

// stop the nsq consumers gracefully
func (c *Consumer) stop(ctx context.Context) error {
	var wg sync.WaitGroup
	for _, con := range c.nsqConsumers {
		wg.Add(1)
		con := con
		go func() { // use goroutines to stop all of them ASAP
			defer wg.Done()
			con.Stop()

			select {
			case <-con.StopChan:
			case <-ctx.Done():
			case <-time.After(c.stopTimeout):
			}
		}()
	}
	wg.Wait()
	return nil
}

// WaitTermSig wait for termination signal
func WaitTermSig(handler func(context.Context) error) <-chan struct{} {
	stoppedCh := make(chan struct{})
	go func() {
		signals := make(chan os.Signal, 1)

		// wait for the sigterm
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-signals

		// We received an os signal, shut down.
		if err := handler(context.Background()); err != nil {
			log.Printf("graceful shutdown  fail: %v", err)
		} else {
			log.Println("gracefull shutdown success")
		}

		close(stoppedCh)

	}()
	return stoppedCh
}
