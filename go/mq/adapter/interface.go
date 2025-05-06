package adapter

import (
	"context"
	"encoding/json"
	"time"
)

type IMessage interface {
	Finish()
	RequeueWithoutBackoff(delay time.Duration)
	Requeue(delay time.Duration)
	GetAttempts() int32
	GetBody() []byte
}

type IConsumerAdapter interface {
	RegisterConsumerHandler(consumer ConsumerHandler) error
	Run() error
}

type ConsumerHandler struct {
	Topic       string
	Channel     string
	Concurrent  int
	MaxAttempts int32
	MaxInFlight int
	Enable      bool
	URL         string

	extraConfig interface{}

	Handler Handler
}

func (c *ConsumerHandler) SetExtraConfig(extra interface{}) {
	c.extraConfig = extra
}

func (c ConsumerHandler) ParseExtraConfig(target interface{}) error {
	bytes, err := json.Marshal(c.extraConfig)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, target)
}

type Handler func(ctx context.Context, message IMessage) error
