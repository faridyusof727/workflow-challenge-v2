package workflow

import (
	"context"
	"fmt"
	"strconv"
	"workflow-code-test/api/internal/edge"
	"workflow-code-test/api/internal/node"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryImpl struct {
	pool *pgxpool.Pool
}

// WorkflowWithNodesAndEdges implements Repository.
func (r *RepositoryImpl) WorkflowWithNodesAndEdges(ctx context.Context, workflowID string) (*Workflow, error) {
	// Acquire a connection from the pool
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire database connection: %w", err)
	}
	defer conn.Release()

	args := pgx.NamedArgs{
		"workflowID": workflowID,
	}

	queryNodes := `select
			w.id,
			w.name,
			w.created_at,
			w.updated_at,
			wn.node_id,
			wn.kind,
			wn.position_x,
			wn.position_y,
			wn.data_label,
			wn.data_description,
			wn.data_metadata::jsonb,
			wn.created_at ,
			wn.updated_at 

		from
			workflows w
		join workflow_nodes wn on
			wn.workflow_id = w.id
		and w.id = @workflowID`

	rows, err := conn.Query(ctx, queryNodes, args)
	if err != nil {
		return nil, fmt.Errorf("failed to query nodes: %w", err)
	}
	defer rows.Close()

	workflow := Workflow{}

	for rows.Next() {
		var node node.Node

		err := rows.Scan(
			&workflow.ID,
			&workflow.Name,
			&workflow.CreatedAt,
			&workflow.UpdatedAt,
			&node.ID,
			&node.Kind,
			&node.Position.X,
			&node.Position.Y,
			&node.Data.Label,
			&node.Data.Description,
			&node.Data.Metadata,
			&node.CreatedAt,
			&node.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow: %w", err)
		}

		workflow.Nodes = append(workflow.Nodes, node)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("queryNodes: failed to iterate over rows: %w", err)
	}

	queryEdges := `select
			we.node_source,
			we.node_target,
			we.kind,
			we.is_animated,
			we.is_source_handle,
			we."style" ,
			we."label",
			we.label_style,
			we.created_at,
			we.updated_at
		from
			workflows w
		join workflow_edges we on
			we.workflow_id = w.id
		and w.id = @workflowID`

	rows, err = conn.Query(ctx, queryEdges, args)
	if err != nil {
		return nil, fmt.Errorf("failed to query edges: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var edge edge.Edge
		var isSourceHandle *bool

		err := rows.Scan(
			&edge.Source,
			&edge.Target,
			&edge.Kind,
			&edge.Animated,
			&isSourceHandle,
			&edge.Style,
			&edge.Label,
			&edge.LabelStyle,
			&edge.CreatedAt,
			&edge.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan workflow: %w", err)
		}

		// Convert nullable boolean to nullable string pointer
		if isSourceHandle != nil {
			val := strconv.FormatBool(*isSourceHandle)
			edge.SourceHandle = &val
		}

		edge.ID = fmt.Sprintf("%s-%s", edge.Source, edge.Target)
		workflow.Edges = append(workflow.Edges, edge)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("queryEdges: failed to iterate over rows: %w", err)
	}

	if workflow.ID == "" {
		return nil, fmt.Errorf("workflow not found")
	}

	return &workflow, nil
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &RepositoryImpl{
		pool: pool,
	}
}
