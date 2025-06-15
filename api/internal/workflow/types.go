package workflow

import (
	"time"
	"workflow-code-test/api/internal/edge"
	"workflow-code-test/api/internal/node"
)

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

type Workflow struct {
	ID        string      `json:"id"`
	Nodes     []node.Node `json:"nodes"`
	Edges     []edge.Edge `json:"edges"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}
