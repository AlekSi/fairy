package main

import (
	"encoding/json"
	"github.com/fairy-project/fairy/fairy"
	"log"
	"net/http"
	"sync"
)

var (
	topics = make(map[string]*fairy.Topic)
	m      sync.Mutex
)

func publish(rw http.ResponseWriter, req *http.Request) {
	topic := req.URL.Path

	var msg fairy.Message
	err := json.NewDecoder(req.Body).Decode(&msg)
	req.Body.Close()
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	m.Lock()
	t, ok := topics[topic]
	if !ok {
		t = fairy.NewTopic()
		topics[topic] = t
	}
	m.Unlock()

	log.Printf("%s: publish to %s: %v", req.RemoteAddr, topic, msg)
	t.Publish(msg)

	rw.WriteHeader(201)
}

func subscribe(rw http.ResponseWriter, req *http.Request) {
	topic := req.URL.Path
	req.Body.Close()

	m.Lock()
	t, ok := topics[topic]
	if !ok {
		t = fairy.NewTopic()
		topics[topic] = t
	}
	m.Unlock()

	id := req.RemoteAddr
	log.Printf("%s: subscribe to %s", id, topic)
	c := t.GetChannel(id)
	msg := <-c
	log.Printf("%s: received %v", id, msg)

	json.NewEncoder(rw).Encode(msg)
}

func handler(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		subscribe(rw, req)
	case "POST":
		publish(rw, req)
	default:
		http.Error(rw, "bad method", 400)
	}
}

func main() {
	adminPrefix := "/meta/admin/"
	http.Handle(adminPrefix, http.StripPrefix(adminPrefix, http.FileServer(http.Dir("public"))))

	http.HandleFunc("/", handler)

	addr := ":8081"
	log.Printf("Listening on %s ...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
