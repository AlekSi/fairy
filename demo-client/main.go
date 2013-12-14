package main

import (
	"bytes"
	"encoding/json"
	"github.com/manveru/faker"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func get(url string) {
	for {
		resp, err := http.Get(url)
		check(err)

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		check(err)

		var m map[string]string
		err = json.Unmarshal(body, &m)
		check(err)

		log.Println(m)
	}
}

func post(url string, fake *faker.Faker) {
	for {
		m := map[string]string{"m": fake.Sentence(1, true)}
		log.Printf("Posted: %v \n", m)

		b, err := json.Marshal(m)
		check(err)

		var buf bytes.Buffer
		buf.Write([]byte(b))

		resp, err := http.Post(url, "application/json", &buf)
		check(err)

		if resp.StatusCode != 201 {
			log.Fatal(resp.Status)
		}
		resp.Body.Close()

		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	url := "http://localhost:8081/topic"
	fake, err := faker.New("en")
	check(err)

	for i := 0; i < 1; i++ {
		go post(url, fake)
	}

	get(url)
}
