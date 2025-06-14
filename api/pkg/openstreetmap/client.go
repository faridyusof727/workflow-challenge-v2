package openstreetmap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"workflow-code-test/api/pkg/helper"
)

type Impl struct{}

// LatLngByCity implements Client.
func (i *Impl) LatLngByCity(city string) (float64, float64, error) {
	cityUrlEncoded := url.QueryEscape(city)
	resp, err := http.Get(fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json", cityUrlEncoded))
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get city resp: %w", err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to read resp body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("failed to get city resp with status: %s, body: %s", resp.Status, string(resBody))
	}

	var cities []City
	if err := json.Unmarshal(resBody, &cities); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal resp body: %w", err)
	}

	matchedCity, found := helper.Find(cities, func(city City) bool {
		return city.Addresstype == "city"
	})
	if !found {
		return 0, 0, fmt.Errorf("failed to find city")
	}

	latitude, err := strconv.ParseFloat(matchedCity.Lat, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse latitude: %w", err)
	}

	longitude, err := strconv.ParseFloat(matchedCity.Lon, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse longitude: %w", err)
	}

	return latitude, longitude, nil
}

func NewClient() Client {
	return &Impl{}
}
