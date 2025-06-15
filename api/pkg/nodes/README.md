# Nodes Package

The `nodes` package provides a pluggable system for executing different types of workflow nodes. Each node represents a specific operation that can be executed as part of a workflow.

## Architecture

The package follows a standardized interface pattern where all nodes implement the `NodeExecutor` interface:

```go
type NodeExecutor interface {
    ID() string
    SetArgs(args map[string]any)
    ValidateAndParse() error
    Execute(ctx context.Context) (any, error)
}
```

## Available Nodes

### Condition Node (`condition`)

Evaluates conditional expressions with configurable operators.

**Input Arguments:**

- `expression` (string): Template expression with placeholders
- `threshold` (float64): Comparison threshold value
- `operator` (Operator): Comparison operator (`greater_than`, `less_than`, `equal_to`, `less_than_or_equal`, `greater_than_or_equal`)
- `temperature` (float64): Temperature value to compare

**Example:**

```go
executor := &condition.Executor{}
executor.SetArgs(map[string]any{
    "expression": "{{temperature}} {{operator}} {{threshold}}",
    "threshold": 25.0,
    "operator": condition.GreaterThanOperator,
    "temperature": 28.5,
})
err := executor.ValidateAndParse()
result, err := executor.Execute(ctx)
```

### Weather API Node (`weather-api`)

Retrieves temperature data for a given city using geocoding and weather APIs.

**Input Arguments:**

- `city` (string): City name to get weather for

**Dependencies:**

- `GeoClient`: OpenStreetMap client for geocoding
- `WeatherClient`: OpenWeather client for temperature data

**Example:**

```go
executor := &weatherapi.Executor{
    Opts: &weatherapi.Options{
        GeoClient: geoClient,
        WeatherClient: weatherClient,
    },
}
executor.SetArgs(map[string]any{
    "city": "London",
})
err := executor.ValidateAndParse()
result, err := executor.Execute(ctx)
```

## Usage

### Loading Nodes

Use the service registry to load nodes by ID:

```go
import "workflow-code-test/api/pkg/nodes"

// Load a specific node by ID
executor := nodes.LoadNode("condition")
if executor == nil {
    // Node not found
}
```

### Executing Nodes

Standard execution pattern:

```go
// 1. Load the node
executor := nodes.LoadNode("weather-api")

// 2. Set arguments
executor.SetArgs(map[string]any{
    "city": "New York",
})

// 3. Validate and parse inputs
if err := executor.ValidateAndParse(); err != nil {
    // Handle validation error
}

// 4. Execute
result, err := executor.Execute(context.Background())
if err != nil {
    // Handle execution error
}
```

## Adding a New Node

To add a new node type, follow these steps:

### 1. Create Node Package

Create a new directory under `pkg/nodes/` (e.g., `pkg/nodes/mynewnode/`).

### 2. Define Types

Create `types.go` with input/output structures:

```go
package mynewnode

type Inputs struct {
    Field1 string  `json:"field1"`
    Field2 int     `json:"field2"`
}

type Outputs struct {
    Result string `json:"result"`
}
```

### 3. Implement Executor

Create `executor.go`:

```go
package mynewnode

import (
    "context"
    "fmt"
)

type Executor struct {
    args   map[string]any
    inputs Inputs
}

func (e *Executor) ID() string {
    return "my-new-node"
}

func (e *Executor) SetArgs(args map[string]any) {
    e.args = args
}

func (e *Executor) ValidateAndParse() error {
    field1, ok := e.args["field1"].(string)
    if !ok {
        return fmt.Errorf("%s: validation failed to get field1 where it should be string", e.ID())
    }

    field2, ok := e.args["field2"].(int)
    if !ok {
        return fmt.Errorf("%s: validation failed to get field2 where it should be int", e.ID())
    }

    e.inputs = Inputs{
        Field1: field1,
        Field2: field2,
    }

    return nil
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
    // Implement your node logic here
    result := fmt.Sprintf("Processed: %s with %d", e.inputs.Field1, e.inputs.Field2)
    
    return Outputs{
        Result: result,
    }, nil
}
```

### 4. Add Tests

Create `executor_test.go`:

```go
package mynewnode

import (
    "context"
    "testing"
    "github.com/stretchr/testify/require"
)

func TestExecutor_Execute(t *testing.T) {
    executor := &Executor{}
    executor.SetArgs(map[string]any{
        "field1": "test",
        "field2": 42,
    })

    err := executor.ValidateAndParse()
    require.NoError(t, err)

    result, err := executor.Execute(context.Background())
    require.NoError(t, err)
    require.NotNil(t, result)

    outputs, ok := result.(Outputs)
    require.True(t, ok)
    require.Equal(t, "Processed: test with 42", outputs.Result)
}
```

### 5. Register Node

Add your node to the registry in `service.go`:

```go
import (
    "workflow-code-test/api/pkg/nodes/mynewnode"
    // ... other imports
)

var nodeFactories = []types.NodeExecutor{
    &condition.Executor{},
    &weatherapi.Executor{},
    &mynewnode.Executor{}, // Add your node here
}
```

### 6. Update Dependencies

If your node requires external dependencies, add them to the `Options` struct pattern:

```go
type Options struct {
    ExternalClient SomeClient
}

type Executor struct {
    Opts   *Options
    args   map[string]any
    inputs Inputs
}
```

## Best Practices

1. **Validation**: Always validate input types and required fields in `ValidateAndParse()`
2. **Error Handling**: Include the node ID in error messages for debugging
3. **Testing**: Write comprehensive tests covering success and failure scenarios
4. **Documentation**: Use JSON tags for input structs to match argument keys
5. **Context**: Always respect the context parameter for cancellation
6. **Type Safety**: Use type assertions with proper error handling

## Error Handling

All nodes should return descriptive errors with the node ID:

```go
return fmt.Errorf("%s: specific error description: %w", e.ID(), err)
```

This ensures errors can be traced back to their source node in complex workflows.
