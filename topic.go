package main

import (
	"github.com/fairy-project/fairy/common"
	"sync"
)

type Topic struct {
	subscribers map[string]chan common.Message
	rw          sync.RWMutex
}

func NewTopic() *Topic {
	return &Topic{subscribers: make(map[string]chan common.Message)}
}

func (t *Topic) GetChannel(subscriberId string) (c chan common.Message) {
	t.rw.Lock()
	defer t.rw.Unlock()

	c, present := t.subscribers[subscriberId]
	if present {
		return
	}

	c = make(chan common.Message)
	t.subscribers[subscriberId] = c
	return
}

func (t *Topic) Unsubscribe(subscriberId string) {
	t.rw.Lock()
	defer t.rw.Unlock()

	delete(t.subscribers, subscriberId)
}

func (t *Topic) Publish(m common.Message) {
	t.rw.Lock()
	defer t.rw.Unlock()

	for _, c := range t.subscribers {
		c <- m
	}
}
