package fairy

import (
	"sync"
)

type Hub struct {
	topics map[string]*Topic // key - topic
	m      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{topics: make(map[string]*Topic)}
}

func (h *Hub) GetTopic(topic string) (t *Topic) {
	h.m.Lock()
	defer h.m.Unlock()

	t = h.topics[topic]
	if t == nil {
		t = NewTopic()
		h.topics[topic] = t
	}
	return
}
