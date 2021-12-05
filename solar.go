package solar

import (
	"math"
	"time"
)


func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

func rad2deg(rad float64) float64 {
	return rad * 180 / math.Pi
}

// returns the number of minutes to add to mean solar time to get actual solar time.
func EquationOfTime(day float64) float64 {
	b := 2 * math.Pi / 364.0 * (day - 81)
	return 9.87 * math.Sin(2 * b) - 7.53 * math.Cos(b) - 1.5 * math.Sin(b)
}

// sun-earth distance is in astronomical units
func GetAberationCorrection(sunEarthDistance float64) float64 {
	return  -20.4898 / (3600.0 * sunEarthDistance)
}

// Common calculations for altitude and azimuth
func GetTopocentricPosition(lat, lon, elevation float64, when time.Time) (float64, float64) {
	// location-dependent calculations
	projectedRadialDistance := GetProjectedRadialDistance(elevation, lat)
	projectedAxialDistance := GetProjectedAxialDistance(elevation, lat)

	// time-dependent calculations
	jd := GetJulianSolarDay(when)
	jde := GetJulianEphemerisDay(when)
	jce := GetJulianEphemerisCentury(jde)
	jme := GetJulianEphemerisMillenium(jce)
	geocentricLatitude := GetGeocentricLatitude(jme)
	geocentricLongitude := GetGeocentricLongitude(jme)
	sunEarthDistance := GetSunEarthDistance(jme)
	aberrationCorrection := GetAberationCorrection(sunEarthDistance)
	equatorialHorizontalParallax := GetEquatorialHorizontalParallax(sunEarthDistance)
	nutation := GetNutation(jce)
	apparentSiderealTime := GetApparentSiderealTime(jd, jme, nutation)
	trueEclipticObliquity := GetTrueEclipticObliquity(jme, nutation)

	// calculations dependent on location and time
	apparentSunLongitude := GetApparentSunLongitude(geocentricLongitude, nutation, aberrationCorrection)
	geocentricSunRightAscension := GetGeocentricSunRightAscension(apparentSunLongitude, trueEclipticObliquity, geocentricLatitude)
	geocentricSunDeclination := GetGeocentricSunDeclination(apparentSunLongitude, trueEclipticObliquity, geocentricLatitude)
	localHourAngle := GetLocalHourAngle(apparentSiderealTime, lon, geocentricSunRightAscension)
	parallaxSunRightAscension := GetParallaxSunRightAscension(projectedRadialDistance, equatorialHorizontalParallax, localHourAngle, geocentricSunDeclination)
	topocentricLocalHourAngle := GetTopocentricLocalHourAngle(localHourAngle, parallaxSunRightAscension)
	topocentricSunDeclination := GetTopocentricSunDeclination(geocentricSunDeclination, projectedAxialDistance, equatorialHorizontalParallax, parallaxSunRightAscension, localHourAngle)
	return topocentricSunDeclination, topocentricLocalHourAngle
}

/*
Given location, time and atmospheric conditions
temperature in Kelvin and pressure in Pascal
	
returns (azimuth, altitude) of sun in degrees.

Same as a combination of get_azimuth and get_altitude
*/
func GetPosition(lat, lon, elevation float64, when time.Time, temperature, pressure *float64) (float64, float64) {
	var temp, pres float64
	if temperature == nil {
		temp = StandardTemperature
	} else {
		temp = *temperature
	}
	if pressure == nil {
		pres = StandardPressure
	} else {
		pres = *pressure
	}
	topocentricSunDeclination, topocentricLocalHourAngle := GetTopocentricPosition(lat, lon, elevation, when)
	topocentricElevationAngle := GetTopocentricElevationAngle(lat, topocentricSunDeclination, topocentricLocalHourAngle)
	refractionCorrection := GetRefractionCorrection(pres, temp, topocentricElevationAngle)
	altitudeDeg := topocentricElevationAngle + refractionCorrection
	azimuthDeg := GetTopocentricAzimuthAngle(topocentricLocalHourAngle, lat, topocentricSunDeclination)
	return altitudeDeg, azimuthDeg
}

