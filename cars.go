package cars

import "math"

const (
	lapseRate           float64 = 0.0065
	standardPressure    float64 = 1013.25
	meterFeetRate       float64 = 0.3048
	hpaFeetRate         float64 = 27.3
	feetMeterRate       float64 = 3.2808
	gravityAcceleration float64 = 9.80665
	gasConstant         float64 = 287.056
)

// CalculateFromADSB computes all altitude values for a vehicle reporting pressure
// altitude via ADS-B. The caller is responsible for supplying geoid undulation (N)
// and terrain height (DSM) obtained from whatever external source they use.
func CalculateFromADSB(in ADSBInput) AltitudeResult {
	hQneFeet := in.HObsQne * feetMeterRate
	qnh := qnhFromQNE(hQneFeet, in.SensorTemperature, in.QNHAirport)
	hOrt := hOrtFromQNE(hQneFeet, in.SensorPressure, in.SensorTemperature, in.SensorElevation)
	hEll := hEllFromQNE(hOrt, in.N)
	asl := aslFromQNE(hOrt, in.DSM)
	return AltitudeResult{
		HObsQne: in.HObsQne,
		HObsQnh: qnh,
		HOrt:    hOrt,
		HEll:    hEll,
		HAsl:    asl,
		HAgl:    asl,
	}
}

// CalculateFromGNSS computes all altitude values for a vehicle reporting ellipsoidal
// height via GNSS. The caller is responsible for supplying geoid undulation (N) and
// terrain height (DSM) obtained from whatever external source they use.
func CalculateFromGNSS(in GNSSInput) AltitudeResult {
	hOrt := hOrtFromGNSS(in.HEll, in.N)
	asl := aslFromGNSS(hOrt, in.DSM)
	qnh := qnhFromGNSS(in.SensorTemperature, in.SensorPressure, in.SensorElevation, in.QNHAirport, hOrt)
	qne := qneFromGNSS(hOrt, in.SensorTemperature, in.SensorPressure, in.SensorElevation)
	return AltitudeResult{
		HObsQne: qne,
		HObsQnh: qnh,
		HOrt:    hOrt,
		HEll:    in.HEll,
		HAsl:    asl,
		HAgl:    asl,
	}
}

// --- internal calculations: QNE path ---

func qnhFromQNE(hQneFeet float64, sensorTemperature float64, qnhAirPort float64) float64 {
	pow := math.Pow(((hQneFeet*meterFeetRate*lapseRate)/(sensorTemperature+273.15))+1, gravityAcceleration/(gasConstant*lapseRate))
	return ((sensorTemperature + 273.15) / lapseRate) * (math.Pow(qnhAirPort/(standardPressure/pow), (lapseRate*gasConstant)/gravityAcceleration) - 1)
}

func hOrtFromQNE(hQneFeet float64, sensorPressure float64, sensorTemperature float64, sensorElevation float64) float64 {
	pow := math.Pow(((hQneFeet*meterFeetRate*lapseRate)/(sensorTemperature+273.15))+1, gravityAcceleration/(gasConstant*lapseRate))
	return ((sensorPressure - (standardPressure / pow)) * hpaFeetRate / feetMeterRate) + sensorElevation
}

func hEllFromQNE(geoidAlt float64, n float64) float64 {
	return geoidAlt + n
}

func aslFromQNE(geoidAlt float64, dsmHeight float64) float64 {
	return geoidAlt - dsmHeight
}

// --- internal calculations: GNSS path ---

func hOrtFromGNSS(hEll float64, n float64) float64 {
	return hEll - n
}

func aslFromGNSS(hOrt float64, hDsm float64) float64 {
	return hOrt - hDsm
}

func qnhFromGNSS(sensorTemperature float64, sensorPressure float64, sensorElevation float64, qnhAirPort float64, hOrt float64) float64 {
	pow := math.Pow(qnhAirPort/(sensorPressure-((hOrt-sensorElevation)*feetMeterRate/hpaFeetRate)), (lapseRate*gasConstant)/gravityAcceleration) - 1
	return ((sensorTemperature + 273.15) / lapseRate) * pow
}

func qneFromGNSS(hOrt float64, sensorTemperature float64, sensorPressure float64, sensorElevation float64) float64 {
	pow := math.Pow(standardPressure/(sensorPressure-((hOrt-sensorElevation)*feetMeterRate/hpaFeetRate)), (lapseRate*gasConstant)/gravityAcceleration) - 1
	return ((sensorTemperature + 273.15) / lapseRate) * pow
}
