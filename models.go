package cars

// ADSBInput holds all parameters needed to calculate altitudes from an ADS-B source.
// HObsQne is the observed pressure altitude in meters (QNE reference).
// SensorTemperature is in Celsius. QNHAirport and SensorPressure are in hPa.
// SensorElevation, N, and DSM are in meters.
type ADSBInput struct {
	HObsQne          float64
	SensorTemperature float64
	QNHAirport        float64
	SensorPressure    float64
	SensorElevation   float64
	N                 float64
	DSM               float64
}

// GNSSInput holds all parameters needed to calculate altitudes from a GNSS source.
// HEll is the ellipsoidal height in meters.
// SensorTemperature is in Celsius. QNHAirport and SensorPressure are in hPa.
// SensorElevation, N, and DSM are in meters.
type GNSSInput struct {
	HEll              float64
	SensorTemperature float64
	QNHAirport        float64
	SensorPressure    float64
	SensorElevation   float64
	N                 float64
	DSM               float64
}

// AltitudeResult contains all calculated altitude values in meters.
type AltitudeResult struct {
	HObsQne float64 `json:"h_obs_qne"`
	HObsQnh float64 `json:"h_obs_qnh"`
	HOrt    float64 `json:"h_ort"`
	HEll    float64 `json:"h_ell"`
	HAsl    float64 `json:"h_asl"`
	HAgl    float64 `json:"h_agl"`
}
