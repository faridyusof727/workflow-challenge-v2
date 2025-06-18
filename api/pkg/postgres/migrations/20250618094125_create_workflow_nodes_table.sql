-- +goose Up
-- +goose StatementBegin
CREATE TABLE workflow_nodes (
    workflow_id uuid NOT NULL REFERENCES workflows(id) ON DELETE CASCADE,
    node_id varchar NOT NULL,
    kind varchar NOT NULL,
    position_x integer NOT NULL,
    position_y integer NOT NULL,
    data_label varchar NOT NULL,
    data_description varchar NOT NULL,
    data_metadata jsonb,
    created_at timestamptz DEFAULT now() NOT NULL,
    updated_at timestamptz DEFAULT now() NOT NULL,
    
    CONSTRAINT workflow_nodes_pkey PRIMARY KEY (workflow_id, node_id),
    CONSTRAINT workflow_nodes_fk FOREIGN KEY (node_id) REFERENCES nodes(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workflow_nodes;
-- +goose StatementEnd
