package main

import (
	"fmt"
	"log"
	"net/http"
)

type Topic struct {
	chans []chan interface{}
}

// Метод добавления топика в массив
func (t *Topic) Subscribe(c chan interface{}) {
	t.chans = append(t.chans, c)
}

// Метод удаления топика в массив
func (t *Topic) Unsubscribe(c chan interface{}) {

}

var topics = make(map[string]Topic)

func handler(rw http.ResponseWriter, req *http.Request) {
	topic := req.URL.Path // get topic name
	switch req.Method {

	case "GET":
		t, ok := topics[topic]
		if !ok {
			t = Topic{chans: make([]chan interface{}, 0)}
			topics[topic] = t
		}

		c := make(chan interface{})
		t.Subscribe(c)
		// v := <-c
		// rw.Write(v)

	case "POST":
		rw.WriteHeader(201)
	}
	fmt.Fprintf(rw, topic)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
