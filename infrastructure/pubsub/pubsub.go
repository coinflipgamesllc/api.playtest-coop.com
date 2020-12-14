package pubsub

import (
	"sync"
)

type Message struct {
	Topic string
	Data  interface{}
}

type EventChan chan Message

type EventChans []EventChan

type EventBus struct {
	Subscribers map[string]EventChans
	lock        sync.RWMutex
}

var (
	Instance = &EventBus{Subscribers: map[string]EventChans{}}
)

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

func (b *EventBus) Subscribe(topic string, ch EventChan) {
	b.lock.Lock()

	if prev, found := b.Subscribers[topic]; found {
		b.Subscribers[topic] = append(prev, ch)
	} else {
		b.Subscribers[topic] = append([]EventChan{}, ch)
	}

	b.lock.Unlock()
}
