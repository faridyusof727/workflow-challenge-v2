package workflow

import (
	"time"
	"workflow-code-test/api/internal/edge"
	"workflow-code-test/api/internal/node"
)

// ExecutionStatus represents the status of an execution or step
type ExecutionStatus string

const (
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
)

// ExecutionResult represents the immediate result of starting a workflow execution
type ExecutionResult struct {
	Status     ExecutionStatus `json:"status"`
	ExecutedAt time.Time       `json:"executedAt"`
	Steps      []Step          `json:"steps"`
}

type Step struct {
	NodeID      string         `json:"nodeId"`
	Type        string         `json:"type"`
	Label       string         `json:"label"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Output      map[string]any `json:"output"`
}

type EmailDraft struct {
	To        string    `json:"to"`
	From      string    `json:"from"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	Timestamp time.Time `json:"timestamp"`
}

type Workflow struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Nodes     []node.Node `json:"nodes"`
	Edges     []edge.Edge `json:"edges"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}
