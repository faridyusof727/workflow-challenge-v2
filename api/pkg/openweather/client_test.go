package openweather_test

import (
	"testing"
	"workflow-code-test/api/pkg/openweather"

	"github.com/stretchr/testify/require"
)

func TestTemperatureInCelsiusByLatLng(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lng     float64
		maxTemp float64
		minTemp float64
		wantErr bool
	}{
		{
			name:    "KL coordinates",
			lat:     3.1516964,
			lng:     101.6942371,
			maxTemp: 50.0,
			minTemp: 10.0,
			wantErr: false,
		},
		{
			name:    "Sydney coordinates",
			lat:     -33.8698439,
			lng:     151.2082848,
			maxTemp: 50.0,
			minTemp: -8,
			wantErr: false,
		},
		{
			name:    "Valid coordinates - New York",
			lat:     40.7128,
			lng:     -74.0060,
			maxTemp: 45,
			minTemp: -15,
			wantErr: false,
		},
		{
			name:    "Invalid coordinates - out of range lat",
			lat:     91.0,
			lng:     0.0,
			wantErr: true,
		},
		{
			name:    "Invalid coordinates - out of range lng",
			lat:     0.0,
			lng:     181.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTemp, err := openweather.NewClient().TemperatureInCelsiusByLatLng(tt.lat, tt.lng)

			if tt.wantErr {
				require.Error(t, err, "error on test case %s", tt.name)
				return
			}

			require.NoError(t, err, "error on test case %s", tt.name)
			require.NotZero(t, gotTemp, "temperature should not be zero on test case %s", tt.name)
			require.GreaterOrEqual(t, gotTemp, tt.minTemp, "temperature should be greater or equal to %f on test case %s", tt.minTemp, tt.name)
			require.LessOrEqual(t, gotTemp, tt.maxTemp, "temperature should be less or equal to %f on test case %s", tt.maxTemp, tt.name)
		})
	}
}
