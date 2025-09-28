package main

import (
	"log"
	"net/http"

	"github.com/Bu1raj/byte-forge-backend/internal/api"
	"github.com/Bu1raj/byte-forge-backend/internal/queue"
)

//TODO need to store these in vault
const BROKER = "localhost:9092"
const TOPIC = "submissions"

func initializeKafkaProducer() *queue.Producer{
	return queue.NewProducer(BROKER, TOPIC)
}

func main() {
	http.Handle("/submit", &api.SubmitHandler{
		Producer: initializeKafkaProducer(),
	})
	// http.HandleFunc("/result/", api.ResultHandler)

	log.Println("listening :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
