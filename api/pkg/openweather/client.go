package openweather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Impl struct{}

// TemperatureInCelsiusByLatLng implements Client.
func (i *Impl) TemperatureInCelsiusByLatLng(lat float64, lng float64) (float64, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m", lat, lng))
	if err != nil {
		return 0, fmt.Errorf("failed to get temperature resp: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read resp body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get temperature resp with status: %s, body: %s", resp.Status, string(resBody))
	}

	var temperature Temperature
	if err := json.Unmarshal(resBody, &temperature); err != nil {
		return 0, fmt.Errorf("failed to unmarshal resp body: %w", err)
	}

	return float64(temperature.Current.Temperature2M), nil
}

func NewClient() Client {
	return &Impl{}
}
