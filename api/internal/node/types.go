package node

import "time"

type Node struct {
	ID         string    `json:"id"`
	WorkflowID string    `json:"workflowId"`
	Kind       string    `json:"type"`
	Position   Position  `json:"position"`
	Data       Data      `json:"data"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Data struct {
	Label       string         `json:"label"`
	Description string         `json:"description"`
	Metadata    map[string]any `json:"metadata"`
}
