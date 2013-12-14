package fairy

import (
	"sync"
)

type Message map[string]interface{}

type Topic struct {
	subscribers map[string]chan Message
	rw          sync.RWMutex
}

func NewTopic() *Topic {
	return &Topic{subscribers: make(map[string]chan Message)}
}

func (t *Topic) GetChannel(subscriberId string) (c chan Message) {
	t.rw.Lock()
	defer t.rw.Unlock()

	c, present := t.subscribers[subscriberId]
	if present {
		return
	}

	c = make(chan Message)
	t.subscribers[subscriberId] = c
	return
}

func (t *Topic) Unsubscribe(subscriberId string) {
	t.rw.Lock()
	defer t.rw.Unlock()

	delete(t.subscribers, subscriberId)
}

func (t *Topic) Publish(m Message) {
	t.rw.RLock()
	defer t.rw.RUnlock()

	for _, c := range t.subscribers {
		c <- m
	}
}
