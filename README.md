# icarus-cars

Pure Go library for computing CARS (Comprehensive Altitude Reference System) altitudes.
No external dependencies — all geoid/terrain data must be supplied by the caller.

## Altitude types

| Field | Name | Description |
|-------|------|-------------|
| `HObsQne` | Observed QNE altitude | Pressure altitude referenced to standard atmosphere (1013.25 hPa), in meters |
| `HObsQnh` | Observed QNH altitude | Pressure altitude corrected to local QNH, in meters |
| `HOrt` | Orthometric height | Height above the geoid (mean sea level surface), in meters |
| `HEll` | Ellipsoidal height | Height above the WGS-84 ellipsoid, in meters |
| `HAsl` | Height above sea level | Orthometric height minus terrain (DSM), in meters |
| `HAgl` | Height above ground level | Same as HAsl (above the surface model), in meters |

## Installation

```bash
go get icarus-cars
```

> Once published, replace the module path with the real one, e.g. `github.com/your-org/icarus-cars`.

## Usage

### Vehicle with ADS-B transponder

Use `CalculateFromADSB` when the vehicle reports its pressure altitude (QNE) via an ADS-B transponder.

```go
import cars "icarus-cars"

result := cars.CalculateFromADSB(cars.ADSBInput{
    HObsQne:           100.0,    // m   — pressure altitude from transponder
    SensorTemperature: 15.0,     // °C  — temperature at nearest weather station
    QNHAirport:        1015.0,   // hPa — QNH from nearest airport
    SensorPressure:    1014.1,   // hPa — pressure at weather station
    SensorElevation:   250.0,    // m   — elevation of weather station
    N:                 45.2,     // m   — geoid undulation at vehicle position
    DSM:               252.0,    // m   — terrain height at vehicle position
})

fmt.Printf("Orthometric height: %.2f m\n", result.HOrt)
fmt.Printf("Height above ground: %.2f m\n", result.HAgl)
```

### Vehicle with GNSS receiver

Use `CalculateFromGNSS` when the vehicle reports its ellipsoidal height from a GNSS receiver (GPS/GLONASS/Galileo).

```go
result := cars.CalculateFromGNSS(cars.GNSSInput{
    HEll:              395.0,    // m   — ellipsoidal height from GNSS
    SensorTemperature: 18.0,     // °C  — temperature at nearest weather station
    QNHAirport:        1017.0,   // hPa — QNH from nearest airport
    SensorPressure:    1016.77,  // hPa — pressure at weather station
    SensorElevation:   88.0,     // m   — elevation of weather station
    N:                 44.8,     // m   — geoid undulation at vehicle position
    DSM:               90.0,     // m   — terrain height at vehicle position
})

fmt.Printf("Orthometric height: %.2f m\n", result.HOrt)
fmt.Printf("Height above ground: %.2f m\n", result.HAgl)
```

## Input parameters

All inputs use SI units unless noted.

| Parameter | Unit | Description |
|-----------|------|-------------|
| `HObsQne` | m | Pressure altitude at standard atmosphere (from ADS-B) |
| `HEll` | m | Ellipsoidal height (from GNSS receiver) |
| `SensorTemperature` | °C | Temperature measured at the nearest weather station |
| `QNHAirport` | hPa | QNH (local pressure at MSL) from the nearest airport METAR |
| `SensorPressure` | hPa | Atmospheric pressure measured at the weather station |
| `SensorElevation` | m | Elevation of the weather station above MSL |
| `N` | m | Geoid undulation at the vehicle's lat/lon (e.g. from EGM96/EGM2008) |
| `DSM` | m | Terrain height at the vehicle's lat/lon from a digital surface model |

`N` and `DSM` are the two values that typically come from an external geoid/terrain service.
This library does not fetch them — fetch them yourself and pass them in.

## Running the example

```bash
cd example
go run main.go
```

Expected output:

```
=== ADS-B example ===
  H_obs_QNE  (pressure alt, QNE ref)  :   100.00 m
  H_obs_QNH  (pressure alt, QNH ref)  :   ...
  H_ort      (orthometric / geoid alt) :   ...
  H_ell      (ellipsoidal alt)         :   ...
  H_ASL      (above sea level)         :   ...
  H_AGL      (above ground level)      :   ...

=== GNSS example ===
  ...
```
