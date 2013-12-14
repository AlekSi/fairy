package fairy

import (
	"log"
	"sync"
	"time"
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
		t = NewTopic(1000)

		go func() {
			for {
				time.Sleep(time.Second)
				log.Printf("%s: %#v", topic, t)
			}
		}()

		h.topics[topic] = t
	}
	return
}
