-- +goose Up
-- +goose StatementBegin
CREATE TABLE workflows (
	id uuid DEFAULT gen_random_uuid() NOT NULL,
    "name" varchar NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	
	CONSTRAINT workflows_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workflows;
-- +goose StatementEnd