/*
See also the faster, but less accurate, get_altitude_fast()
temperature in Kelvin and pressure in Pascal
*/
func GetAltitude(lat, lon, elevation float64, when time.Time, temperature, pressure *float64) float64 {
	var temp, pres float64
	if temperature == nil {
		temp = StandardTemperature
	} else {
		temp = *temperature
	}
	if pressure == nil {
		pres = StandardPressure
	} else {
		pres = *pressure
	}
	topocentricSunDeclination, topocentricLocalHourAngle := GetTopocentricPosition(lat, lon, elevation, when)
	topocentricElevationAngle := GetTopocentricElevationAngle(lat, topocentricSunDeclination, topocentricLocalHourAngle)
	refractionCorrection := GetRefractionCorrection(pres, temp, topocentricElevationAngle)
	return topocentricElevationAngle + refractionCorrection
}






func GetAltitudeFast(latitudeDeg, longitudeDeg float64, when time.Time) float64 {
	// expect 19 degrees for GetAltitude(42.364908,-71.112828,0,time.Date(2007, time.February, 18, 20, 13, 1, 130320000),nil, nil)
	day := when.YearDay()
	declinationRad := deg2rad(GetDeclination(float64(day)))
	latitudeRad := deg2rad(latitudeDeg)
	hourAngle := GetHourAngle(when, longitudeDeg)
	firstTerm := math.Cos(latitudeRad) * math.Cos(declinationRad) * math.Cos(deg2rad(hourAngle))
	secondTerm := math.Sin(latitudeRad) * math.Sin(declinationRad)
	return rad2deg(math.Asin(firstTerm + secondTerm))
}


func GetApparentSiderealTime(jd, jme float64, nutation map[string]float64) float64 {
	return GetMeanSiderealTime(jd) + nutation["longitude"] * math.Cos(GetTrueEclipticObliquity(jme, nutation))
}

func GetApparentSunLongitude(geocentricLongitude float64, nutation map[string]float64, abCorrection float64) float64 {
	return geocentricLongitude + nutation["longitude"] + abCorrection
}

func GetAzimuth(latitudeDeg, longitudeDeg, elevation float64, when time.Time) float64 {
	topocentricSunDeclination, topocentricLocalHourAngle := GetTopocentricPosition(latitudeDeg, longitudeDeg, elevation, when)
	return GetTopocentricAzimuthAngle(topocentricLocalHourAngle, latitudeDeg, topocentricSunDeclination)
}

func GetAzimuthFast(latitudeDeg, longitudeDeg float64, when time.Time) float64 {
	// expect 230 degrees for GetAzimuth(42.364908,-71.112828,0,time.Date(2007, time.February, 18, 20, 18, 0, 0))
	day := when.YearDay()
	declinationRad := deg2rad(GetDeclination(float64(day)))
	latitudeRad := deg2rad(latitudeDeg)
	hourAngleRad := deg2rad(GetHourAngle(when, longitudeDeg))
	altitudeRad := deg2rad(GetAltitudeFast(latitudeDeg, longitudeDeg, when))
	azimuthRad := math.Asin(-math.Cos(declinationRad) * math.Sin(hourAngleRad) / math.Cos(altitudeRad))

	if math.Cos(hourAngleRad) * math.Tan(latitudeRad) >= math.Tan(declinationRad) {
		return 180 - rad2deg(azimuthRad)
	}
	azimuthDeg := rad2deg(azimuthRad)
	if azimuthRad < 0 {
		azimuthDeg += 360
	}
	return azimuthDeg
}

/*
computes a polynomial with time-varying coefficients from the given constant
coefficients array and the current Julian millennium.
*/
func GetCoeff(jme float64, coeffs [][][3]float64) float64 {
	result := 0.0
	x := 1.0
	for _, line := range coeffs {
		c := 0.0
		for _, l := range line {
			c += l[0] * math.Cos(l[1] + l[2] * jme)
		}
		result += c * x
		x *= jme
	}
	return result
}

/*
The declination of the sun is the angle between
Earth's equatorial plane and a line between the Earth and the sun.
The declination of the sun varies between 23.45 degrees and -23.45 degrees,
hitting zero on the equinoxes and peaking on the solstices.
*/
func GetDeclination(day float64) float64 {
	return EarthAxisInclination * math.Sin((2 * math.Pi / 365.0) * (day - 81))
}

func GetEquatorialHorizontalParallax(sunEarthDistance float64) float64 {
	return 8.794 / (3600 / sunEarthDistance)
}

func GetFlattenedLatitude(latitude float64) float64 {
	latitudeRad := deg2rad(latitude)
	return rad2deg(math.Atan(0.99664719 * math.Tan(latitudeRad)))
}

// Geocentric functions calculate angles relative to the center of the earth.
func GetGeocentricLatitude(jme float64) float64 {
	return -1 * GetHeliocentricLatitude(jme)
}

