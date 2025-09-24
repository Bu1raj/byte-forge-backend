package models

type SubmitReq struct {
	Language      string `json:"language"`
	Code          string `json:"code"`
	TimeoutSecond int    `json:"timeout_seconds"` // TODO:why do we need this?
}

type Result struct {
	ID        string `json:"id"`
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	ExitErr   string `json:"exit_error,omitempty"`
	TimedOut  bool   `json:"timed_out"`
	Duration  string `json:"duration"`
	Completed bool   `json:"completed"`
}
