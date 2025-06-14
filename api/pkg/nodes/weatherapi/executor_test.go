package weatherapi_test

import (
	"context"
	"errors"
	"testing"
	"workflow-code-test/api/pkg/nodes/weatherapi"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockGeoClient struct {
	mock.Mock
}

func (m *MockGeoClient) LatLngByCity(city string) (float64, float64, error) {
	args := m.Called(city)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

type MockWeatherClient struct {
	mock.Mock
}

func (m *MockWeatherClient) TemperatureInCelsiusByLatLng(lat, lng float64) (float64, error) {
	args := m.Called(lat, lng)
	return args.Get(0).(float64), args.Error(1)
}

func TestExecutor_ID(t *testing.T) {
	mockGeoClient := &MockGeoClient{}
	mockWeatherClient := &MockWeatherClient{}
	opts := &weatherapi.Options{Input: "test-city"}

	executor := weatherapi.NewExecutor(opts, mockGeoClient, mockWeatherClient)

	require.Equal(t, "weather-api", executor.ID())
}

func TestExecutor_Execute(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedLat      float64
		expectedLng      float64
		expectedTemp     float64
		geoError         error
		weatherError     error
		expectedOutput   map[string]interface{}
		expectedErrorMsg string
	}{
		{
			name:         "successful execution",
			input:        "Kuala Lumpur",
			expectedLat:  3.1516964,
			expectedLng:  101.6942371,
			expectedTemp: 28.5,
			expectedOutput: map[string]interface{}{
				"temperature": 28.5,
			},
		},
		{
			name:         "successful execution - Sydney",
			input:        "Sydney",
			expectedLat:  -33.8698439,
			expectedLng:  151.2082848,
			expectedTemp: 22.3,
			expectedOutput: map[string]interface{}{
				"temperature": 22.3,
			},
		},
		{
			name:             "geo client error",
			input:            "Invalid City",
			expectedLat:      0.0,
			expectedLng:      0.0,
			geoError:         errors.New("city not found"),
			expectedErrorMsg: "failed to get lat lng: city not found",
		},
		{
			name:             "weather client error",
			input:            "Valid City",
			expectedLat:      40.7128,
			expectedLng:      -74.0060,
			weatherError:     errors.New("weather service unavailable"),
			expectedErrorMsg: "failed to get weather: weather service unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGeoClient := &MockGeoClient{}
			mockWeatherClient := &MockWeatherClient{}

			opts := &weatherapi.Options{Input: tt.input}
			executor := weatherapi.NewExecutor(opts, mockGeoClient, mockWeatherClient)

			mockGeoClient.On("LatLngByCity", tt.input).Return(tt.expectedLat, tt.expectedLng, tt.geoError)

			if tt.geoError == nil {
				mockWeatherClient.On("TemperatureInCelsiusByLatLng", tt.expectedLat, tt.expectedLng).Return(tt.expectedTemp, tt.weatherError)
			}

			ctx := context.Background()
			outputs, err := executor.Execute(ctx)

			if tt.expectedErrorMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErrorMsg)
				require.Nil(t, outputs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedOutput, outputs)
			}

			mockGeoClient.AssertExpectations(t)
			if tt.geoError == nil {
				mockWeatherClient.AssertExpectations(t)
			}
		})
	}
}

func TestNewExecutor(t *testing.T) {
	mockGeoClient := &MockGeoClient{}
	mockWeatherClient := &MockWeatherClient{}
	opts := &weatherapi.Options{Input: "test-city"}

	executor := weatherapi.NewExecutor(opts, mockGeoClient, mockWeatherClient)

	require.NotNil(t, executor)
	require.Equal(t, "weather-api", executor.ID())
}
