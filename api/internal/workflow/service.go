package workflow

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"strconv"
	"time"
	"workflow-code-test/api/internal/node"
	"workflow-code-test/api/pkg/nodes"
	"workflow-code-test/api/pkg/nodes/types"
)

const (
	startNode = "start"
	endNode   = "end"
)

type ServiceImpl struct {
	repo        Repository
	nodeService *nodes.Service
	log         *slog.Logger
}

// optimizedWorkflow contains pre-built indexes for lookups
type optimizedWorkflow struct {
	*Workflow
	edgesBySource map[string]map[bool]string // source -> handle -> target
	nodesById     map[string]node.Node       // nodeId -> node
}

type outData struct {
	nextNode node.Node
}

type inData struct {
	source             string
	sourceHandleResult bool
}

// Alias for better readability in refactored code
type executionState = inData

// Workflow implements Service.
func (s *ServiceImpl) Workflow(ctx context.Context, workflowID string) (*Workflow, error) {
	workflow, err := s.repo.WorkflowWithNodesAndEdges(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	return workflow, nil
}

func (s *ServiceImpl) Execute(ctx context.Context, workflowID string, executionInput *ExecutionInput) (*ExecutionResult, error) {
	wf, err := s.loadWorkflow(ctx, workflowID)
	if err != nil {
		return nil, err
	}

	return s.executeWorkflow(ctx, wf, executionInput)
}

func (s *ServiceImpl) loadWorkflow(ctx context.Context, workflowID string) (*Workflow, error) {
	return s.repo.WorkflowWithNodesAndEdges(ctx, workflowID)
}

func (s *ServiceImpl) executeWorkflow(ctx context.Context, wf *Workflow, executionInput *ExecutionInput) (*ExecutionResult, error) {
	executionResult := &ExecutionResult{
		ExecutedAt: time.Now(),
		Steps:      make([]Step, 0),
	}

	input := executionInput.FormData

	// Build optimized workflow structure for lookups
	optimizedWf := s.buildOptimizedWorkflow(wf)

	// Add start node to execution steps
	executionResult.Steps = append(executionResult.Steps, Step{
		NodeID: startNode,
		Type:   startNode,
		Label:  startNode,
		Status: StepStatusCompleted,
	})

	executionState := &executionState{
		source:             startNode,
		sourceHandleResult: false,
	}

	nextNodeData := outData{}
	for nextOptimized(optimizedWf, executionState, &nextNodeData) {
		step, err := s.executeNode(ctx, nextNodeData.nextNode, input)
		if err != nil {
			executionResult.Status = ExecutionStatusFailed
			return executionResult, err
		}

		// Update input with node output
		if step.Output != nil {
			maps.Copy(input, step.Output)
		}

		// Update execution state for next iteration
		s.updateExecutionState(step, executionState)

		executionResult.Steps = append(executionResult.Steps, *step)
	}

	// Add end node to execution steps
	executionResult.Steps = append(executionResult.Steps, Step{
		NodeID: endNode,
		Type:   endNode,
		Label:  endNode,
		Status: StepStatusCompleted,
	})
	executionResult.Status = ExecutionStatusCompleted
	return executionResult, nil
}

func (s *ServiceImpl) executeNode(ctx context.Context, node node.Node, input map[string]any) (*Step, error) {
	s.log.Info("starting node execution",
		slog.Any("node", node),
		slog.Any("input", input),
	)

	// Merge node metadata into input
	if node.Data.Metadata != nil {
		maps.Copy(input, node.Data.Metadata)
	}

	// Get and validate executor
	executor := s.nodeService.LoadNode(node.ID)
	if executor == nil {
		return nil, fmt.Errorf("executor not found with ID: %v", node.ID)
	}

	// Configure executor with input and validation
	if err := s.configureExecutor(executor, node, input); err != nil {
		return nil, fmt.Errorf("failed to configure executor for node %v: %w", node.ID, err)
	}

	// Execute node
	output, err := executor.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute node %v: %w", node.ID, err)
	}

	// Process output and create step
	outputMap := s.processNodeOutput(output)
	step := &Step{
		NodeID:      node.ID,
		Type:        node.Kind,
		Label:       node.Data.Label,
		Status:      StepStatusCompleted,
		Description: node.Data.Description,
		Output:      outputMap,
	}

	return step, nil
}

