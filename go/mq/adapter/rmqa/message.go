package rmqa

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	amqp.Delivery
	ch          *amqp.Channel
	topic       string
	attempts    int32
	maxAttempts int32
	requeued    bool
}

func (m *Message) Finish() {
	m.Ack(false)
}

func (m *Message) RequeueWithoutBackoff(delay time.Duration) {
	_ = m.republish(m.topic)
	m.Nack(false, false)
	m.requeued = true
}

func (m *Message) GetAttempts() int32 {
	if m.Headers == nil {
		m.Headers = make(amqp.Table)
		m.Headers["attempts"] = int32(0)
		return 0
	}

	if attempts, ok := m.Headers["attempts"]; ok {
		switch v := attempts.(type) {
		case int32:
			return v
		case int64:
			return int32(v)
		case int:
			return int32(v)
		default:
			return 0
		}
	}

	return 0
}

func (m *Message) Requeue(delay time.Duration) {
	attempts := m.GetAttempts()
	if attempts >= m.maxAttempts {
		m.Reject(false)
		return
	}

	_ = m.republish(m.topic)
	m.Nack(false, false)
	m.requeued = true
}

func (m *Message) republish(queueName string) error {

	// Modify headers
	newHeaders := m.Headers
	newHeaders["attempts"] = m.GetAttempts()

	return m.ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			Headers:       newHeaders,
			ContentType:   m.ContentType,
			Body:          m.Body,
			DeliveryMode:  m.DeliveryMode,
			CorrelationId: m.CorrelationId,
			ReplyTo:       m.ReplyTo,
			MessageId:     m.MessageId,
			Timestamp:     m.Timestamp,
			Type:          m.Type,
			AppId:         m.AppId,
		},
	)
}

func (m *Message) increaseAttempts() {
	attempts := m.GetAttempts()
	m.attempts = attempts + 1
	m.Headers["attempts"] = m.attempts
}

func (m *Message) GetBody() []byte {
	return m.Body
}
