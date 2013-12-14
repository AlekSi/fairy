package fairy

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Topic struct {
	subscribers  map[string]chan Message // key - subscriberId
	bufSize      int
	pub, pubSkip uint64
	rw           sync.Mutex
}

// check interface
var _ fmt.GoStringer = &Topic{}

func NewTopic(bufSize int) *Topic {
	return &Topic{subscribers: make(map[string]chan Message), bufSize: bufSize}
}

func (t *Topic) GoString() (res string) {
	res = fmt.Sprintf("pub=%d pubSkip=%d", atomic.LoadUint64(&t.pub), atomic.LoadUint64(&t.pubSkip))
	return
}

func (t *Topic) GetChannel(subscriberId string) (c chan Message) {
	t.rw.Lock()
	defer t.rw.Unlock()

	c, present := t.subscribers[subscriberId]
	if present {
		return
	}

	c = make(chan Message, t.bufSize)
	t.subscribers[subscriberId] = c
	return
}

func (t *Topic) Unsubscribe(subscriberId string) {
	t.rw.Lock()
	defer t.rw.Unlock()

	delete(t.subscribers, subscriberId)
}

func (t *Topic) Publish(m Message) {
	t.rw.Lock()
	defer t.rw.Unlock()

	for _, c := range t.subscribers {
		select {
		case c <- m:
			atomic.AddUint64(&t.pub, 1)
		default:
			atomic.AddUint64(&t.pubSkip, 1)
		}
	}
}
