package di

import (
	"workflow-code-test/api/pkg/mailer"
	"workflow-code-test/api/pkg/nodes"
	"workflow-code-test/api/pkg/openstreetmap"
	"workflow-code-test/api/pkg/openweather"
)

// nodeService initializes and returns a new nodes.Service.
// For simplicity, OpenStreetMap, OpenWeather, and Mailer clients are
// initialized together here. In future iterations, these dependencies
// should be initialized individually to allow for more granular control
// and easier testing.
func (s *serviceImpl) nodeService() *nodes.Service {
	return nodes.NewService(openstreetmap.NewClient(), openweather.NewClient(), mailer.NewNoopClient())
}
