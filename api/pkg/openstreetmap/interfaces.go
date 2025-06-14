package openstreetmap

type Client interface {
	// LatLngByCity retrieves the latitude and longitude coordinates for a given city.
	// It returns the latitude, longitude, and an error if the city cannot be located.
	LatLngByCity(city string) (float64, float64, error)
}
