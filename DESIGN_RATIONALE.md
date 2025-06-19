# Design Rationale: Workflow Execution Engine

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Key Design Decisions](#key-design-decisions)
3. [Trade-offs & Assumptions](#trade-offs--assumptions)
4. [Extension Points](#extension-points)
5. [Performance Considerations](#performance-considerations)
6. [Testing Strategy](#testing-strategy)

## Architecture Overview

### High-Level Architecture Pattern

The system follows **Clean Architecture** principles with clear separation of concerns:

* **HTTP Layer** (`/api`)
* **Domain Layer** (`/internal/*`)
  * Following Hexagonal Pattern (port-adapter)
  * Port
    * `port.go`
  * Adapter
    * `handler.go`
    * `repository.go`
  * Core Login
    * `service.go`
    * `types.go`
* **Entrypoints** (`/cmd`)
  * The use of the `/cmd` directory in Go projects is a convention , not a requirement — but it's widely adopted and idiomatic, especially for larger or production-ready applications.
  * Using `/cmd` is an idiomatic and scalable way to organize Go applications. It keeps your main packages isolated, improves modularity, and makes your project easier to understand and grow.
* **Reusable Package**
  * Most of the pkgs are reusable. But interestingly, I would like to explain about `di`, `config`, `nodes`.
  * `di` - I've created a Dependency Injection (DI) package that centralizes the management of our application's dependencies. Instead of using global variable like the initial approach, I believe it's better to use DI pattern as it tends to provide more benefits:
    * It improves testability
      * If you can refer to `api/pkg/mailer/noop.go`, this is just an example of creating `noop`, `mock` services to assist with testing.
    * Enables loose coupling
    * Makes dependencies explicit
    * Supports clean architecture
    * Allows for better lifecycle management
  * `config` - The initial approach is using `os.Getenv` to get environment variables. While `os.Getenv` is great for getting started, I believe parsing envs to a dedicated struct would help us to with defaults, types, and easy to test.
  * `nodes` - The nodes package provides a pluggable system for executing different types of workflow nodes. See [api/pkg/nodes/README.md](api/pkg/nodes/README.md)

> The file structure is flatter, as initial development followed the Go idiom of 'flat until it hurts'—starting simple and only introducing hierarchy when necessary.
>
> Note: The phrase "flat until it hurts" is a well-known idiom in the Go community. It refers to the practice of keeping directory structures shallow and simple early in a project's life cycle. The idea is to avoid premature abstraction or over-engineering by not introducing unnecessary package divisions or deep nesting.
> Instead, you keep most code in one or a few packages (often main or a single top-level package), and only split things into sub-packages when the structure becomes unwieldy—when it "hurts" to keep it flat anymore.

**Benefits:**

* **Testability**: Each layer can be tested independently
* **Maintainability**: Clear responsibility boundaries
* **Extensibility**: Easy to add new features without affecting existing code

## Key Design Decisions

### 1. Node Executor Architecture

#### Decision: Plugin-Based Strategy Pattern

```go
type NodeExecutor interface {
    Execute(ctx context.Context) (any, error)
    SetArgs(args map[string]any)
    ValidateAndParse(inputFields []string) error
    SetOutputFields(outputFields []string)
}
```

**Rationale:**

* **Open/Closed Principle**: Easy to add new node types without modifying existing code
* **Separation of Concerns**: Each executor handles one specific node type
* **Testability**: Each executor can be tested independently
* **Dependency Injection**: External services (weather API, email) injected at runtime

### 2. Workflow Execution Engine

#### Decision: Graph Traversal with State Management

```go
type executionState struct {
    source             string
    sourceHandleResult bool
}
```

**Rationale:**

* **Conditional Logic**: Support for if/else branching based on node outputs
* **State Preservation**: Maintains execution context between nodes
* **Error Handling**: Can stop execution at any point and return partial results

## Trade-offs & Assumptions

### Database Trade-offs

#### Trade-off: Normalized vs Denormalized Schema

* **Chosen:** Normalized 4-table design
* **Benefits:** Data integrity, no duplication, flexible queries
* **Costs:** More complex joins, slightly slower for simple reads
* **Assumption:** Workflow complexity will grow, making normalized design worthwhile

#### Trade-off: JSONB vs Strongly Typed Columns

* **Chosen:** JSONB for node metadata, typed columns for core fields
* **Benefits:** Schema flexibility without sacrificing performance
* **Costs:** Less database-level validation
* **Assumption:** Node types will evolve frequently, requiring schema flexibility

### Execution Engine Trade-offs

#### Trade-off: Synchronous vs Asynchronous Execution

* **Chosen:** Synchronous execution
* **Benefits:** Simpler error handling, immediate results
* **Costs:** Cannot handle long-running workflows efficiently
* **Assumption:** Current workflows are short-duration (weather checks)

#### Trade-off: Single-threaded vs Parallel Node Execution

* **Chosen:** Single-threaded sequential execution
* **Benefits:** Simpler state management, easier debugging
* **Costs:** Cannot leverage parallelism for independent nodes
* **Assumption:** Node dependencies require sequential execution

### Error Handling Trade-offs

#### Trade-off: Fail-fast vs Continue-on-error

* **Chosen:** Fail-fast (stop execution on first error)
* **Benefits:** Clear failure points, prevents cascade failures
* **Costs:** Cannot complete partial workflows
* **Assumption:** All nodes in workflow are required for meaningful result

## Extension Points

### 1. New Node Types

Adding a new node type requires only:

```go
// 1. Create new executor
type NewNodeExecutor struct {
    // node-specific fields
}

func (e *NewNodeExecutor) Execute(ctx context.Context) (any, error) {
    // implementation
}

// 2. Register in node service
func (s *Service) LoadNode(nodeType string) types.NodeExecutor {
    switch nodeType {
    case "new-node-type":
        return &NewNodeExecutor{}
    // existing cases...
    }
}
```

**No changes needed to:**

* Database schema
* Workflow execution engine
* API handlers
* Core business logic

### 2. Complex Conditional Logic

Current system supports boolean conditions. Extensions could include:

```go
// Support for complex expressions
type AdvancedConditionExecutor struct {
    ExpressionEngine ExpressionEvaluator
}

// Support for multiple output paths
type MultiPathNode struct {
    OutputPaths map[string]bool
}
```

### 3. External Service Integration

Pattern for adding new external services:

```go
// New service interface
type NotificationService interface {
    Send(ctx context.Context, payload any) error
}

// Inject into executors
type SlackExecutor struct {
    slack NotificationService
}
```

### 4. Workflow Validation

Extension point for complex workflow validation:

```go
type WorkflowValidator interface {
    ValidateStructure(workflow *Workflow) error
    ValidateNodes(nodes []node.Node) error
    ValidateEdges(edges []edge.Edge) error
}
```

### 5. Monitoring & Observability

Built-in extension points for monitoring:

```go
// Execution middleware
type ExecutionMiddleware interface {
    BeforeExecution(ctx context.Context, step *Step) error
    AfterExecution(ctx context.Context, step *Step, result any, err error)
}

// Metrics collection
type MetricsCollector interface {
    RecordExecution(workflowID string, duration time.Duration, success bool)
    RecordNodeExecution(nodeType string, duration time.Duration)
}
```

## Performance Considerations

### Database Performance

1. **Indexing Strategy:**

   ```sql
   CREATE INDEX idx_workflows_id ON workflows(id);
   CREATE INDEX idx_workflow_nodes_workflow_id ON workflow_nodes(workflow_id);
   CREATE INDEX idx_workflow_edges_source ON workflow_edges(source);
   ```

2. **Query Optimization:**
   * Single query to load complete workflow with JOINs
   * JSONB indexing for metadata queries when needed

### Memory Management

1. **Execution Context:**
   * Minimal state preservation during execution
   * Immediate cleanup after workflow completion

2. **Node Output Chaining:**
   * Efficient map copying using `maps.Copy()`
   * No unnecessary data duplication

### API Performance

1. **Response Optimization:**
   * Structured JSON responses
   * Minimal data transfer
   * Proper HTTP status codes

## Testing Strategy

### Unit Testing

* **Node Executors:** Table-driven tests for all scenarios
* **Service Layer:** Mock repositories for isolated testing
* **Validation Logic:** Comprehensive input validation tests

### Integration Testing

* **Database Operations:** Real PostgreSQL instance
* **External APIs:** Mock weather service for predictable tests
* **Workflow Execution:** End-to-end execution scenarios

### Test Structure Example

```go
func TestWeatherAPIExecutor_Execute(t *testing.T) {
    tests := []struct {
        name        string
        city        string
        mockTemp    float64
        expectError bool
    }{
        {"valid city", "London", 25.5, false},
        {"invalid city", "InvalidCity", 0, true},
    }
    // table-driven test implementation
}
```

## Future Considerations

### Scalability Improvements

1. **Async Execution:** Message queue for long-running workflows
2. **Horizontal Scaling:** Stateless design enables easy scaling
3. **Caching:** Redis for frequently accessed workflow definitions

### Advanced Features

1. **Workflow Versioning:** Track changes to workflow definitions
2. **Execution History:** Optional persistence of execution logs
3. **Scheduling:** Cron-like scheduling for recurring workflows
4. **Visual Debugging:** Enhanced execution tracing and visualization