func GetGeocentricLongitude(jme float64) float64 {
	return math.Mod(GetHeliocentricLongitude(jme) + 180, 360)
}

func GetGeocentricSunDeclination(apparentSunLongitude, trueEclipticObliquity, geocentricLatitude float64) float64 {
	apparentSunLongitudeRad := deg2rad(apparentSunLongitude)
	trueEclipticObliquityRad := deg2rad(trueEclipticObliquity)
	geocentricLatitudeRad := deg2rad(geocentricLatitude)
	a := math.Sin(geocentricLatitudeRad) * math.Cos(trueEclipticObliquityRad)
	b := math.Cos(geocentricLatitudeRad) * math.Sin(trueEclipticObliquityRad) * math.Sin(apparentSunLongitudeRad)
	delta := math.Asin(a + b)
	return rad2deg(delta)
}

func GetGeocentricSunRightAscension(apparentSunLongitude, trueEclipticObliquity, geocentricLatitude float64) float64 {
	apparentSunLongitudeRad := deg2rad(apparentSunLongitude)
	trueEclipticObliquityRad := deg2rad(trueEclipticObliquity)
	geocentricLatitudeRad := deg2rad(geocentricLatitude)

	a := math.Sin(apparentSunLongitudeRad) * math.Cos(trueEclipticObliquityRad)
	b := math.Tan(geocentricLatitudeRad) * math.Sin(trueEclipticObliquityRad)
	c := math.Cos(apparentSunLongitudeRad)
	alpha := math.Atan2((a - b),  c)
	return math.Mod(rad2deg(alpha), 360)
}

// Heliocentric functions calculate angles relative to the center of the sun.
func GetHeliocentricLatitude(jme float64) float64 {
	return rad2deg(GetCoeff(jme, HeliocentricLatitudeCoeffs) / 1e8)
}

func GetHeliocentricLongitude(jme float64) float64 {
	return math.Mod(rad2deg(GetCoeff(jme, HeliocentricLongitudeCoeffs) / 1e8), 360)
}

func GetHourAngle(when time.Time, longitudeDeg float64) float64 {
	solarTime := GetSolarTime(longitudeDeg, when)
	return 15.0 * (solarTime - 12.0)
}

func GetIncidenceAngle(topocentricZenithAngle, slope, slopeOrientation, topocentricAzimuthAngle float64) float64 {
	tzaRad := deg2rad(topocentricZenithAngle)
	slopeRad := deg2rad(slope)
	soRad := deg2rad(slopeOrientation)
	taaRad := deg2rad(topocentricAzimuthAngle)
	return rad2deg(math.Acos(math.Cos(tzaRad) * math.Cos(slopeRad) + math.Sin(slopeRad) * math.Sin(tzaRad) * math.Cos(taaRad - math.Pi - soRad)))
}

func GetLocalHourAngle(apparentSiderealTime, longitude, geocentricSunRightAscension float64) float64 {
	return math.Mod(apparentSiderealTime + longitude - geocentricSunRightAscension, 360)
}

// This function doesn't agree with Andreas and Reda as well as it should. Works to ~5 sig figs in current unit test
func GetMeanSiderealTime(jd float64) float64 {
	jc := GetJulianCentury(jd)
	siderealTime :=  280.46061837 + (360.98564736629 * (jd - 2451545.0)) + 0.000387933 * jc * jc * (1 - jc / 38710000)
	return math.Mod(siderealTime, 360)
}

func sum(vals []float64) float64 {
	s := float64(0)
	for _, v := range vals {
		s += v
	}
	return s
}

func GetNutation(jce float64) map[string]float64 {
	abcd := NutationCoefficients
	nutationLong := []float64{}
	nutationOblique := []float64{}
	p := GetAberrationCoeffs()
	// order here matters
	keys := []string{
		"MeanElongationOfMoon",
		"MeanAnomalyOfSun",
		"MeanAnomalyOfMoon",
		"ArgumentOfLatitudeOfMoon",
		"LongitudeOfAscendingNode",
	}
	x := make([]float64, len(keys))
	for i, k := range keys {
		x[i] = p[k](jce)
	}
	y := AberrationSinTerms
	for i := range abcd {
		sigmaxy := 0.0
		for j := range x {
			sigmaxy += x[j] * y[i][j]
		}
		nutationLong = append(nutationLong, (abcd[i][0] + (abcd[i][1] * jce)) * math.Sin(deg2rad(sigmaxy)))
		nutationOblique = append(nutationOblique, (abcd[i][2] + (abcd[i][3] * jce)) * math.Cos(deg2rad(sigmaxy)))
	}
	// 36000000 scales from 0.0001 arcseconds to degrees
	return map[string]float64{
		"longitude": sum(nutationLong) / 36000000.0,
		"obliquity": sum(nutationOblique) / 36000000.0,
	}
}

