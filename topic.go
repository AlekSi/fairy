package main

import (
	"github.com/fairy-project/fairy/common"
	"sync"
)

type Topic struct {
	subscribers map[string]chan common.Message
	m           sync.Mutex
}

func NewTopic() Topic {
	return Topic{subscribers: make(map[string]chan common.Message)}
}

func (t *Topic) GetChannel(subscriberId string) (c chan common.Message) {
	t.m.Lock()
	defer t.m.Unlock()

	c, present := t.subscribers[subscriberId]
	if present {
		return
	}

	c = make(chan common.Message)
	t.subscribers[subscriberId] = c
	return
}

func (t *Topic) Unsubscribe(subscriberId string) {
	t.m.Lock()
	defer t.m.Unlock()

	delete(t.subscribers, subscriberId)
}

func (t *Topic) Publish(m common.Message) {
	for _, c := range t.subscribers {
		c <- m
	}
}
