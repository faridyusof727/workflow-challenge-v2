-- +goose Up
-- +goose StatementBegin
CREATE TABLE nodes (
	id varchar NOT NULL,
    is_default boolean DEFAULT false NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	
	CONSTRAINT nodes_pkey PRIMARY KEY (id)
);

INSERT INTO nodes (id, is_default) VALUES ('start', true);
INSERT INTO nodes (id, is_default) VALUES ('end', true);
INSERT INTO nodes (id, is_default) VALUES ('form', false);
INSERT INTO nodes (id, is_default) VALUES ('weather-api', false);
INSERT INTO nodes (id, is_default) VALUES ('condition', false);
INSERT INTO nodes (id, is_default) VALUES ('email', false);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE nodes;
-- +goose StatementEnd
