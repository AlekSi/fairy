package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Message map[string]interface{}

type Topic struct {
	chans []chan Message
}

// Метод добавления топика в массив
func (t *Topic) Subscribe(c chan Message) {
	t.chans = append(t.chans, c)
}

// Метод удаления топика в массив
func (t *Topic) Unsubscribe(c chan Message) {

}

var topics = make(map[string]Topic)

func handler(rw http.ResponseWriter, req *http.Request) {
	topic := req.URL.Path // get topic name
	switch req.Method {

	case "GET":
		t, ok := topics[topic]
		if !ok {
			t = Topic{chans: make([]chan Message, 0)}
			topics[topic] = t
		}

		c := make(chan Message)
		t.Subscribe(c)
		message := <-c
		t.Unsubscribe(c)

		b, err := json.Marshal(message)
		if err == nil {
			rw.Write(b)
		}

	case "POST":
		body, err := ioutil.ReadAll(req.Body)
		if err == nil {
			fmt.Fprintf(rw, string(body))
		}
		//   []byte -> Message
		// for c := range t.chans {
		// 	c <- message
		// }
		rw.WriteHeader(201)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
