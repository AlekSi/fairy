package main

import (
	"bytes"
	"encoding/json"
	"github.com/manveru/faker"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

const (
	C  = 10
	N1 = 100
)

var (
	RecvNum int64
	SentNum int64
)

func get(url string) {
	for {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		var msg map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&msg)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()

		atomic.AddInt64(&RecvNum, 1)
		log.Print(msg)
	}
}

func post(url string) {
	fake, err := faker.New("en")
	if err != nil {
		log.Fatal(err)
	}

	var (
		msg  = make(map[string]interface{})
		buf  bytes.Buffer
		resp *http.Response
	)

	for n1 := 0; n1 < N1; n1++ {
		msg["n"] = atomic.AddInt64(&SentNum, 1)
		msg["m"] = fake.Sentence(1, true)
		msg["ts"] = time.Now().Format("15:04:05.000")

		buf.Reset()
		err = json.NewEncoder(&buf).Encode(msg)
		if err != nil {
			log.Fatal(err)
		}

		resp, err = http.Post(url, "application/json", &buf)
		if err != nil {
			log.Fatal(err)
		}
		resp.Body.Close()

		if resp.StatusCode != 201 {
			log.Fatal(resp.Status)
		}
	}
}

func main() {
	url := "http://localhost:8081/topic"
	log.Printf("Sending %d messages to %s with concurrency %d ...", N1*C, url, C)

	for i := 0; i < C; i++ {
		go post(url)
	}

	get(url)
}
