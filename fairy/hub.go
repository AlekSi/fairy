package fairy

import (
	"sync"
)

type Hub struct {
	topics map[string]*Topic // key - topic
	rw     sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{topics: make(map[string]*Topic)}
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

	t = NewTopic(1000)
	h.topics[topic] = t
	return
}
