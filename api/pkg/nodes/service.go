package nodes

import (
	"workflow-code-test/api/pkg/nodes/condition"
	"workflow-code-test/api/pkg/nodes/types"
	"workflow-code-test/api/pkg/nodes/weatherapi"
)

// TODO: inject dependencies through di later
var nodeFactories = []types.NodeExecutor{
	&condition.Executor{},
	&weatherapi.Executor{},
}

func LoadNode(id string) types.NodeExecutor {
	for _, node := range nodeFactories {
		if node.ID() == id {
			return node
		}
	}

	return nil
}
