package weatherapi

import (
	"context"
	"fmt"

	"workflow-code-test/api/pkg/nodes/types"
	"workflow-code-test/api/pkg/openstreetmap"
	"workflow-code-test/api/pkg/openweather"
)

type Inputs struct {
	City string
}

type Outputs struct {
	Temperature float64 `json:"temperature"`
}

type Options struct {
	Inputs        *Inputs
	GeoClient     openstreetmap.Client
	WeatherClient openweather.Client
}

type Executor struct {
	opts *Options
}

// Validate implements nodes.NodeExecutor.
func (e *Executor) Validate() error {
	panic("unimplemented")
}

// ID implements NodeExecutor.
func (e *Executor) ID() string {
	return "weather-api"
}

func NewExecutor(opts *Options) types.NodeExecutor {
	return &Executor{
		opts: opts,
	}
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
	lat, lng, err := e.opts.GeoClient.LatLngByCity(e.opts.Inputs.City)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get lat lng: %w", e.ID(), err)
	}

	temperature, err := e.opts.WeatherClient.TemperatureInCelsiusByLatLng(lat, lng)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get weather: %w", e.ID(), err)
	}

	return &Outputs{
		Temperature: temperature,
	}, nil
}
