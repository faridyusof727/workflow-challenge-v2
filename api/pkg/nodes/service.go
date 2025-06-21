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
	nodeFactories map[string]types.NodeExecutor
}

func NewService(
	geoClient openstreetmap.Client,
	weatherClient openweather.Client,
	mailClient mailer.Client,
) *Service {
	condition := types.NodeExecutor(&condition.Executor{})
	weatherapi := types.NodeExecutor(&weatherapi.Executor{
		Opts: &weatherapi.Options{
			GeoClient:     geoClient,
			WeatherClient: weatherClient,
		},
	})
	form := types.NodeExecutor(&form.Executor{})
	email := types.NodeExecutor(&email.Executor{
		Opts: &email.Options{
			MailClient: mailClient,
		},
	})

	nodeFactories := map[string]types.NodeExecutor{
		condition.ID():  condition,
		weatherapi.ID(): weatherapi,
		form.ID():       form,
		email.ID():      email,
	}

	return &Service{
		nodeFactories: nodeFactories,
	}
}

func (s *Service) LoadNode(id string) types.NodeExecutor {
	if node, ok := s.nodeFactories[id]; ok {
		return node
	}

	return nil
}
