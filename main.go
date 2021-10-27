package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/robfig/cron"
)

type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	b := payload{}
	json.Unmarshal(body, &b)
	w.Header().Set("Content-Type", "application/json")
	code := getRandomCode()
	now := time.Now()
	fmt.Printf("%v | %v			%v\n", now.Unix(), code, b.Uuid)
	w.WriteHeader(code)
}

func getRandomCode() int {
	val := rand.Intn(4)
	code := http.StatusOK
	if val == 1 {
		code = http.StatusInternalServerError
	}
	return code
}

func main() {
	fmt.Println("Running...")
	go cronProducer()
	s := &server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":3070", nil))
}

func cronProducer() {
	c := cron.New()
	c.AddFunc("@every 1s", produceToKafka)
	c.Start()
}

type redpandaPayload struct {
	Records []record `json:"records"`
}

type record struct {
	Value interface{} `json:"value"`
}

type payload struct {
	Uuid string `json:"uuid"`
}

func createPayload(id uuid.UUID) []byte {
	r := payload{
		Uuid: id.String(),
	}
	payload := redpandaPayload{Records: []record{{
		Value: r,
	}}}
	postBody, _ := json.Marshal(payload)
	return postBody
}

func produceToKafka() {
	id := uuid.New()
	payload := createPayload(id)

	responseBody := bytes.NewBuffer(payload)
	http.Post("http://localhost:8082/topics/test_benthos_perf", "application/vnd.kafka.json.v2+json", responseBody)
	// fmt.Println("produced to kafka", resp.StatusCode)
}
