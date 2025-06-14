package openweather

type Temperature struct {
	Latitude             float64      `json:"latitude"`
	Longitude            float64      `json:"longitude"`
	GenerationtimeMS     float64      `json:"generationtime_ms"`
	UTCOffsetSeconds     int64        `json:"utc_offset_seconds"`
	Timezone             string       `json:"timezone"`
	TimezoneAbbreviation string       `json:"timezone_abbreviation"`
	Elevation            float64      `json:"elevation"`
	CurrentUnits         CurrentUnits `json:"current_units"`
	Current              Current      `json:"current"`
}

type Current struct {
	Time          string  `json:"time"`
	Interval      int64   `json:"interval"`
	Temperature2M float64 `json:"temperature_2m"`
}

type CurrentUnits struct {
	Time          string `json:"time"`
	Interval      string `json:"interval"`
	Temperature2M string `json:"temperature_2m"`
}
