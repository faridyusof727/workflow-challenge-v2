package weatherapi

import (
	"context"
	"fmt"

	"workflow-code-test/api/pkg/nodes"
	"workflow-code-test/api/pkg/openstreetmap"
	"workflow-code-test/api/pkg/openweather"
)

type Inputs struct {
	City string
}

type Outputs struct {
	Temperature float64 `json:"temperature"`
}

type WeatherExecutor = nodes.NodeExecutor[Outputs]

type Executor struct {
	inputs        *Inputs
	geoClient     openstreetmap.Client
	weatherClient openweather.Client
}

// Validate implements nodes.NodeExecutor.
func (e *Executor) Validate() error {
	panic("unimplemented")
}

// ID implements NodeExecutor.
func (e *Executor) ID() string {
	return "weather-api"
}

func NewExecutor(inputs *Inputs, geoClient openstreetmap.Client, weatherClient openweather.Client) WeatherExecutor {
	return &Executor{
		inputs:        inputs,
		geoClient:     geoClient,
		weatherClient: weatherClient,
	}
}

func (e *Executor) Execute(ctx context.Context) (*Outputs, error) {
	lat, lng, err := e.geoClient.LatLngByCity(e.inputs.City)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get lat lng: %w", e.ID(), err)
	}

	temperature, err := e.weatherClient.TemperatureInCelsiusByLatLng(lat, lng)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get weather: %w", e.ID(), err)
	}

	return &Outputs{
		Temperature: temperature,
	}, nil
}
