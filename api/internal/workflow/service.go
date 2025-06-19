package workflow

import (
	"context"
	"fmt"
	"maps"
	"time"
	"workflow-code-test/api/internal/edge"
	"workflow-code-test/api/internal/node"
	"workflow-code-test/api/pkg/helper"
	"workflow-code-test/api/pkg/mailer"
	"workflow-code-test/api/pkg/nodes"
	"workflow-code-test/api/pkg/openstreetmap"
	"workflow-code-test/api/pkg/openweather"
)

const (
	startNode = "start"
	endNode   = "end"
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

func (s *ServiceImpl) Execute(ctx context.Context, workflowID string, executionInput *ExecutionInput) (*ExecutionResult, error) {
	wf, err := s.repo.WorkflowWithNodesAndEdges(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	// TODO: move to upper layer
	nodeService := nodes.NewService(openstreetmap.NewClient(), openweather.NewClient(), mailer.NewNoopClient())

	input := executionInput.FormData

	var node node.Node
	executionNode := startNode
	for next(wf.Edges, wf.Nodes, &executionNode, &node) {
		if node.Data.Metadata != nil {
			maps.Copy(input, node.Data.Metadata)
		}

		executor := nodeService.LoadNode(node.ID)
		if executor == nil {
			fmt.Printf("executor with ID: %s not found", node.ID)
			continue
		}

		fmt.Printf("input: %+v\n", input)
		executor.SetArgs(input)
		if inputVars, ok := node.Data.Metadata["inputVariables"].([]any); ok {
			fields := make([]string, len(inputVars))
			for i, v := range inputVars {
				fields[i] = fmt.Sprintf("%v", v)
			}

			err := executor.ValidateAndParse(fields)
			if err != nil {
				return nil, err
			}
		} else {
			err := executor.ValidateAndParse(nil)
			if err != nil {
				return nil, err
			}
		}

		if outputVars, ok := node.Data.Metadata["outputVariables"].([]any); ok {
			fields := make([]string, len(outputVars))
			for i, v := range outputVars {
				fields[i] = fmt.Sprintf("%v", v)
			}
			executor.SetOutputFields(fields)
		}

		output, err := executor.Execute(ctx)
		if err != nil {
			return nil, err
		}

		fmt.Printf("output: %+v\n\n\n\n\n\n\n", output)
		maps.Copy(input, output.(map[string]any))

		executionNode = node.ID
	}

	return &ExecutionResult{
		Status:     ExecutionStatusCompleted,
		ExecutedAt: time.Now(),
		Steps:      []Step{}, // TODO: set actual steps
	}, nil
}

func next(edges []edge.Edge, nodes []node.Node, source *string, nextNode *node.Node) bool {
	edge, found := helper.Find(edges, func(item edge.Edge) bool {
		return item.Source == *source
	})
	if !found {
		return false
	}

	nxNode, found := helper.Find(nodes, func(item node.Node) bool {
		return item.ID == edge.Target
	})
	if !found {
		return false
	}

	if nextNode.ID == endNode {
		return false
	}

	*nextNode = nxNode

	return true
}

func NewService(repo Repository) Service {
	return &ServiceImpl{repo: repo}
}
