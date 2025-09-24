package api

import (
	"github.com/Bu1raj/byte-forge-backend/internal/executor"
	"github.com/Bu1raj/byte-forge-backend/internal/models"
	"github.com/Bu1raj/byte-forge-backend/pkg/utils"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// we might want to use something else here later
var results sync.Map // map[string]models.Result

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SubmitReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	// might need some work here, like max timeout
	if req.TimeoutSecond <= 0 {
		req.TimeoutSecond = 3
	}

	id := utils.NewID()
	tmpDir, err := os.MkdirTemp("", "submission-*") // can use the id as prefix?
	if err != nil {
		http.Error(w, "server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filename := executor.FilenameForLang(req.Language)
	if filename == "" {
		http.Error(w, "unsupported language", http.StatusBadRequest)
		return
	}
	if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(req.Code), fs.FileMode(0644)); err != nil {
		http.Error(w, "write error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	results.Store(id, models.Result{ID: id, Completed: false})

	go func(subID, dir, lang string, timeoutSec int) {
		res := executor.RunSubmission(subID, dir, lang, time.Duration(timeoutSec)*time.Second)
		results.Store(subID, res)
		_ = os.RemoveAll(dir)
	}(id, tmpDir, req.Language, req.TimeoutSecond)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func ResultHandler(w http.ResponseWriter, r *http.Request) {
	id := filepath.Base(r.URL.Path)
	val, ok := results.Load(id)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(val)
}
