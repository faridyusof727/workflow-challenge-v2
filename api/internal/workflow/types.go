package workflow

import "time"

// ExecutionResult represents the immediate result of starting a workflow execution
type ExecutionResult struct {
	Status     ExecutionStatus `json:"status"`
	ExecutedAt time.Time       `json:"executedAt"`
	// TODO: Add node steps, output and status
	Steps []any `json:"steps"`
}

// ExecutionStatus represents the status of an execution or step
type ExecutionStatus string

const (
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
)
