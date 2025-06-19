package workflow

import (
	"context"
	"fmt"
	"time"
	"workflow-code-test/api/pkg/mailer"
	"workflow-code-test/api/pkg/nodes"
	"workflow-code-test/api/pkg/nodes/condition"
	"workflow-code-test/api/pkg/openstreetmap"
	"workflow-code-test/api/pkg/openweather"
)

type ServiceImpl struct {
	repo Repository
}

// Workflow implements Service.
func (s *ServiceImpl) Workflow(ctx context.Context, workflowID string) (*Workflow, error) {
	workflow, err := s.repo.WorkflowWithNodesAndEdges(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	return workflow, nil
}

func (s *ServiceImpl) Execute(ctx context.Context, workflowID string) (*ExecutionResult, error) {
	wf, err := s.repo.WorkflowWithNodesAndEdges(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	// result := &ExecutionResult{}
	nodeService := nodes.NewService(openstreetmap.NewClient(), openweather.NewClient(), mailer.NewNoopClient())

	input := map[string]any{
		"name":  "John Doe",
		"email": "example@example.com",
		"city":  "London",
	}

	for _, node := range wf.Nodes {
		if node.ID == "start" || node.ID == "end" {
			continue
		}
		if node.ID == "condition" {
			input["expression"] = node.Data.Metadata["conditionExpression"]
			input["threshold"] = float64(30)
			input["operator"] = condition.EqualToOperator
		}

		fmt.Printf("running node %+v\n", node.ID)
		fmt.Printf("setter %+v\n", input)
		executor := nodeService.LoadNode(node.ID)
		if executor == nil {
			fmt.Printf("executor with ID: %s not found", node.ID)
			continue
		}
		
		executor.SetArgs(input)
		err := executor.ValidateAndParse()
		if err != nil {
			return nil, err
		}

		output, err := executor.Execute(ctx)
		if err != nil {
			return nil, err
		}

		fmt.Printf("output: %+v\n", output)
		for key, o := range output.(map[string]any) {
			input[key] = o
		}
	}

	return &ExecutionResult{
		Status:     ExecutionStatusCompleted,
		ExecutedAt: time.Now(),
		Steps:      []Step{}, // TODO: set actual steps
	}, nil
}

func NewService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}
