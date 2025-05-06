package rmqa

import (
	"github.com/pkg/errors"
	"github.com/reyhanfahlevi/pkg/go/log"
	"github.com/reyhanfahlevi/pkg/go/mq/adapter"
)

type ConsumerManager struct {
	consumers []*Consumer
}

func NewManager() *ConsumerManager {
	return &ConsumerManager{}
}

func (c *ConsumerManager) RegisterConsumerHandler(consumer adapter.ConsumerHandler) error {
	conn := NewConsumer(consumer)
	c.consumers = append(c.consumers, conn)
	return nil
}

func (c *ConsumerManager) Run() error {
	for _, conn := range c.consumers {
		err := conn.Run()
		if err != nil {
			log.Error(errors.Wrapf(err, "%s failed to start", conn.handler.Topic))
			continue
		}
	}

	return nil
}
