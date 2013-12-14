package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"github.com/fairy-project/fairy/fairy"
	"io"
	"log"
	"net/http"
)

const CookieName = "subscriber_id"

var hub = fairy.NewHub()

func publish(rw http.ResponseWriter, req *http.Request) {
	topic := req.URL.Path

	var msg fairy.Message
	err := json.NewDecoder(req.Body).Decode(&msg)
	req.Body.Close()
	if err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	t := hub.GetTopic(topic)
	// log.Printf("[ %20s ] %s: PUBLISH %v", req.RemoteAddr, topic, msg)
	t.Publish(msg)

	rw.WriteHeader(201)
}

func subscribe(rw http.ResponseWriter, req *http.Request) {
	topic := req.URL.Path
	req.Body.Close()

	cookie, _ := req.Cookie(CookieName)
	if cookie == nil {
		b := make([]byte, 8)
		_, err := io.ReadFull(rand.Reader, b)
		if err != nil {
			log.Fatal(err)
		}

		cookie = &http.Cookie{Name: CookieName, Value: hex.EncodeToString(b)}
		http.SetCookie(rw, cookie)
	}

	id := cookie.Value
	t := hub.GetTopic(topic)
	// log.Printf("[ %20s ] %s: %s: GET", req.RemoteAddr, id, topic)
	c := t.GetChannel(id)
	msg := <-c
	// log.Printf("[ %20s ] %s: %s: GOT %v", req.RemoteAddr, id, topic, msg)

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
