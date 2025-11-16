package models

type SubmissionsRequest struct {
	LanguageId string `json:"language_id"`
	SourceCode string `json:"source_code"`
	Stdin      string `json:"stdin,omitempty"`
}

type SubmissionResponse struct {
	Stdout     string `json:"stdout"`
	Stderr     string `json:"stderr"`
	Time       string `json:"time"`
	Memory     int    `json:"memory"`
	Status     Status `json:"status"`
	ExitCode   int    `json:"exit_code"`
	ExitSignal string `json:"exit_signal,omitempty"`
}

type Status struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

type KafkaCodeSubmissionsPayload struct {
	ID      string             `json:"id"`
	Request SubmissionsRequest `json:"submit_request"`
}

type KafkaCodeResultsPayload struct {
	ID     string             `json:"id"`
	Result SubmissionResponse `json:"result"`
}