func GetParallaxSunRightAscension(projectedRadialDistance, equatorialHorizontalParallax, localHourAngle, geocentricSunDeclination float64) float64 {
	prd := projectedRadialDistance
	ehpRad := deg2rad(equatorialHorizontalParallax)
	lhaRad := deg2rad(localHourAngle)
	gsdRad := deg2rad(geocentricSunDeclination)
	a := -1 * prd * math.Sin(ehpRad) * math.Sin(lhaRad)
	b :=  math.Cos(gsdRad) - prd * math.Sin(ehpRad) * math.Cos(lhaRad)
	parallax := math.Atan2(a, b)
	return rad2deg(parallax)
}

func GetProjectedRadialDistance(elevation, latitude float64) float64 {
	flattenedLatitudeRad := deg2rad(GetFlattenedLatitude(latitude))
	latitudeRad := deg2rad(latitude)
	return math.Cos(flattenedLatitudeRad) + (elevation * math.Cos(latitudeRad) / EarthRadius)
}

func GetProjectedAxialDistance(elevation, latitude float64) float64 {
	flattenedLatitudeRad := deg2rad(GetFlattenedLatitude(latitude))
	latitudeRad := deg2rad(latitude)
	return 0.99664719 * math.Sin(flattenedLatitudeRad) + (elevation * math.Sin(latitudeRad) / EarthRadius)
}

func GetSunEarthDistance(jme float64) float64 {
	return GetCoeff(jme, SunEarthDistanceCoeffs) / 1e8
}

func GetRefractionCorrection(pressure, temperature, topocentricElevationAngle float64) float64 {
	// function and default values according to original NREL SPA C code
	// http://www.nrel.gov/midc/spa/

	sunRadius := 0.26667
	atmosRefract := 0.5667
	tea := topocentricElevationAngle

	// Approximation only valid if sun is not well below horizon
	// This approximation could be improved; see history at https://github.com/pingswept/pysolar/pull/23
	// Better method could come from Auer and Standish [2000]:
	// http://iopscience.iop.org/1538-3881/119/5/2472/pdf/1538-3881_119_5_2472.pdf

	a := pressure * 2.830 * 1.02
	b := 1010.0 * temperature * 60.0 * math.Tan(deg2rad(tea + (10.3/(tea + 5.11))))

	if tea >= -1.0*(sunRadius + atmosRefract) {
		return a / b
	}
	return 0
}

/*
returns solar time in hours for the specified longitude and time,
accurate only to the nearest minute.
*/
func GetSolarTime(longitudeDeg float64, when time.Time) float64 {
	return float64(when.Hour()) * 60 + float64(when.Minute()) + 4 * longitudeDeg + EquationOfTime(float64(when.YearDay())) / 60
}

// Topocentric functions calculate angles relative to a location on the surface of the earth.

// West is negative, East is positive, Masters p. 395
func GetTopocentricAzimuthAngle(topocentricLocalHourAngle, latitude, topocentricSunDeclination float64) float64 {
	tlhaRad := deg2rad(topocentricLocalHourAngle)
	latitudeRad := deg2rad(latitude)
	tsdRad := deg2rad(topocentricSunDeclination)
	a := math.Sin(tlhaRad)
	b := math.Cos(tlhaRad) * math.Sin(latitudeRad) - math.Tan(tsdRad) * math.Cos(latitudeRad)
	return math.Mod(180.0 + rad2deg(math.Atan2(a, b)), 360)
}

func GetTopocentricElevationAngle(latitude, topocentricSunDeclination, topocentricLocalHourAngle float64) float64 {
	latitudeRad := deg2rad(latitude)
	tsdRad := deg2rad(topocentricSunDeclination)
	tlhaRad := deg2rad(topocentricLocalHourAngle)
	return rad2deg(math.Asin((math.Sin(latitudeRad) * math.Sin(tsdRad)) + math.Cos(latitudeRad) * math.Cos(tsdRad) * math.Cos(tlhaRad)))
}

func GetTopocentricLocalHourAngle(localHourAngle, parallaxSunRightAscension float64) float64 {
	return localHourAngle - parallaxSunRightAscension
}

