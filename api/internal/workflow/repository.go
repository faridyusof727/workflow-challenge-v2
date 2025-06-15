package workflow

import (
	"context"
	"workflow-code-test/api/internal/edge"
	"workflow-code-test/api/internal/node"
)

type RepositoryImpl struct {
}

// WorkflowWithNodesAndEdges implements Repository.
func (r *RepositoryImpl) WorkflowWithNodesAndEdges(ctx context.Context, workflowID string) (*Workflow, error) {
	// TODO: Placeholder implementation
	return &Workflow{
		ID: "550e8400-e29b-41d4-a716-446655440000",
		Nodes: []node.Node{
			{
				ID:         "start",
				WorkflowID: "550e8400-e29b-41d4-a716-446655440000",
				Kind:       "start",
				Position: node.Position{
					X: -160,
					Y: 300,
				},
				Data: node.Data{
					Label:       "Start",
					Description: "Begin weather check workflow",
					Metadata: map[string]any{
						"hasHandles": map[string]bool{
							"source": true,
							"target": false,
						},
					},
				},
			},
			{
				ID:         "form",
				WorkflowID: "550e8400-e29b-41d4-a716-446655440000",
				Kind:       "form",
				Position: node.Position{
					X: 152,
					Y: 304,
				},
				Data: node.Data{
					Label:       "User Input",
					Description: "Process collected data - name, email, location",
					Metadata: map[string]any{
						"hasHandles": map[string]bool{
							"source": true,
							"target": false,
						},
						"inputFields":     []string{"name", "email", "city"},
						"outputVariables": []string{"name", "email", "city"},
					},
				},
			},
			{
				ID:         "weather-api",
				WorkflowID: "550e8400-e29b-41d4-a716-446655440000",
				Kind:       "integration",
				Position: node.Position{
					X: 152,
					Y: 304,
				},
				Data: node.Data{
					Label:       "Weather API",
					Description: "Fetch weather data for the provided city",
					Metadata: map[string]any{
						"hasHandles": map[string]bool{
							"source": true,
							"target": false,
						},
						"inputFields":     []string{"city"},
						"outputVariables": []string{"temperature"},
					},
				},
			},
		},
		Edges: []edge.Edge{
			{
				ID:       "e1",
				Source:   "start",
				Target:   "form",
				Kind:     "smoothstep",
				Animated: true,
				Label:    "Start",
				Style: map[string]any{
					"stroke":      "#10b981",
					"strokeWidth": 3,
				},
			},
		},
	}, nil
}

func NewRepository() Repository {
	return &RepositoryImpl{}
}
