package wt

import "github.com/eclipse/paho.golang/paho"

const (
	SubscriberGroupTag = `group:"subscribers"`
)

// Subscriber is an interface for MQTT subscriber for WeTriage messages.
type Subscriber interface {
	// Option returns subscription options for MQTT. Mainly the subscription topic.
	Option() paho.SubscribeOptions
	// Handle is called when a message is received.
	Handle(publish *paho.Publish) error
}

// Payload replicates the WeTriage message envelope structure.
type Payload[T any] struct {
	Id        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
	Topic     string `json:"topic"`
	Content   T      `json:"content"`
}
