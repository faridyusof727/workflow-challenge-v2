# Nodes Package

The `nodes` package provides a pluggable system for executing different types of workflow nodes. Each node represents a
specific operation that can be executed as part of a workflow, supporting conditional branching, external API integration,
email notifications, and form data processing.

## Architecture

The package follows a standardized interface pattern where all nodes implement the `NodeExecutor` interface:

```go
type NodeExecutor interface {
    Execute(ctx context.Context) (any, error)           // Executes the node's primary logic
    ID() string                                         // Returns unique node identifier  
    SetArgs(args map[string]any)                        // Sets input arguments/parameters
    SetOutputFields(fields []string)                    // Specifies which output fields to return
    ValidateAndParse(argsCheck []string) error          // Validates and parses configuration
}
```

### Service Registry

The node service manages node factories and provides dependency injection:

```go
type Service struct {
    nodeFactories map[string]types.NodeExecutor
}

// Load a node by ID with all dependencies injected
func (s *Service) LoadNode(id string) types.NodeExecutor
```

## Available Nodes

### 1. Form Node (`form`)

**Purpose**: Entry point for workflow execution, captures and passes through user form data.

**Input Arguments**: Any key-value pairs from user form submission

**Output**: Fields specified in `outputFields` configuration, extracted from input arguments

**Dependencies**: None

**Example:**

```go
executor := service.LoadNode("form")
executor.SetArgs(map[string]any{
    "name": "John Doe",
    "email": "john@example.com", 
    "city": "London",
    "threshold": "25",
})
executor.SetOutputFields([]string{"name", "email", "city", "threshold"})
err := executor.ValidateAndParse([]string{"name", "email", "city", "threshold"})
result, err := executor.Execute(ctx)
// Returns: {"name": "John Doe", "email": "john@example.com", "city": "London", "threshold": "25"}
```

### 2. Weather API Node (`weather-api`)

**Purpose**: Retrieves current temperature for a specified city using geocoding and weather APIs.

**Input Arguments**:

- `city` (string): City name to get weather for

**Output**: Single temperature field (formatted as string with 2 decimal places)

**Dependencies**:

- `GeoClient`: OpenStreetMap client for geocoding
- `WeatherClient`: OpenWeather client for temperature data

**Process Flow**:

1. Convert city name to latitude/longitude coordinates
2. Retrieve current temperature using coordinates  
3. Format temperature as string with 2 decimal places
4. Return in specified output field

**Example:**

```go
executor := service.LoadNode("weather-api")
executor.SetArgs(map[string]any{
    "city": "London",
})
executor.SetOutputFields([]string{"temperature"})
err := executor.ValidateAndParse([]string{"city"})
result, err := executor.Execute(ctx)
// Returns: {"temperature": "22.50"}
```

### 3. Condition Node (`condition`)

**Purpose**: Evaluates conditional expressions using template placeholders and mathematical operators.

**Input Arguments**:

- `conditionExpression` (string): Template with placeholders (e.g., "{{temperature}} {{operator}} {{threshold}}")
- `operator` (string): One of `greater_than`, `less_than`, `equals`, `greater_than_or_equal`, `less_than_or_equal`
- Template variables: Any values referenced in the expression placeholders

**Output**: Single boolean field specified in `outputFields`

**Dependencies**: Uses expr-lang library for expression evaluation

**Template System**: Uses `{{key}}` syntax for placeholder replacement

**Example:**

```go
executor := service.LoadNode("condition")
executor.SetArgs(map[string]any{
    "conditionExpression": "{{temperature}} {{operator}} {{threshold}}",
    "operator": "greater_than",
    "temperature": "28.50",
    "threshold": "25",
})
executor.SetOutputFields([]string{"result"})
err := executor.ValidateAndParse([]string{"conditionExpression", "operator"})
result, err := executor.Execute(ctx)
// Returns: {"result": true} (if 28.50 > 25)
```

### 4. Email Node (`email`)

**Purpose**: Sends emails with template support for dynamic content generation.