func (s *ServiceImpl) configureExecutor(executor types.NodeExecutor, node node.Node, input map[string]any) error {
	executor.SetArgs(input)

	// Handle input variables
	inputFields := s.extractInputFields(node.Data.Metadata)
	if err := executor.ValidateAndParse(inputFields); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Handle output variables
	outputFields := s.extractOutputFields(node.Data.Metadata)
	if len(outputFields) > 0 {
		executor.SetOutputFields(outputFields)
	}

	return nil
}

func (s *ServiceImpl) extractInputFields(metadata map[string]any) []string {
	if inputVars, ok := metadata["inputVariables"].([]any); ok {
		fields := make([]string, len(inputVars))
		for i, v := range inputVars {
			fields[i] = fmt.Sprintf("%v", v)
		}
		return fields
	}
	return nil
}

func (s *ServiceImpl) extractOutputFields(metadata map[string]any) []string {
	if outputVars, ok := metadata["outputVariables"].([]any); ok {
		fields := make([]string, 0, len(outputVars))
		for _, v := range outputVars {
			fields = append(fields, fmt.Sprintf("%v", v))
		}
		return fields
	}
	return nil
}

func (s *ServiceImpl) processNodeOutput(output any) map[string]any {
	if outputMap, ok := output.(map[string]any); ok {
		return outputMap
	}
	return nil
}

func (s *ServiceImpl) updateExecutionState(step *Step, state *executionState) {
	if step.Output != nil {
		// Look for boolean output fields that affect routing
		for _, val := range step.Output {
			if result, ok := val.(bool); ok {
				// Update state with boolean result for conditional routing
				state.sourceHandleResult = result
				break // Take first boolean for now
			}
		}
	}
	state.source = step.NodeID
}

// buildOptimizedWorkflow creates optimized data structures for lookups
func (s *ServiceImpl) buildOptimizedWorkflow(wf *Workflow) *optimizedWorkflow {
	optimized := &optimizedWorkflow{
		Workflow:      wf,
		edgesBySource: make(map[string]map[bool]string),
		nodesById:     make(map[string]node.Node),
	}

	// Build nodes index
	for _, n := range wf.Nodes {
		optimized.nodesById[n.ID] = n
	}

	// Build edges index
	for _, e := range wf.Edges {
		if optimized.edgesBySource[e.Source] == nil {
			optimized.edgesBySource[e.Source] = make(map[bool]string)
		}

		// Parse handle as boolean
		handle := false
		if e.SourceHandle != nil {
			if parsed, err := strconv.ParseBool(*e.SourceHandle); err == nil {
				handle = parsed
			}
		}

		optimized.edgesBySource[e.Source][handle] = e.Target
	}

	return optimized
}

// nextOptimized performs lookup using pre-built indexes
func nextOptimized(wf *optimizedWorkflow, in *inData, out *outData) bool {
	sourceEdges, sourceExists := wf.edgesBySource[in.source]
	if !sourceExists {
		return false
	}

	targetNodeID, edgeExists := sourceEdges[in.sourceHandleResult]
	if !edgeExists {
		return false
	}

	// Check if we've reached the end
	if targetNodeID == endNode {
		return false
	}

	// lookup for the actual node
	nextNode, nodeExists := wf.nodesById[targetNodeID]
	if !nodeExists {
		return false
	}

	*out = outData{
		nextNode: nextNode,
	}

	return true
}

func NewService(repo Repository, nodeService *nodes.Service, log *slog.Logger) Service {
	return &ServiceImpl{
		repo:        repo,
		nodeService: nodeService,
		log:         log,
	}
}
