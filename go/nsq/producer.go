package nsq

import (
	"encoding/json"
	"strings"

	"github.com/nsqio/go-nsq"
)

// Publisher is struct for publisher
type Publisher struct {
	producer *nsq.Producer
	prefix   string
}

// NewPublisher will create new publisher instance
// leaf the prefix empty
func NewPublisher(publishAddress, prefix string) (*Publisher, error) {

	config := nsq.NewConfig()
	prod, err := nsq.NewProducer(publishAddress, config)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		producer: prod,
		prefix:   prefix,
	}, nil
}

// Publish will publish the data using json format, by default will always use the prefix in the topic
func (p *Publisher) Publish(topic string, data interface{}) error {
	topic = strings.Join([]string{p.prefix, topic}, "_")
	return p.PublishWithoutPrefix(topic, data)
}

// PublishWithoutPrefix will publish the data using json format without prefix in the topic
func (p *Publisher) PublishWithoutPrefix(topic string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return p.producer.Publish(topic, payload)
}
