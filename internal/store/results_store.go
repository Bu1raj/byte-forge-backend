package store

import (
	"fmt"
	"sync"

	"github.com/Bu1raj/byte-forge-backend/internal/models"
)

var (
	results sync.Map // map[string]models.SubmissionResult
)

// StoreResult stores the result for a given job ID.
func StoreResult(id string, result interface{}) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if res, ok := result.(models.Result); ok {
		results.Store(id, res)
		return nil
	}
	return fmt.Errorf("invalid result type")
}

// GetResult retrieves the result for a given job ID.
// It returns the result and a boolean indicating if it exists.
func GetResult(id string) (models.Result, bool) {
	val, ok := results.Load(id)
	if !ok {
		return models.Result{}, false
	}
	if res, ok := val.(models.Result); ok {
		return res, true
	}
	return models.Result{}, false
}