**Input Arguments**:

- `email` (string): Recipient email address
- `emailTemplate` (map[string]any): Template object with `subject` and `body` fields
- Template variables: Values for placeholder replacement in email body

**Output**: Single boolean field indicating email sent status

**Dependencies**:

- `MailClient`: Email service client for sending emails

**Template System**: Uses `{{key}}` syntax for placeholder replacement in email body

**Example:**

```go
executor := service.LoadNode("email")
executor.SetArgs(map[string]any{
    "email": "john@example.com",
    "emailTemplate": map[string]any{
        "subject": "Weather Alert",
        "body": "Hello {{name}}, the temperature in {{city}} is {{temperature}}°C",
    },
    "name": "John Doe",
    "city": "London", 
    "temperature": "28.50",
})
executor.SetOutputFields([]string{"emailSent"})
err := executor.ValidateAndParse([]string{"email", "emailTemplate"})
result, err := executor.Execute(ctx)
// Returns: {"emailSent": true}
```

## Usage

### Service Initialization

Initialize the node service with required dependencies:

```go
import (
    "workflow-code-test/api/pkg/nodes"
    "workflow-code-test/api/pkg/openstreetmap"
    "workflow-code-test/api/pkg/openweather" 
    "workflow-code-test/api/pkg/mailer"
)

// Initialize external clients
geoClient := openstreetmap.NewClient()
weatherClient := openweather.NewClient(apiKey)
mailClient := mailer.NewClient(smtpConfig)

// Create node service with dependencies
nodeService := nodes.NewService(geoClient, weatherClient, mailClient)
```

### Loading Nodes

Load nodes by ID using the service registry:

```go
// Load a specific node by ID (returns nil if not found)
executor := nodeService.LoadNode("weather-api")
if executor == nil {
    return fmt.Errorf("node type not supported: weather-api")
}
```

### Standard Execution Pattern

All nodes follow the same execution pattern:

```go
// 1. Load the node
executor := nodeService.LoadNode("condition")
if executor == nil {
    return fmt.Errorf("node not found")
}

// 2. Set input arguments  
executor.SetArgs(map[string]any{
    "conditionExpression": "{{temperature}} > {{threshold}}",
    "operator": "greater_than",
    "temperature": "28.50",
    "threshold": "25",
})

// 3. Configure output fields
executor.SetOutputFields([]string{"result"})

// 4. Validate required input fields
requiredFields := []string{"conditionExpression", "operator"}
if err := executor.ValidateAndParse(requiredFields); err != nil {
    return fmt.Errorf("validation failed: %w", err)
}

// 5. Execute with context
result, err := executor.Execute(context.Background())
if err != nil {
    return fmt.Errorf("execution failed: %w", err)
}

// 6. Process result (returns map[string]any)
if resultMap, ok := result.(map[string]any); ok {
    conditionResult := resultMap["result"].(bool)
    fmt.Printf("Condition result: %v\n", conditionResult)
}
```

### Workflow Integration

Nodes are integrated into workflows through the workflow service:

```go
// In workflow execution
step, err := s.executeNode(ctx, nextNode, input)
if err != nil {
    return nil, fmt.Errorf("node execution failed: %w", err)
}

// Node output is merged back into workflow input for next nodes
if step.Output != nil {
    maps.Copy(input, step.Output)
}
```

## Adding a New Node

To add a new node type, follow these steps:

### 1. Create Node Package

Create a new directory under `pkg/nodes/` (e.g., `pkg/nodes/notification/`).

### 2. Implement Executor

Create `executor.go` implementing the `NodeExecutor` interface:

