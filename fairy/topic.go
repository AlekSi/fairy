package fairy

import (
	"fmt"
	"sync"
)

type Topic struct {
	subscribers map[string]chan Message // key - subscriberId
	rw          sync.RWMutex
}

// check interface
var _ fmt.GoStringer = &Topic{}

func NewTopic() *Topic {
	return &Topic{subscribers: make(map[string]chan Message)}
}

func (t *Topic) GoString() (res string) {
	t.rw.RLock()

	m := make(map[string][2]int, len(t.subscribers))
	for id, c := range t.subscribers {
		m[id] = [2]int{len(c), cap(c)}
	}

	t.rw.RUnlock()

	for id, pair := range m {
		res += fmt.Sprintf("%s(len=%d, cap=%d) ", id, pair[0], pair[1])
	}
	return
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
		select {
		case c <- m:
		default:
		}
	}
}
