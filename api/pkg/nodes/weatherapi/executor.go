package weatherapi

import (
	"context"
	"fmt"

	"workflow-code-test/api/pkg/nodes"
	"workflow-code-test/api/pkg/openstreetmap"
	"workflow-code-test/api/pkg/openweather"
)

type Options struct {
	// Input represents the city as now
	Input string
}

type Executor struct {
	opts          *Options
	geoClient     openstreetmap.Client
	weatherClient openweather.Client
}

// ID implements NodeExecutor.
func (c *Executor) ID() string {
	return "weather-api"
}

func NewExecutor(opts *Options, geoClient openstreetmap.Client, weatherClient openweather.Client) nodes.NodeExecutor {
	return &Executor{
		opts:          opts,
		geoClient:     geoClient,
		weatherClient: weatherClient,
	}
}

func (c *Executor) Execute(ctx context.Context) (map[string]interface{}, error) {
	lat, lng, err := c.geoClient.LatLngByCity(c.opts.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to get lat lng: %w", err)
	}

	temperature, err := c.weatherClient.TemperatureInCelsiusByLatLng(lat, lng)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather: %w", err)
	}

	return map[string]interface{}{
		"temperature": temperature,
	}, nil
}
