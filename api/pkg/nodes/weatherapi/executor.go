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

	args         map[string]any
	outputFields []string
}

func (e *Executor) SetArgs(args map[string]any) {
	e.args = args
}

func (e *Executor) SetOutputFields(fields []string) {
	e.outputFields = fields
}

// Validate implements nodes.NodeExecutor.
func (e *Executor) ValidateAndParse(argsCheck []string) error {
	for _, key := range argsCheck {
		_, ok := e.args[key].(string)
		if !ok {
			return fmt.Errorf("%s: validation key failed, key: %v", e.ID(), key)
		}
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

	// Hardcoded for now to explicitly there should be one output from the api request
	if len(e.outputFields) != 1 {
		return nil, fmt.Errorf("%s: output should only contain one variable, outputs: %+v", e.ID(), e.outputFields)
	}

	result := map[string]any{}
	for _, field := range e.outputFields {
		result[field] = fmt.Sprintf("%.2f", temperature)
	}

	return result, nil
}