func GetTopocentricSunDeclination(geocentricSunDeclination, projectedAxialDistance, equatorialHorizontalParallax, parallaxSunRightAscension, localHourAngle float64) float64 {
	gsdRad := deg2rad(geocentricSunDeclination)
	pad := projectedAxialDistance
	ehpRad := deg2rad(equatorialHorizontalParallax)
	psraRad := deg2rad(parallaxSunRightAscension)
	lhaRad := deg2rad(localHourAngle)
	a := (math.Sin(gsdRad) - pad * math.Sin(ehpRad)) * math.Cos(psraRad)
	b := math.Cos(gsdRad) - (pad * math.Sin(ehpRad) * math.Cos(lhaRad))
	return rad2deg(math.Atan2(a, b))
}

func GetTopocentricSunRightAscension(projectedRadialDistance, equatorialHorizontalParallax, localHourAngle, apparentSunLongitude, trueEclipticObliquity, geocentricLatitude float64) float64 {
	gsd := GetGeocentricSunDeclination(apparentSunLongitude, trueEclipticObliquity, geocentricLatitude)
	psra := GetParallaxSunRightAscension(projectedRadialDistance, equatorialHorizontalParallax, localHourAngle, gsd)
	gsra := GetGeocentricSunRightAscension(apparentSunLongitude, trueEclipticObliquity, geocentricLatitude)
	return psra + gsra
}

func GetTopocentricZenithAngle(latitude, topocentricSunDeclination, topocentricLocalHourAngle, pressure, temperature float64) float64 {
	tea := GetTopocentricElevationAngle(latitude, topocentricSunDeclination, topocentricLocalHourAngle)
	return 90 - tea - GetRefractionCorrection(pressure, temperature, tea)
}

func GetTrueEclipticObliquity(jme float64, nutation map[string]float64) float64 {
	u := jme / 10.0
	u2 := u * u
	u3 := u2 * u
	u4 := u3 * u
	u5 := u4 * u
	u6 := u5 * u
	u7 := u6 * u
	u8 := u7 * u
	u9 := u8 * u
	u10 := u9 * u
	meanObliquity := 84381.448 - (4680.93 * u) - (1.55 * u2) + (1999.25 * u3) - (51.38 * u4) -(249.67 * u5) - (39.05 * u6) + (7.12 * u7) + (27.87 * u8) + (5.79 * u9) + (2.45 * u10)
	return (meanObliquity / 3600.0) + nutation["obliquity"]
}








/*

func GetFlattenedLatitude(lat float64) float64 {
	latitudeRad := deg2rad(lat)
	return lat2deg(math.Atan(0.99664719 * math.Tan(latitudeRad)))
}

func GetProjectedRadialDistance(elevation, lat float64) float64 {
	flattenedLatitudeRad := deg2rad(GetFlattenedLatitude(lat))
	latitudeRad := deg2rad(lat)
	return math.Cos(flattenedLatitudeRad) + (elevation * math.Cos(latitudeRad) / EarthRadius)
}

func GetProjectedAxialDistance(elevation, lat float64) float64 {
	flattenedLatitudeRad := deg2rad(GetFlattenedLatitude(lat))
	latitudeRad := deg2rad(lat)
	return 0.99664719 * math.Sin(flattenedLatitudeRad) + (elevation * math.Sin(latitudeRad) / EarthRadius)
}

func GetHeliocentricLatitude(jme float64) float64 {
	return rad2deg(GetCoeff(jme, HeliocentricLatitudeCoeffs) / 1e8)
}

func GetHeliocentricLongitude(jme float64) float64 {
	return rad2deg(GetCoeff(jme, HeliocentricLongitudeCoeffs) / 1e8) % 360.0
}

func GetGeocentricLatitude(jme float64) float64 {
	return -1 * GetHeliocentricLatitude(jme)
}

func GetGeocentricLongitude
func GetSunEarthDistance
func GetEquatorialHorizontalParalax
func GetNutation
func GetApparentSiderealTime
func GetTrueEclipticObliquity
func GetApparentSunLongitude
func GetGeocentricSunRightAscension
func GetGeocentricSunDeclination
func GetLocalHourAngle
func GetParallaxSunRightAscension
func GetTopocentricLocalHourAngle
func GetTopocentricSunDeclination

func GetTopocentricElevationAngle(lat, topocentricSunDeclination, topocentricLocalHourAngle float64) float64 {
}

func GetRefractionCorrection(pressure, temperature, topocentricElevationAngle float64) float64 {
}
*/
