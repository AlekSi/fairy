package fairy

import (
	"sync"
)

type Hub struct {
	topics  map[string]*Topic // key - topic
	bufSize int
	rw      sync.RWMutex
}

func NewHub(bufSize int) *Hub {
	return &Hub{
		topics:  make(map[string]*Topic),
		bufSize: bufSize,
	}
}

func (h *Hub) GetTopic(topic string) (t *Topic) {
	// optimistic path
	h.rw.RLock()
	t = h.topics[topic]
	h.rw.RUnlock()
	if t != nil {
		return
	}

	h.rw.Lock()
	defer h.rw.Unlock()

	// try again with write lock
	t = h.topics[topic]
	if t != nil {
		return
	}

	t = NewTopic(h.bufSize)
	h.topics[topic] = t
	return
}
