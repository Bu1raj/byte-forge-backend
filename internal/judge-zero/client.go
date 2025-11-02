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
	body, _ := json.Marshal(req)
	fields := "stdout,stderr,time,memory,status,exit_code,exit_signal"
	url := fmt.Sprintf("http://172.19.6.225:2358/submissions?base64_encoded=false&wait=true&fields=%s", fields)
	res, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return models.SubmissionResponse{}, err
	}
	defer res.Body.Close()

	var response models.SubmissionResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return models.SubmissionResponse{}, err
	}

	fmt.Printf("Judge0 Response: %+v\n", response)

	return response, nil
}
