package cars_test

import (
	"math"
	"testing"

	cars "icarus-cars"
)

const delta = 0.01 // 1 cm tolerance for floating-point comparisons

func withinDelta(t *testing.T, label string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > delta {
		t.Errorf("%s: got %.6f, want %.6f (delta %.6f)", label, got, want, math.Abs(got-want))
	}
}

// ---------------------------------------------------------------------------
// CalculateFromADSB
// ---------------------------------------------------------------------------

func TestCalculateFromADSB_KnownValues(t *testing.T) {
	in := cars.ADSBInput{
		HObsQne:           100.0,
		SensorTemperature: 15.0,
		QNHAirport:        1015.0,
		SensorPressure:    1014.1,
		SensorElevation:   250.0,
		N:                 45.2,
		DSM:               252.0,
	}
	r := cars.CalculateFromADSB(in)

	withinDelta(t, "HObsQne", r.HObsQne, 100.00)
	withinDelta(t, "HObsQnh", r.HObsQnh, 114.59)
	withinDelta(t, "HOrt", r.HOrt, 356.33)
	withinDelta(t, "HEll", r.HEll, 401.53)
	withinDelta(t, "HAsl", r.HAsl, 104.33)
	withinDelta(t, "HAgl", r.HAgl, 104.33)
}

// HObsQne from the transponder must come back unchanged.
func TestCalculateFromADSB_HObsQnePassthrough(t *testing.T) {
	for _, qne := range []float64{0, 50, 100, 500, 1000} {
		in := cars.ADSBInput{
			HObsQne:           qne,
			SensorTemperature: 15.0,
			QNHAirport:        1013.25,
			SensorPressure:    1013.25,
			SensorElevation:   0.0,
			N:                 0.0,
			DSM:               0.0,
		}
		r := cars.CalculateFromADSB(in)
		if r.HObsQne != qne {
			t.Errorf("HObsQne passthrough: got %.6f, want %.6f", r.HObsQne, qne)
		}
	}
}

// HEll must equal HOrt + N (geoid undulation relationship).
func TestCalculateFromADSB_HEllEqualsHOrtPlusN(t *testing.T) {
	cases := []struct{ n float64 }{
		{0}, {10}, {45.2}, {-5},
	}
	for _, tc := range cases {
		in := cars.ADSBInput{
			HObsQne:           100.0,
			SensorTemperature: 15.0,
			QNHAirport:        1015.0,
			SensorPressure:    1014.1,
			SensorElevation:   250.0,
			N:                 tc.n,
			DSM:               0.0,
		}
		r := cars.CalculateFromADSB(in)
		if math.Abs(r.HEll-(r.HOrt+tc.n)) > 1e-9 {
			t.Errorf("N=%.1f: HEll (%.6f) != HOrt (%.6f) + N (%.1f)", tc.n, r.HEll, r.HOrt, tc.n)
		}
	}
}

// HAsl must equal HOrt − DSM; HAgl must equal HAsl.
func TestCalculateFromADSB_HAslAndHAgl(t *testing.T) {
	cases := []struct{ dsm float64 }{
		{0}, {50}, {252}, {1000},
	}
	for _, tc := range cases {
		in := cars.ADSBInput{
			HObsQne:           100.0,
			SensorTemperature: 15.0,
			QNHAirport:        1015.0,
			SensorPressure:    1014.1,
			SensorElevation:   250.0,
			N:                 45.2,
			DSM:               tc.dsm,
		}
		r := cars.CalculateFromADSB(in)
		if math.Abs(r.HAsl-(r.HOrt-tc.dsm)) > 1e-9 {
			t.Errorf("DSM=%.1f: HAsl (%.6f) != HOrt (%.6f) - DSM", tc.dsm, r.HAsl, r.HOrt)
		}
		if r.HAgl != r.HAsl {
			t.Errorf("DSM=%.1f: HAgl (%.6f) != HAsl (%.6f)", tc.dsm, r.HAgl, r.HAsl)
		}
	}
}

// ---------------------------------------------------------------------------
// CalculateFromGNSS
// ---------------------------------------------------------------------------

func TestCalculateFromGNSS_KnownValues(t *testing.T) {
	in := cars.GNSSInput{
		HEll:              395.0,
		SensorTemperature: 18.0,
		QNHAirport:        1017.0,
		SensorPressure:    1016.77,
		SensorElevation:   88.0,
		N:                 44.8,
		DSM:               90.0,
	}
	r := cars.CalculateFromGNSS(in)

	withinDelta(t, "HObsQne", r.HObsQne, 239.37)
	withinDelta(t, "HObsQnh", r.HObsQnh, 271.04)
	withinDelta(t, "HOrt", r.HOrt, 350.20)
	withinDelta(t, "HEll", r.HEll, 395.00)
	withinDelta(t, "HAsl", r.HAsl, 260.20)
	withinDelta(t, "HAgl", r.HAgl, 260.20)
}

// HEll from the GNSS receiver must come back unchanged.
func TestCalculateFromGNSS_HEllPassthrough(t *testing.T) {
	for _, hEll := range []float64{0, 100, 395, 1000, 8849} {
		in := cars.GNSSInput{
			HEll:              hEll,
			SensorTemperature: 15.0,
			QNHAirport:        1013.25,
			SensorPressure:    1013.25,
			SensorElevation:   0.0,
			N:                 0.0,
			DSM:               0.0,
		}
		r := cars.CalculateFromGNSS(in)
		if r.HEll != hEll {
			t.Errorf("HEll passthrough: got %.6f, want %.6f", r.HEll, hEll)
		}
	}
}

// HOrt must equal HEll − N (exact arithmetic).
func TestCalculateFromGNSS_HOrtEqualsHEllMinusN(t *testing.T) {
	cases := []struct{ hEll, n float64 }{
		{395.0, 44.8},
		{100.0, 0.0},
		{200.0, -10.0},
		{500.0, 50.0},
	}
	for _, tc := range cases {
		in := cars.GNSSInput{
			HEll:              tc.hEll,
			SensorTemperature: 15.0,
			QNHAirport:        1013.25,
			SensorPressure:    1013.25,
			SensorElevation:   0.0,
			N:                 tc.n,
			DSM:               0.0,
		}
		r := cars.CalculateFromGNSS(in)
		want := tc.hEll - tc.n
		if math.Abs(r.HOrt-want) > 1e-9 {
			t.Errorf("HEll=%.1f N=%.1f: HOrt (%.6f) != HEll-N (%.6f)", tc.hEll, tc.n, r.HOrt, want)
		}
	}
}

// HAsl must equal HOrt − DSM; HAgl must equal HAsl.
func TestCalculateFromGNSS_HAslAndHAgl(t *testing.T) {
	cases := []struct{ dsm float64 }{
		{0}, {90}, {250}, {350},
	}
	for _, tc := range cases {
		in := cars.GNSSInput{
			HEll:              395.0,
			SensorTemperature: 18.0,
			QNHAirport:        1017.0,
			SensorPressure:    1016.77,
			SensorElevation:   88.0,
			N:                 44.8,
			DSM:               tc.dsm,
		}
		r := cars.CalculateFromGNSS(in)
		want := r.HOrt - tc.dsm
		if math.Abs(r.HAsl-want) > 1e-9 {
			t.Errorf("DSM=%.1f: HAsl (%.6f) != HOrt (%.6f) - DSM", tc.dsm, r.HAsl, r.HOrt)
		}
		if r.HAgl != r.HAsl {
			t.Errorf("DSM=%.1f: HAgl (%.6f) != HAsl (%.6f)", tc.dsm, r.HAgl, r.HAsl)
		}
	}
}
