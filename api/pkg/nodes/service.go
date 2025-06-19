package nodes

import (
	"workflow-code-test/api/pkg/mailer"
	"workflow-code-test/api/pkg/nodes/condition"
	"workflow-code-test/api/pkg/nodes/email"
	"workflow-code-test/api/pkg/nodes/form"
	"workflow-code-test/api/pkg/nodes/types"
	"workflow-code-test/api/pkg/nodes/weatherapi"
	"workflow-code-test/api/pkg/openstreetmap"
	"workflow-code-test/api/pkg/openweather"
)

type Service struct {
	nodeFactories []types.NodeExecutor
}

func NewService(
	geoClient openstreetmap.Client,
	weatherClient openweather.Client,
	mailClient mailer.Client,
) *Service {
	return &Service{
		nodeFactories: []types.NodeExecutor{
			&condition.Executor{},
			&weatherapi.Executor{
				Opts: &weatherapi.Options{
					GeoClient:     geoClient,
					WeatherClient: weatherClient,
				},
			},
			&form.Executor{},
			&email.Executor{
				Opts: &email.Options{
					MailClient: mailClient,
				},
			},
		},
	}
}

func (s *Service) LoadNode(id string) types.NodeExecutor {
	for _, node := range s.nodeFactories {
		if node.ID() == id {
			return node
		}
	}

	return nil
}
