package main

import (
	"encoding/json"
	"errors"
	"github.com/fairy-project/fairy/common"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var (
	ErrAlreadySubscribed = errors.New("already subscribed")
)

type Topic struct {
	subscribers map[string]chan common.Message
	m           sync.Mutex
}

func (t *Topic) Subscribe(id string, c chan common.Message) (err error) {
	t.m.Lock()
	defer t.m.Unlock()

	_, present := t.subscribers[id]
	if present {
		err = ErrAlreadySubscribed
		return
	}
	t.subscribers[id] = c
}

func (t *Topic) Unsubscribe(id string) {
	t.m.Lock()
	defer t.m.Unlock()

	delete(t.subscribers, id)
}

func (t *Topic) Publish(m common.Message) {
	for c := range t.subscribers {
		c <- m
	}
}