```go
package notification

import (
    "context"
    "fmt"
    "workflow-code-test/api/pkg/nodes/types"
)

type Options struct {
    // External dependencies (if needed)
    NotificationClient NotificationService
}

type Executor struct {
    Opts         *Options
    args         map[string]any
    outputFields []string
}

func (e *Executor) ID() string {
    return "notification"
}

func (e *Executor) SetArgs(args map[string]any) {
    e.args = args
}

func (e *Executor) SetOutputFields(fields []string) {
    e.outputFields = fields
}

func (e *Executor) ValidateAndParse(argsCheck []string) error {
    // Validate required input fields
    for _, key := range argsCheck {
        if _, exists := e.args[key]; !exists {
            return fmt.Errorf("%s: required field missing: %s", e.ID(), key)
        }
    }
    
    // Validate specific field types
    if message, ok := e.args["message"].(string); !ok || message == "" {
        return fmt.Errorf("%s: message must be a non-empty string", e.ID())
    }
    
    return nil
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
    // Extract inputs
    message := e.args["message"].(string)
    recipient := e.args["recipient"].(string)
    
    // Perform node logic
    err := e.Opts.NotificationClient.Send(ctx, recipient, message)
    if err != nil {
        return nil, fmt.Errorf("%s: failed to send notification: %w", e.ID(), err)
    }
    
    // Build output based on configured fields
    result := make(map[string]any)
    for _, field := range e.outputFields {
        switch field {
        case "success":
            result[field] = true
        case "timestamp":
            result[field] = time.Now().Format(time.RFC3339)
        default:
            // Include other values from args if requested
            if value, exists := e.args[field]; exists {
                result[field] = value
            }
        }
    }
    
    return result, nil
}
```

### 3. Add Tests

Create `executor_test.go` with comprehensive test coverage:

```go
package notification

import (
    "context"
    "testing"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
)

// Mock notification service
type MockNotificationService struct {
    mock.Mock
}

func (m *MockNotificationService) Send(ctx context.Context, recipient, message string) error {
    args := m.Called(ctx, recipient, message)
    return args.Error(0)
}

func TestExecutor_Execute(t *testing.T) {
    tests := []struct {
        name          string
        args          map[string]any
        outputFields  []string
        argsCheck     []string
        setupMock     func(*MockNotificationService)
        expectedError string
        expectedOutput map[string]any
    }{
        {
            name: "successful notification",
            args: map[string]any{
                "message":   "Test message",
                "recipient": "user@example.com",
            },
            outputFields: []string{"success"},
            argsCheck:    []string{"message", "recipient"},
            setupMock: func(m *MockNotificationService) {
                m.On("Send", mock.Anything, "user@example.com", "Test message").Return(nil)
            },
            expectedOutput: map[string]any{"success": true},
        },
        {
            name: "validation failure - missing message",
            args: map[string]any{
                "recipient": "user@example.com",
            },
            argsCheck:     []string{"message", "recipient"},
            expectedError: "notification: required field missing: message",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockService := &MockNotificationService{}
            if tt.setupMock != nil {
                tt.setupMock(mockService)
            }

            executor := &Executor{
                Opts: &Options{
                    NotificationClient: mockService,
                },
            }
            executor.SetArgs(tt.args)
            executor.SetOutputFields(tt.outputFields)

            err := executor.ValidateAndParse(tt.argsCheck)
            if tt.expectedError != "" {
                require.Error(t, err)
                require.Contains(t, err.Error(), tt.expectedError)
                return
            }
            require.NoError(t, err)

            result, err := executor.Execute(context.Background())
            require.NoError(t, err)
            
            if tt.expectedOutput != nil {
                require.Equal(t, tt.expectedOutput, result)
            }

            mockService.AssertExpectations(t)
        })
    }
}
```

### 4. Register Node in Service

Add your node to the node service factory in `service.go`:

```go
func NewService(geo openstreetmap.Client, weather openweather.Client, mail mailer.Client, notif notification.Service) *Service {
    return &Service{
        nodeFactories: map[string]types.NodeExecutor{
            "form":         &form.Executor{},
            "weather-api":  &weatherapi.Executor{Opts: &weatherapi.Options{GeoClient: geo, WeatherClient: weather}},
            "condition":    &condition.Executor{},
            "email":        &email.Executor{Opts: &email.Options{MailClient: mail}},
            "notification": &notification.Executor{Opts: &notification.Options{NotificationClient: notif}}, // Add here
        },
    }
}
```

