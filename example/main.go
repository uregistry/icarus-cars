package main

import (
	"fmt"
	cars "icarus-cars"
)

func main() {
	exampleFromADSB()
	fmt.Println()
	exampleFromGNSS()
}

// exampleFromADSB demonstrates calculating altitudes for a vehicle that
// reports its pressure altitude via ADS-B transponder.
//
// The workflow is:
//  1. Receive HObsQne (pressure altitude) from the vehicle's transponder.
//  2. Fetch atmospheric data from a nearby weather station:
//     SensorPressure, SensorTemperature, SensorElevation.
//  3. Fetch QNH from the nearest airport.
//  4. Fetch N (geoid undulation) and DSM (terrain height) for the vehicle's
//     position from your terrain/geoid service — these are the only two values
//     that previously required an external HTTP call; now you supply them here.
//  5. Call CalculateFromADSB and use the result.
func exampleFromADSB() {
	fmt.Println("=== ADS-B example ===")

	// Step 1: pressure altitude reported by the vehicle transponder (meters).
	hObsQne := 100.0

	// Step 2: nearest weather station readings.
	sensorPressure    := 1014.1 // hPa — atmospheric pressure at station
	sensorTemperature := 15.0   // °C  — temperature at station
	sensorElevation   := 250.0  // m   — elevation of the station above MSL

	// Step 3: QNH broadcast by the nearest airport (hPa).
	qnhAirport := 1015.0

	// Step 4: values from your terrain/geoid data source for the vehicle's lat/lon.
	n   := 45.2  // m — geoid undulation (EGM96 or similar)
	dsm := 252.0 // m — terrain height from digital surface model

	result := cars.CalculateFromADSB(cars.ADSBInput{
		HObsQne:           hObsQne,
		SensorTemperature: sensorTemperature,
		QNHAirport:        qnhAirport,
		SensorPressure:    sensorPressure,
		SensorElevation:   sensorElevation,
		N:                 n,
		DSM:               dsm,
	})

	printResult(result)
}

// exampleFromGNSS demonstrates calculating altitudes for a vehicle that
// reports its position via GNSS (GPS/GLONASS/Galileo).
//
// The workflow is:
//  1. Receive HEll (ellipsoidal height) from the vehicle's GNSS receiver.
//  2. Fetch atmospheric data from a nearby weather station.
//  3. Fetch QNH from the nearest airport.
//  4. Fetch N and DSM for the vehicle's position from your terrain/geoid service.
//  5. Call CalculateFromGNSS and use the result.
func exampleFromGNSS() {
	fmt.Println("=== GNSS example ===")

	// Step 1: ellipsoidal height reported by the vehicle's GNSS receiver (meters).
	hEll := 395.0

	// Step 2: nearest weather station readings.
	sensorPressure    := 1016.77 // hPa
	sensorTemperature := 18.0    // °C
	sensorElevation   := 88.0    // m

	// Step 3: QNH from nearest airport (hPa).
	qnhAirport := 1017.0

	// Step 4: values from your terrain/geoid data source.
	n   := 44.8 // m — geoid undulation
	dsm := 90.0 // m — terrain height from DSM

	result := cars.CalculateFromGNSS(cars.GNSSInput{
		HEll:              hEll,
		SensorTemperature: sensorTemperature,
		QNHAirport:        qnhAirport,
		SensorPressure:    sensorPressure,
		SensorElevation:   sensorElevation,
		N:                 n,
		DSM:               dsm,
	})

	printResult(result)
}

func printResult(r cars.AltitudeResult) {
	fmt.Printf("  H_obs_QNE  (pressure alt, QNE ref)  : %8.2f m\n", r.HObsQne)
	fmt.Printf("  H_obs_QNH  (pressure alt, QNH ref)  : %8.2f m\n", r.HObsQnh)
	fmt.Printf("  H_ort      (orthometric / geoid alt) : %8.2f m\n", r.HOrt)
	fmt.Printf("  H_ell      (ellipsoidal alt)         : %8.2f m\n", r.HEll)
	fmt.Printf("  H_ASL      (above sea level)         : %8.2f m\n", r.HAsl)
	fmt.Printf("  H_AGL      (above ground level)      : %8.2f m\n", r.HAgl)
}
