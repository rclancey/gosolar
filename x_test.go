package solar

import (
	"math"
	"testing"
	"time"
)

func TestGetAltitude(t *testing.T) {
	lat := 34.2245872
	lon := -118.0574345
	elev := float64(1742)
	tz, _ := time.LoadLocation("America/Los_Angeles")
	when := time.Date(2021, time.December, 4, 16, 25, 0, 0, tz).In(time.UTC)
	temp := 290.35
	pres := float64(101862)
	alt := GetAltitude(lat, lon, elev, when, &temp, &pres)
	exp := float64(2.5575097320893)
	if math.Abs(alt - exp) > 1e-6 {
		t.Errorf("expected %f, got %f", exp, alt)
	}
}

func TestGetRadiationDirect(t *testing.T) {
	tz, _ := time.LoadLocation("America/Los_Angeles")
	when := time.Date(2021, time.December, 4, 16, 25, 0, 0, tz).In(time.UTC)
	alt := 2.5575097320893
	rad := GetRadiationDirect(when, alt)
	exp := float64(47.51274557840454)
	if math.Abs(rad - exp) > 1e-6 {
		t.Errorf("expected %f, got %f", exp, rad)
	}
}