### 5. Update Dependency Injection

If your node requires external services, update the service constructor and dependency injection:

```go
// In main.go or service initialization
notificationService := notification.NewService(config.NotificationConfig)
nodeService := nodes.NewService(geoClient, weatherClient, mailClient, notificationService)
```

## Best Practices

### 1. Input Validation

- Always validate required fields in `ValidateAndParse(argsCheck []string)`
- Validate field types with proper error messages
- Check for empty or invalid values

```go
func (e *Executor) ValidateAndParse(argsCheck []string) error {
    for _, key := range argsCheck {
        if _, exists := e.args[key]; !exists {
            return fmt.Errorf("%s: required field missing: %s", e.ID(), key)
        }
    }
    return nil
}
```

### 2. Error Handling

- Include node ID in all error messages for debugging
- Use error wrapping with `fmt.Errorf("%w", err)` for context
- Provide descriptive error messages

```go
return fmt.Errorf("%s: failed to process data: %w", e.ID(), err)
```

### 3. Output Field Configuration

- Build output dynamically based on `outputFields` configuration
- Allow flexible output field selection
- Provide default values for common fields

```go
result := make(map[string]any)
for _, field := range e.outputFields {
    switch field {
    case "success":
        result[field] = true
    case "data":
        result[field] = processedData
    }
}
```

### 4. Context Handling

- Always respect context cancellation in `Execute(ctx context.Context)`
- Pass context to external service calls
- Handle timeout scenarios gracefully

```go
func (e *Executor) Execute(ctx context.Context) (any, error) {
    select {
    case <-ctx.Done():
        return nil, fmt.Errorf("%s: execution cancelled: %w", e.ID(), ctx.Err())
    default:
        // Continue execution
    }
}
```

### 5. Testing Strategy

- Write table-driven tests covering all scenarios
- Test validation failures and edge cases
- Mock external dependencies
- Verify error messages and output formats

### 6. Dependencies

- Use dependency injection through Options struct
- Keep dependencies as interfaces for testability
- Initialize dependencies in service constructor

### 7. Template Support

- Use consistent `{{key}}` placeholder syntax
- Support template replacement in text fields
- Validate template variables against available data

## Error Handling Patterns

### Node-Specific Errors

```go
// Validation errors
return fmt.Errorf("%s: validation failed for field %s: expected %s, got %T", 
    e.ID(), fieldName, expectedType, actualValue)

// Execution errors  
return fmt.Errorf("%s: execution failed: %w", e.ID(), err)

// External service errors
return fmt.Errorf("%s: external service call failed: %w", e.ID(), err)
```

### Error Propagation

Errors are propagated up through the workflow execution engine, which:

1. Stops execution immediately on any node failure
2. Records the failed step in execution results
3. Returns partial execution results showing progress up to failure point

## Integration with Workflow Engine

### Node Configuration

Nodes receive configuration from workflow definitions:

```go
// Workflow service extracts from node metadata
inputFields := extractInputFields(node.Data.Metadata)
outputFields := extractOutputFields(node.Data.Metadata)

// Configure executor
executor.SetArgs(workflowInput)
executor.SetOutputFields(outputFields)
err := executor.ValidateAndParse(inputFields)
```

### Output Chaining

Node outputs are merged into the workflow input for subsequent nodes:

```go
// Execute node and get output
result, err := executor.Execute(ctx)

// Merge output into workflow input for next nodes
if step.Output != nil {
    maps.Copy(workflowInput, step.Output)
}
```

### Conditional Routing

Boolean node outputs control workflow branching:

```go
// Condition node returns boolean result
conditionResult := result.(map[string]any)["result"].(bool)

// Workflow engine uses boolean for routing decisions
if conditionResult {
    // Follow "true" path
} else {
    // Follow "false" path  
}
```

This architecture ensures consistent node behavior while providing flexibility for different node types and use cases.
