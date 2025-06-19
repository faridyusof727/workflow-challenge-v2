package workflow

import (
	"context"
	"fmt"
	"maps"
	"strconv"
	"time"
	"workflow-code-test/api/internal/edge"
	"workflow-code-test/api/internal/node"
	"workflow-code-test/api/pkg/helper"
	"workflow-code-test/api/pkg/nodes"
)

const (
	startNode = "start"
	endNode   = "end"
)

type ServiceImpl struct {
	repo        Repository
	nodeService *nodes.Service
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
	executionResult := &ExecutionResult{
		ExecutedAt: time.Now(),
		Steps:      make([]Step, 0),
	}

	wf, err := s.repo.WorkflowWithNodesAndEdges(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	input := executionInput.FormData

	executionResult.Steps = append(executionResult.Steps, Step{
		NodeID: startNode,
		Type:   startNode,
		Label:  startNode,
		Status: StepStatusCompleted,
	})
	in := inData{
		source:             startNode,
		sourceHandleResult: false,
	}
	out := outData{}
	for next(wf.Edges, wf.Nodes, &in, &out) {
		fmt.Printf("input: %+v\n", out)
		node := out.nextNode
		if node.Data.Metadata != nil {
			maps.Copy(input, node.Data.Metadata)
		}

		executor := s.nodeService.LoadNode(node.ID)
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

		outputFields := make([]string, 0)
		if outputVars, ok := node.Data.Metadata["outputVariables"].([]any); ok {
			for _, v := range outputVars {
				outputFields = append(outputFields, fmt.Sprintf("%v", v))
			}
			executor.SetOutputFields(outputFields)
		}

		output, err := executor.Execute(ctx)
		if err != nil {
			return nil, err
		}

		fmt.Printf("output: %+v\n\n\n\n\n\n\n", output)
		maps.Copy(input, output.(map[string]any))

		for key, val := range output.(map[string]any) {
			for _, f := range outputFields {
				if result, ok := val.(bool); ok {
					if key == f {
						in.sourceHandleResult = result
					}
				}
			}
		}

		in.source = node.ID
		executionResult.Steps = append(executionResult.Steps, Step{
			NodeID:      node.ID,
			Type:        node.Kind,
			Label:       node.Data.Label,
			Status:      StepStatusCompleted,
			Description: node.Data.Description,
			Output:      output.(map[string]any),
		})
	}

	return executionResult, nil
}

type outData struct {
	nextNode         node.Node
	needSourceHandle bool
}

type inData struct {
	source             string
	sourceHandleResult bool
}

func next(edges []edge.Edge, nodes []node.Node, in *inData, out *outData) bool {
	edge, found := helper.Find(edges, func(item edge.Edge) bool {
		// check here
		handle := false
		if item.SourceHandle != nil {
			n, err := strconv.ParseBool(*item.SourceHandle)
			if err != nil {
				// LOG
			}
			handle = n
		}

		return item.Source == in.source && in.sourceHandleResult == handle
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

	if out.nextNode.ID == endNode {
		return false
	}

	var needHandle bool
	if edge.SourceHandle != nil {
		n, err := strconv.ParseBool(*edge.SourceHandle)
		if err != nil {
			// LOG
		}
		needHandle = n
	}

	*out = outData{
		nextNode:         nxNode,
		needSourceHandle: needHandle,
	}

	return true
}

func NewService(repo Repository, nodeService *nodes.Service) Service {
	return &ServiceImpl{
		repo:        repo,
		nodeService: nodeService,
	}
}
