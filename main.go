package main

import (
	"fmt"
	"log"
	"net/http"
)

type Topic struct {
	chans []chan interface{}
}

func (t *Topic) Subscribe(c chan interface{}) {

}

func (t *Topic) Unsubscribe(c chan interface{}) {

}

var topics = make(map[string]Topic)

func handler(rw http.ResponseWriter, req *http.Request) {
	topic := req.URL.Path // get topic name
	switch req.Method {

	case "GET":
		// save_topic(rw, topic, topics)
		t, ok := topics[topic]
		if !ok {
			t = Topic{chans: make([]chan interface{})}
			topics[topic] = t
		}

		c := make(chan interface{})
		t.Subscribe(c)
		v := <-c
		rw.Write(v)

	case "POST":
		rw.WriteHeader(201)
	}
	fmt.Fprintf(rw, topic)
}

// func save_topic(rw http.ResponseWriter, topic string, topics Topic) {
// 	fmt.Fprintf(rw, topic)

// }

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
