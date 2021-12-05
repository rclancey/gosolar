package solar

/* From pysolar.radiation */

import (
	"math"
	"time"
)

func GetAirMassRatio(altitudeDeg float64) float64 {
	return 1 / math.Sin(deg2rad(altitudeDeg))
}

func GetApparentExtraterrestrialFlux(day float64) float64 {
	return 1160 + (75 * math.Sin(2 * math.Pi / 365 * (day - 275)))
}

func GetOpticalDepth(day float64) float64 {
	return 0.174 + (0.035 * math.Sin(2 * math.Pi / 365 * (day - 100)))
}

func GetRadiationDirect(when time.Time, altitudeDeg float64) float64 {
	if altitudeDeg <= 0 {
		return 0
	}
	day := float64(when.YearDay())
	flux := GetApparentExtraterrestrialFlux(day)
	depth := GetOpticalDepth(day)
	ratio := GetAirMassRatio(altitudeDeg)
	return flux * math.Exp(-1 * depth * ratio)
}
