package api

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/Bu1raj/byte-forge-backend/internal/models"
	"github.com/Bu1raj/byte-forge-backend/pkg/utils"
)

// SubmitHandler handles code submission requests.
func (srv *Server) SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SubmissionsRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: need to look into timeout handling with judge0
	// if req.TimeoutSecond <= 0 {
	// 	req.TimeoutSecond = 3
	// }

	jobId := utils.NewID()
	payload := models.KafkaCodeSubmissionsPayload{
		ID:      jobId,
		Request: req,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "failed to marshal payload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	submissions_producer, ok := srv.Store.Kafka.GetProducer("submissions")
	if !ok {
		// TODO: initialize producer if not exists, and add to store, and send the message
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
func (srv *Server) ResultHandler(w http.ResponseWriter, r *http.Request) {
	var result models.SubmissionResponse
	id := filepath.Base(r.URL.Path)

	err := srv.Store.Redis.Get(r.Context(), id, &result)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
