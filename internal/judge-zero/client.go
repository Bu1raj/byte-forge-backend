package judgezero

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Bu1raj/byte-forge-backend/internal/models"
)

// getJudge0Timeout returns the timeout duration for Judge0 API calls.
// It reads from the JUDGE0_TIMEOUT_SECONDS environment variable, defaulting to 30 seconds.
// Since Judge0 runs code synchronously (wait=true), the timeout should be set based on
// the maximum expected execution time for user code, plus network overhead.
// Recommended values: 30-60 seconds for typical code execution scenarios.
func getJudge0Timeout() time.Duration {
	timeoutStr := os.Getenv("JUDGE0_TIMEOUT_SECONDS")
	if timeoutStr == "" {
		return 30 * time.Second // Default: 30 seconds
	}
	
	timeoutSec, err := strconv.Atoi(timeoutStr)
	if err != nil || timeoutSec <= 0 {
		return 30 * time.Second // Fallback to default on invalid value
	}
	
	return time.Duration(timeoutSec) * time.Second
}

// SubmitCode submits the code to Judge0 and returns the submission response.
// The call to Judge0 is synchronous, waiting for the result.
// A timeout is configured to prevent indefinite hangs. The timeout can be adjusted
// via the JUDGE0_TIMEOUT_SECONDS environment variable (default: 30 seconds).
func SubmitCode(req models.SubmissionsRequest) (models.SubmissionResponse, error) {
	body, _ := json.Marshal(req)
	fields := "stdout,stderr,time,memory,status,exit_code,exit_signal"
	url := fmt.Sprintf("http://172.19.6.225:2358/submissions?base64_encoded=false&wait=true&fields=%s", fields)
	
	// Create HTTP client with timeout to prevent indefinite hangs
	client := &http.Client{
		Timeout: getJudge0Timeout(),
	}
	
	// Create request with context for additional control
	ctx := context.Background()
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return models.SubmissionResponse{}, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	
	res, err := client.Do(httpReq)
	if err != nil {
		return models.SubmissionResponse{}, fmt.Errorf("judge0 request failed: %w", err)
	}
	defer res.Body.Close()

	var response models.SubmissionResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return models.SubmissionResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Printf("Judge0 Response: %+v\n", response)

	return response, nil
}
