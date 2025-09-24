package main

import (
	"log"
	"net/http"

	"github.com/Bu1raj/byte-forge-backend/internal/api"
)

func main() {
	http.HandleFunc("/submit", api.SubmitHandler)
	http.HandleFunc("/result/", api.ResultHandler)

	log.Println("listening :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
