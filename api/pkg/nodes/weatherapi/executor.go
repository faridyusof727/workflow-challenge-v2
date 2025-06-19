package weatherapi

import (
	"context"
	"fmt"

	"workflow-code-test/api/pkg/openstreetmap"
	"workflow-code-test/api/pkg/openweather"
)

type Options struct {
	GeoClient     openstreetmap.Client
	WeatherClient openweather.Client
}

type Executor struct {
	Opts *Options

	args   map[string]any
}

func (e *Executor) SetArgs(args map[string]any) {
	e.args = args
}

// Validate implements nodes.NodeExecutor.
func (e *Executor) ValidateAndParse() error {
	_, ok := e.args["city"].(string)
	if !ok {
		return fmt.Errorf("%s: validation failed to get city where it should string", e.ID())
	}

	return nil
}

// ID implements NodeExecutor.
func (e *Executor) ID() string {
	return "weather-api"
}

func (e *Executor) Execute(ctx context.Context) (any, error) {
	lat, lng, err := e.Opts.GeoClient.LatLngByCity(e.args["city"].(string))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get lat lng: %w", e.ID(), err)
	}

	temperature, err := e.Opts.WeatherClient.TemperatureInCelsiusByLatLng(lat, lng)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get weather: %w", e.ID(), err)
	}

	return map[string]any{
		"temperature": temperature,
	}, nil
}
