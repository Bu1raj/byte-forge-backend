package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/Bu1raj/byte-forge-backend/internal/models"
	"github.com/Bu1raj/byte-forge-backend/internal/store"
	"github.com/Bu1raj/byte-forge-backend/pkg/utils"
)

// SubmitHandler handles code submission requests.
func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SubmitReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	// might need some work here, like max timeout
	if req.TimeoutSecond <= 0 {
		req.TimeoutSecond = 3
	}

	jobId := utils.NewID()
	payload := models.KafkaCodeSubmissionsPayload{
		ID:            jobId,
		SubmitRequest: req,
	}
	data, _ := json.Marshal(payload)
	submissions_producer, ok := store.GetProducer("submissions")
	if !ok {
		http.Error(w, "failed to get submissions producer", http.StatusInternalServerError)
		return
	}

	err = submissions_producer.SendMessage(data)
	if err != nil {
		http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": jobId})
}

// ResultHandler handles requests to fetch the result of a code execution.
func ResultHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	val, ok := store.GetResult(id)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(val)
}
