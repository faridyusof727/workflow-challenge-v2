package openstreetmap_test

import (
	"testing"
	"workflow-code-test/api/pkg/openstreetmap"

	"github.com/stretchr/testify/require"
)

func TestLatLngByCity(t *testing.T) {
	tests := []struct {
		name    string
		city    string
		wantLat float64
		wantLng float64
		wantErr bool
	}{
		{
			name:    "KL",
			city:    "Kuala Lumpur",
			wantLat: 3.1516964,
			wantLng: 101.6942371,
			wantErr: false,
		},
		{
			name:    "Sydney",
			city:    "Sydney",
			wantLat: -33.8698439,
			wantLng: 151.2082848,
			wantErr: false,
		},
		{
			name:    "gibberish",
			city:    "asdbasdbasd",
			wantLat: 0,
			wantLng: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLat, gotLng, err := openstreetmap.NewClient().LatLngByCity(tt.city)

			if tt.wantErr {
				require.Error(t, err, "error on test case %s", tt.name)
				return
			}

			require.NoError(t, err, "error on test case %s", tt.name)
			require.Equal(t, tt.wantLat, gotLat, "lat on test case %s", tt.name)
			require.Equal(t, tt.wantLng, gotLng, "lng on test case %s", tt.name)
		})
	}
}
