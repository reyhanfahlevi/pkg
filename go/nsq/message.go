package nsq

import (
	"time"

	"github.com/nsqio/go-nsq"
)

// IMessage interface define contract for nsq message in handler
type IMessage interface {
	// Finish sends a FIN command to the nsqd which
	// sent this message
	Finish()
	// RequeueWithoutBackoff sends a REQ command to the nsqd which
	// sent this message, using the supplied delay.
	//
	// Notably, using this method to respond does not trigger a backoff
	// event on the configured Delegate.
	RequeueWithoutBackoff(delay time.Duration)
	// Requeue sends a REQ command to the nsqd which
	// sent this message, using the supplied delay.
	Requeue(delay time.Duration)
	// GetAttempts will get the current attempts
	GetAttempts() uint16
	// GetBody will get the body value
	GetBody() []byte
}

// Message alias for built in nsq message
type Message struct {
	*nsq.Message
}

// GetAttempts return number of how many this message enter the consumer
func (m *Message) GetAttempts() uint16 {
	return m.Attempts
}

// GetBody return body of message
func (m *Message) GetBody() []byte {
	return m.Body
}
