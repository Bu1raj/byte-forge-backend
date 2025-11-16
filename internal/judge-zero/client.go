package judgezero

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bu1raj/byte-forge-backend/internal/models"
)

// SubmitCode submits the code to Judge0 and returns the submission response.
// The call to Judge0 is synchronous, waiting for the result.
func SubmitCode(req models.SubmissionsRequest) (models.SubmissionResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return models.SubmissionResponse{}, err
	}

	fields := "stdout,stderr,time,memory,status,exit_code,exit_signal"
	url := fmt.Sprintf("http://172.19.6.225:2358/submissions?base64_encoded=false&wait=true&fields=%s", fields)
	
	// TODO: this timeout value might need to be configurable
	client := &http.Client{
		Timeout:   5 * 60 * 1e9, // 5 minutes
	}
	res, err := client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return models.SubmissionResponse{}, err
	}
	defer res.Body.Close()

	var response models.SubmissionResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return models.SubmissionResponse{}, err
	}

	return response, nil
}
