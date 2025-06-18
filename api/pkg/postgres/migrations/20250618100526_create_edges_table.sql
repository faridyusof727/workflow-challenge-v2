-- +goose Up
-- +goose StatementBegin
CREATE TYPE edge_kind AS ENUM ('smoothstep');

CREATE TABLE workflow_edges (
    workflow_id uuid NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    node_source varchar NOT NULL,
    node_target varchar NOT NULL,
    kind edge_kind NOT NULL,
    is_animated boolean DEFAULT false NOT NULL,
    is_source_handle boolean DEFAULT NULL,
    label varchar,
    label_style jsonb,
    style jsonb,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    
    CONSTRAINT workflow_edges_pkey PRIMARY KEY (workflow_id, node_source, node_target),
    CONSTRAINT workflow_edges_fk1 FOREIGN KEY (node_source) REFERENCES nodes(id),
    CONSTRAINT workflow_edges_fk2 FOREIGN KEY (node_target) REFERENCES nodes(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workflow_edges;
DROP TYPE edge_kind;
-- +goose StatementEnd
