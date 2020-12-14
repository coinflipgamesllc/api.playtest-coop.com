// Package pubsub is a very simple channel-based event system.
// TODO - consider using a messaging backend instead of in-memory channels?
package pubsub

import (
	"sync"
)

// Message is a wrapper around an event that has occurred with a name.
// e.g: User/Created, { "name": "New user", "email": "email@example.com", ... }
type Message struct {
	Topic string
	Data  interface{}
}

// EventChan is a channel for sending event messages
type EventChan chan Message

// EventChans is a slice of channels
type EventChans []EventChan

// EventBus handles the subscribers for the different events
type EventBus struct {
	Subscribers map[string]EventChans
	lock        sync.RWMutex
}

var (
	// Instance ensures singleton access to the event bus
	Instance = &EventBus{Subscribers: map[string]EventChans{}}
)

// Publish will push a message to all subscribers for a topic
func (b *EventBus) Publish(topic string, event interface{}) {
	b.lock.RLock()

	if chans, found := b.Subscribers[topic]; found {
		channels := append(EventChans{}, chans...)
		go func(msg Message, chans EventChans) {
			for _, ch := range chans {
				ch <- msg
			}
		}(Message{Topic: topic, Data: event}, channels)
	}

	b.lock.RUnlock()
}

// Subscribe will add the specified channel to the list of subscribers for a topic.
func (b *EventBus) Subscribe(topic string, ch EventChan) {
	b.lock.Lock()

	if prev, found := b.Subscribers[topic]; found {
		b.Subscribers[topic] = append(prev, ch)
	} else {
		b.Subscribers[topic] = append([]EventChan{}, ch)
	}

	b.lock.Unlock()
}
