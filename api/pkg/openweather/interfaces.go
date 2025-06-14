package openweather

type Client interface {
	// TemperatureInCelsiusByLatLng retrieves the temperature in Celsius for a given latitude and longitude.
	// It returns the temperature as a float64 and an error if the temperature cannot be retrieved.
	TemperatureInCelsiusByLatLng(lat, lng float64) (float64, error)
}
