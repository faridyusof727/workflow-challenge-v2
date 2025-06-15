package workflow

import (
	"context"
	"net/http"
)

// Service defines the interface for workflow-related operations.
// It provides methods to retrieve workflow details and execute workflows.
type Service interface {
	// Workflow retrieves a workflow by its ID, including its associated nodes and edges.
	// It takes a context and workflow ID as parameters and returns the complete Workflow or an error.
	Workflow(ctx context.Context, workflowID string) (*Workflow, error)

	// Execute runs a workflow with the given ID and returns the execution result.
	// It takes a context and workflow ID as parameters and returns the execution result or an error.
	Execute(ctx context.Context, workflowID string) (*ExecutionResult, error)
}

// Repository is an interface that provides methods to retrieve workflow data,
// including associated nodes and edges, from a data source.
type Repository interface {
	// WorkflowWithNodesAndEdges retrieves a workflow by its ID along with its
	// associated nodes and edges. The context allows for cancellation or timeouts,
	// and the returned *Workflow may be nil if an error occurs.
	//
	// Returns:
	//   - A pointer to the Workflow object if found.
	//   - An error if the retrieval fails or the workflow does not exist.
	WorkflowWithNodesAndEdges(ctx context.Context, workflowID string) (*Workflow, error)
}

// Handler is an interface that defines HTTP handler functions for managing workflows
// and their execution in a web service.
type Handler interface {
	// Workflow handles HTTP requests for retrieving or managing workflow details.
	// It writes the response using the http.ResponseWriter and processes the request
	// using the *http.Request object.
	Workflow(w http.ResponseWriter, r *http.Request)

	// Execute handles HTTP requests to start or control the execution of a workflow.
	// It uses the http.ResponseWriter to send responses and the *http.Request to read
	// input parameters or payload.
	Execute(w http.ResponseWriter, r *http.Request)
}
