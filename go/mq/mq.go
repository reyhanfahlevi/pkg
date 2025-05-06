package mq

import "github.com/reyhanfahlevi/pkg/go/mq/adapter"

type MessageQueue struct {
	consumer adapter.IConsumerAdapter
}

type ConsumerConfig struct {
	Topic       string `json:"topic"`
	Channel     string `json:"channel,omitempty"`
	Concurrent  int    `json:"concurrent,omitempty"`
	MaxAttempts int32  `json:"max_attempts,omitempty"`
	MaxInFlight int    `json:"max_in_flight,omitempty"`
	Enable      bool   `json:"enable,omitempty"`
	URL         string `json:"url"`
	ExtraConfig map[string]interface{}
}

// New will create mq client that can receive consumer
// the consumer can be rmq, nsq, or kafka if needed
func New(consumer adapter.IConsumerAdapter) *MessageQueue {
	return &MessageQueue{consumer: consumer}
}

func (mq *MessageQueue) RegisterConsumerHandler(consumerConfig ConsumerConfig, handler adapter.Handler) error {
	cfg := adapter.ConsumerHandler{
		Topic:       consumerConfig.Topic,
		Channel:     consumerConfig.Channel,
		Concurrent:  consumerConfig.Concurrent,
		MaxAttempts: consumerConfig.MaxAttempts,
		MaxInFlight: consumerConfig.MaxInFlight,
		Enable:      consumerConfig.Enable,
		URL:         consumerConfig.URL,

		Handler: handler,
	}
	cfg.SetExtraConfig(consumerConfig.ExtraConfig)
	return mq.consumer.RegisterConsumerHandler(cfg)
}

func (mq *MessageQueue) RunConsumer() error {
	return mq.consumer.Run()
}
