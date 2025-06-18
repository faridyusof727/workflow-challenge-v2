-- +goose Up
-- +goose StatementBegin
DO $$
DECLARE
    workflow_uuid UUID;
BEGIN
    INSERT INTO workflows (id, name, created_at, updated_at)
    VALUES (gen_random_uuid(), 'My Workflow', now(), now())
    RETURNING id INTO workflow_uuid;

    INSERT INTO workflow_nodes (workflow_id, node_id, kind, position_x, position_y, data_label, data_description, data_metadata, created_at, updated_at)
    SELECT workflow_uuid, vals.node_id, vals.kind, vals.position_x, vals.position_y, vals.data_label, vals.data_description, vals.data_metadata::jsonb, now(), now()
    FROM (VALUES
        ('start', 'start', -160, 300, 'Start', 'Begin weather check workflow', '{"hasHandles": {"source": true, "target": false}}'),
        ('form', 'form', 152, 304, 'User Input', 'Process collected data - name, email, location', '{"hasHandles": {"source": true, "target": true}, "inputFields": ["name", "email", "city"], "outputVariables": ["name", "email", "city"]}'),
        ('weather-api', 'integration', 460, 304, 'Weather API', 'Fetch current temperature for {{city}}', '{"hasHandles": {"source": true, "target": true}, "inputVariables": ["city"], "outputVariables": ["temperature"]}'),
        ('condition', 'condition', 794, 304, 'Check Condition', 'Evaluate temperature threshold', '{"hasHandles": {"source": ["true", "false"], "target": true}, "conditionExpression": "{{temperature}} {{operator}} {{threshold}}", "outputVariables": ["conditionMet"]}'),
        ('email', 'email', 1096, 88, 'Send Alert', 'Email weather alert notification', '{"hasHandles": {"source": true, "target": true}, "inputVariables": ["name", "city", "temperature"], "emailTemplate": {"subject": "Weather Alert", "body": "Weather alert for {{city}}! Temperature is {{temperature}}°C!"}, "outputVariables": ["emailSent"]}'),
        ('end', 'end', 1360, 302, 'End', 'Finish workflow', '{"hasHandles": {"source": false, "target": true}}')
    ) AS vals (
        node_id,
        kind,
        position_x,
        position_y,
        data_label,
        data_description,
        data_metadata
    );

    INSERT INTO workflow_edges (workflow_id, node_source, node_target, kind, is_animated, is_source_handle, label, label_style, style, created_at, updated_at)
    SELECT workflow_uuid, vals.node_source, vals.node_target, vals.kind::edge_kind, vals.is_animated, vals.is_source_handle, vals.label, vals.label_style::jsonb, vals.style::jsonb, now(), now()
    FROM (VALUES
        ('start', 'form', 'smoothstep', true, null, 'Initialize', null, '{"stroke": "#10b981", "strokeWidth": 3}'),
        ('form', 'weather-api', 'smoothstep', true, null, 'Submit Data', null, '{"stroke": "#3b82f6", "strokeWidth": 3}'),
        ('weather-api', 'condition', 'smoothstep', true, null, 'Temperature Data', null, '{"stroke": "#f97316", "strokeWidth": 3}'),
        ('condition', 'email', 'smoothstep', true, true, '✓ Condition Met', '{"fill": "#10b981", "fontWeight": "bold"}', '{"stroke": "#10b981", "strokeWidth": 3}'),
        ('condition', 'end', 'smoothstep', true, false, '✗ No Alert Needed', '{"fill": "#6b7280", "fontWeight": "bold"}', '{"stroke": "#6b7280", "strokeWidth": 3}'),
        ('email', 'end', 'smoothstep', true, false, 'Alert Sent', '{"fill": "#ef4444", "fontWeight": "bold"}', '{"stroke": "#ef4444", "strokeWidth": 2}')
    ) AS vals (
        node_source,
        node_target,
        kind,
        is_animated,
        is_source_handle,
        label,
        label_style,
        style
    );
END $$;  
-- +goose StatementEnd
