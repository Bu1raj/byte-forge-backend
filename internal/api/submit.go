package api

import (
	"encoding/json"
	"net/http"

	"github.com/Bu1raj/byte-forge-backend/internal/models"
	"github.com/Bu1raj/byte-forge-backend/internal/queue"
	"github.com/Bu1raj/byte-forge-backend/pkg/utils"
)

type SubmitHandler struct {
	Producer *queue.Producer
}

// SubmitHandler handles code submission requests.
func (h *SubmitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	err = h.Producer.SendMessage(data)
	if err != nil {
		http.Error(w, "failed to enqueue job", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": jobId})
}

// TODO we would need a consumer here which can listen the results
// and update a persistent store or in-memory store, hence commenting this out for now

// ResultHandler handles requests to fetch the result of a code execution.
// func ResultHandler(w http.ResponseWriter, r *http.Request) {
// 	id := filepath.Base(r.URL.Path)
// 	val, ok := results.Load(id)
// 	if !ok {
// 		http.Error(w, "not found", http.StatusNotFound)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(val)
// }
