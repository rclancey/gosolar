package solar

/*
This file is consists of numerical constants for calculating corrections,
such as the wiggling ("nutation") of the axis of the earth. It also includes
functions for building dictionaries of polynomial functions for rapid
calculation of corrections.

Most of the constants come from a 2005 paper by Reda and Andreas:

I. Reda and A. Andreas, "Solar Position Algorithm for Solar Radiation
Applications," National Renewable Energy Laboratory, NREL/TP-560-34302,
revised November 2005.

http://www.osti.gov/bridge/servlets/purl/15003974-iP3z6k/native/15003974.PDF

However, it seems that Reda and Andreas took the bulk of the constants
(L0, etc.) from Pierre Bretagnon and Gerard Francou's Variations Seculaires
des Orbites Planetaires, or VSOP87:

http://en.wikipedia.org/wiki/Secular_variations_of_the_planetary_orbits#VSOP87

See also ftp://ftp.imcce.fr/pub/ephem/planets/vsop87/VSOP87D.ear
*/
const (
	EarthRadius = float64(6378140) // meters
	EarthAxisInclination = float64(23.45) // degrees
	StandardPressure = float64(101325.00) // pascals
	StandardTemperature = float64(288.15) // kelvin
	CelsiusOffset = float64(273.15)
	EarthTemperatureLapseRate = float64(-0.0065) //change in temperature with height, kevin/meter
	AirGasConstant = float64(8.31432) // N*m/s^2
	EarthGravity = float64(9.80665) // m/s^2
	EarthAtmosphereMolarMass = 0.0289644 // kg/mol
)

var aberationCoeffs map[string]func(float64) float64

func GetAberrationCoeffs() map[string]func(float64) float64 {
	if aberationCoeffs == nil {
		makeCoeff := func(a, b, c, d float64) func(float64) float64 {
			return func(x float64) float64 {
				return a + (b * x) + (c * (x * x)) + (x * x * x) / d
			}
		}
		aberationCoeffs = map[string]func(float64) float64{
			"ArgumentOfLatitudeOfMoon": makeCoeff(93.27191, 483202.017538, -0.0036825, 327270.0),
			"LongitudeOfAscendingNode": makeCoeff(125.04452, -1934.136261, 0.0020708, 450000.0),
			"MeanElongationOfMoon": makeCoeff(297.85036, 445267.111480, -0.0019142, 189474.0),
			"MeanAnomalyOfMoon": makeCoeff(134.96298, 477198.867398, 0.0086972, 56250.0),
			"MeanAnomalyOfSun": makeCoeff(357.52772, 35999.050340, -0.0001603, -300000.0),
		}
	}
	return aberationCoeffs
}

var AberrationSinTerms = [][5]float64{
	[5]float64{0,0,0,0,1},
	[5]float64{-2,0,0,2,2},
	[5]float64{0,0,0,2,2},
	[5]float64{0,0,0,0,2},
	[5]float64{0,1,0,0,0},
	[5]float64{0,0,1,0,0},
	[5]float64{-2,1,0,2,2},
	[5]float64{0,0,0,2,1},
	[5]float64{0,0,1,2,2},
	[5]float64{-2,-1,0,2,2},
	[5]float64{-2,0,1,0,0},
	[5]float64{-2,0,0,2,1},
	[5]float64{0,0,-1,2,2},
	[5]float64{2,0,0,0,0},
	[5]float64{0,0,1,0,1},
	[5]float64{2,0,-1,2,2},
	[5]float64{0,0,-1,0,1},
	[5]float64{0,0,1,2,1},
	[5]float64{-2,0,2,0,0},
	[5]float64{0,0,-2,2,1},
	[5]float64{2,0,0,2,2},
	[5]float64{0,0,2,2,2},
	[5]float64{0,0,2,0,0},
	[5]float64{-2,0,1,2,2},
	[5]float64{0,0,0,2,0},
	[5]float64{-2,0,0,2,0},
	[5]float64{0,0,-1,2,1},
	[5]float64{0,2,0,0,0},
	[5]float64{2,0,-1,0,1},
	[5]float64{-2,2,0,2,2},
	[5]float64{0,1,0,0,1},
	[5]float64{-2,0,1,0,1},
	[5]float64{0,-1,0,0,1},
	[5]float64{0,0,2,-2,0},
	[5]float64{2,0,-1,2,1},
	[5]float64{2,0,1,2,2},
	[5]float64{0,1,0,2,2},
	[5]float64{-2,1,1,0,0},
	[5]float64{0,-1,0,2,2},
	[5]float64{2,0,0,2,1},
	[5]float64{2,0,1,0,0},
	[5]float64{-2,0,2,2,2},
	[5]float64{-2,0,1,2,1},
	[5]float64{2,0,-2,0,1},
	[5]float64{2,0,0,0,1},
	[5]float64{0,-1,1,0,0},
	[5]float64{-2,-1,0,2,1},
	[5]float64{-2,0,0,0,1},
	[5]float64{0,0,2,2,1},
	[5]float64{-2,0,2,0,1},
	[5]float64{-2,1,0,2,1},
	[5]float64{0,0,1,-2,0},
	[5]float64{-1,0,1,0,0},
	[5]float64{-2,1,0,0,0},
	[5]float64{1,0,0,0,0},
	[5]float64{0,0,1,2,0},
	[5]float64{0,0,-2,2,2},
	[5]float64{-1,-1,1,0,0},
	[5]float64{0,1,1,0,0},
	[5]float64{0,-1,1,2,2},
	[5]float64{2,-1,-1,2,2},
	[5]float64{0,0,3,2,2},
	[5]float64{2,-1,0,2,2},
}

var NutationCoefficients = [][4]float64{
	[4]float64{-171996,-174.2,92025,8.9},
	[4]float64{-13187,-1.6,5736,-3.1},
	[4]float64{-2274,-0.2,977,-0.5},
	[4]float64{2062,0.2,-895,0.5},
	[4]float64{1426,-3.4,54,-0.1},
	[4]float64{712,0.1,-7,0},
	[4]float64{-517,1.2,224,-0.6},
	[4]float64{-386,-0.4,200,0},
	[4]float64{-301,0,129,-0.1},
	[4]float64{217,-0.5,-95,0.3},
	[4]float64{-158,0,0,0},
	[4]float64{129,0.1,-70,0},
	[4]float64{123,0,-53,0},
	[4]float64{63,0,0,0},
	[4]float64{63,0.1,-33,0},
	[4]float64{-59,0,26,0},
	[4]float64{-58,-0.1,32,0},
	[4]float64{-51,0,27,0},
	[4]float64{48,0,0,0},
	[4]float64{46,0,-24,0},
	[4]float64{-38,0,16,0},
	[4]float64{-31,0,13,0},
	[4]float64{29,0,0,0},
	[4]float64{29,0,-12,0},
	[4]float64{26,0,0,0},
	[4]float64{-22,0,0,0},
	[4]float64{21,0,-10,0},
	[4]float64{17,-0.1,0,0},
	[4]float64{16,0,-8,0},
	[4]float64{-16,0.1,7,0},
	[4]float64{-15,0,9,0},
	[4]float64{-13,0,7,0},
	[4]float64{-12,0,6,0},
	[4]float64{11,0,0,0},
	[4]float64{-10,0,5,0},
	[4]float64{-8,0,3,0},
	[4]float64{7,0,-3,0},
	[4]float64{-7,0,0,0},
	[4]float64{-7,0,3,0},
	[4]float64{-7,0,3,0},
	[4]float64{6,0,0,0},
	[4]float64{6,0,-3,0},
	[4]float64{6,0,-3,0},
	[4]float64{-6,0,3,0},
	[4]float64{-6,0,3,0},
	[4]float64{5,0,0,0},
	[4]float64{-5,0,3,0},
	[4]float64{-5,0,3,0},
	[4]float64{-5,0,3,0},
	[4]float64{4,0,0,0},
	[4]float64{4,0,0,0},
	[4]float64{4,0,0,0},
	[4]float64{-4,0,0,0},
	[4]float64{-4,0,0,0},
	[4]float64{-4,0,0,0},
	[4]float64{3,0,0,0},
	[4]float64{-3,0,0,0},
	[4]float64{-3,0,0,0},
	[4]float64{-3,0,0,0},
	[4]float64{-3,0,0,0},
	[4]float64{-3,0,0,0},
	[4]float64{-3,0,0,0},
	[4]float64{-3,0,0,0},
}

var HeliocentricLongitudeCoeffs = [][][3]float64{
	[][3]float64{ // L0
		[3]float64{175347046.0,0,0},
		[3]float64{3341656.0,4.6692568,6283.07585},
		[3]float64{34894.0,4.6261,12566.1517},
		[3]float64{3497.0,2.7441,5753.3849},
		[3]float64{3418.0,2.8289,3.5231},
		[3]float64{3136.0,3.6277,77713.7715},
		[3]float64{2676.0,4.4181,7860.4194},
		[3]float64{2343.0,6.1352,3930.2097},
		[3]float64{1324.0,0.7425,11506.7698},
		[3]float64{1273.0,2.0371,529.691},
		[3]float64{1199.0,1.1096,1577.3435},
		[3]float64{990,5.233,5884.927},
		[3]float64{902,2.045,26.298},
		[3]float64{857,3.508,398.149},
		[3]float64{780,1.179,5223.694},
		[3]float64{753,2.533,5507.553},
		[3]float64{505,4.583,18849.228},
		[3]float64{492,4.205,775.523},
		[3]float64{357,2.92,0.067},
		[3]float64{317,5.849,11790.629},
		[3]float64{284,1.899,796.298},
		[3]float64{271,0.315,10977.079},
		[3]float64{243,0.345,5486.778},
		[3]float64{206,4.806,2544.314},
		[3]float64{205,1.869,5573.143},
		[3]float64{202,2.458,6069.777},
		[3]float64{156,0.833,213.299},
		[3]float64{132,3.411,2942.463},
		[3]float64{126,1.083,20.775},
		[3]float64{115,0.645,0.98},
		[3]float64{103,0.636,4694.003},
		[3]float64{102,0.976,15720.839},
		[3]float64{102,4.267,7.114},
		[3]float64{99,6.21,2146.17},
		[3]float64{98,0.68,155.42},
		[3]float64{86,5.98,161000.69},
		[3]float64{85,1.3,6275.96},
		[3]float64{85,3.67,71430.7},
		[3]float64{80,1.81,17260.15},
		[3]float64{79,3.04,12036.46},
		[3]float64{75,1.76,5088.63},
		[3]float64{74,3.5,3154.69},
		[3]float64{74,4.68,801.82},
		[3]float64{70,0.83,9437.76},
		[3]float64{62,3.98,8827.39},
		[3]float64{61,1.82,7084.9},
		[3]float64{57,2.78,6286.6},
		[3]float64{56,4.39,14143.5},
		[3]float64{56,3.47,6279.55},
		[3]float64{52,0.19,12139.55},
		[3]float64{52,1.33,1748.02},
		[3]float64{51,0.28,5856.48},
		[3]float64{49,0.49,1194.45},
		[3]float64{41,5.37,8429.24},
		[3]float64{41,2.4,19651.05},
		[3]float64{39,6.17,10447.39},
		[3]float64{37,6.04,10213.29},
		[3]float64{37,2.57,1059.38},
		[3]float64{36,1.71,2352.87},
		[3]float64{36,1.78,6812.77},
		[3]float64{33,0.59,17789.85},
		[3]float64{30,0.44,83996.85},
		[3]float64{30,2.74,1349.87},
		[3]float64{25,3.16,4690.48},
	},
	[][3]float64{ // L1
		[3]float64{628331966747.0,0,0},
		[3]float64{206059.0,2.678235,6283.07585},
		[3]float64{4303.0,2.6351,12566.1517},
		[3]float64{425.0,1.59,3.523},
		[3]float64{119.0,5.796,26.298},
		[3]float64{109.0,2.966,1577.344},
		[3]float64{93,2.59,18849.23},
		[3]float64{72,1.14,529.69},
		[3]float64{68,1.87,398.15},
		[3]float64{67,4.41,5507.55},
		[3]float64{59,2.89,5223.69},
		[3]float64{56,2.17,155.42},
		[3]float64{45,0.4,796.3},
		[3]float64{36,0.47,775.52},
		[3]float64{29,2.65,7.11},
		[3]float64{21,5.34,0.98},
		[3]float64{19,1.85,5486.78},
		[3]float64{19,4.97,213.3},
		[3]float64{17,2.99,6275.96},
		[3]float64{16,0.03,2544.31},
		[3]float64{16,1.43,2146.17},
		[3]float64{15,1.21,10977.08},
		[3]float64{12,2.83,1748.02},
		[3]float64{12,3.26,5088.63},
		[3]float64{12,5.27,1194.45},
		[3]float64{12,2.08,4694},
		[3]float64{11,0.77,553.57},
		[3]float64{10,1.3,6286.6},
		[3]float64{10,4.24,1349.87},
		[3]float64{9,2.7,242.73},
		[3]float64{9,5.64,951.72},
		[3]float64{8,5.3,2352.87},
		[3]float64{6,2.65,9437.76},
		[3]float64{6,4.67,4690.48},

	},
	[][3]float64{ // L2
		[3]float64{52919.0,0,0},
		[3]float64{8720.0,1.0721,6283.0758},
		[3]float64{309.0,0.867,12566.152},
		[3]float64{27,0.05,3.52},
		[3]float64{16,5.19,26.3},
		[3]float64{16,3.68,155.42},
		[3]float64{10,0.76,18849.23},
		[3]float64{9,2.06,77713.77},
		[3]float64{7,0.83,775.52},
		[3]float64{5,4.66,1577.34},
		[3]float64{4,1.03,7.11},
		[3]float64{4,3.44,5573.14},
		[3]float64{3,5.14,796.3},
		[3]float64{3,6.05,5507.55},
		[3]float64{3,1.19,242.73},
		[3]float64{3,6.12,529.69},
		[3]float64{3,0.31,398.15},
		[3]float64{3,2.28,553.57},
		[3]float64{2,4.38,5223.69},
		[3]float64{2,3.75,0.98},
	},
	[][3]float64{ // L3
		[3]float64{289.0,5.844,6283.076},
		[3]float64{35,0,0},
		[3]float64{17,5.49,12566.15},
		[3]float64{3,5.2,155.42},
		[3]float64{1,4.72,3.52},
		[3]float64{1,5.3,18849.23},
		[3]float64{1,5.97,242.73},
	},
	[][3]float64{ // L4
		[3]float64{114.0,3.142,0},
		[3]float64{8,4.13,6283.08},
		[3]float64{1,3.84,12566.15},
	},
	[][3]float64{ // L5
		[3]float64{1,3.14,0},
	},
}

var HeliocentricLatitudeCoeffs = [][][3]float64{
	[][3]float64{ // B0
		[3]float64{280.0,3.199,84334.662},
		[3]float64{102.0,5.422,5507.553},
		[3]float64{80,3.88,5223.69},
		[3]float64{44,3.7,2352.87},
		[3]float64{32,4,1577.34},
	},
	[][3]float64{ // B1
		[3]float64{9,3.9,5507.55},
		[3]float64{6,1.73,5223.69},
	},
}

var SunEarthDistanceCoeffs = [][][3]float64{
	[][3]float64{ // R0
		[3]float64{100013989.0,0,0},
		[3]float64{1670700.0,3.0984635,6283.07585},
		[3]float64{13956.0,3.05525,12566.1517},
		[3]float64{3084.0,5.1985,77713.7715},
		[3]float64{1628.0,1.1739,5753.3849},
		[3]float64{1576.0,2.8469,7860.4194},
		[3]float64{925.0,5.453,11506.77},
		[3]float64{542.0,4.564,3930.21},
		[3]float64{472.0,3.661,5884.927},
		[3]float64{346.0,0.964,5507.553},
		[3]float64{329.0,5.9,5223.694},
		[3]float64{307.0,0.299,5573.143},
		[3]float64{243.0,4.273,11790.629},
		[3]float64{212.0,5.847,1577.344},
		[3]float64{186.0,5.022,10977.079},
		[3]float64{175.0,3.012,18849.228},
		[3]float64{110.0,5.055,5486.778},
		[3]float64{98,0.89,6069.78},
		[3]float64{86,5.69,15720.84},
		[3]float64{86,1.27,161000.69},
		[3]float64{65,0.27,17260.15},
		[3]float64{63,0.92,529.69},
		[3]float64{57,2.01,83996.85},
		[3]float64{56,5.24,71430.7},
		[3]float64{49,3.25,2544.31},
		[3]float64{47,2.58,775.52},
		[3]float64{45,5.54,9437.76},
		[3]float64{43,6.01,6275.96},
		[3]float64{39,5.36,4694},
		[3]float64{38,2.39,8827.39},
		[3]float64{37,0.83,19651.05},
		[3]float64{37,4.9,12139.55},
		[3]float64{36,1.67,12036.46},
		[3]float64{35,1.84,2942.46},
		[3]float64{33,0.24,7084.9},
		[3]float64{32,0.18,5088.63},
		[3]float64{32,1.78,398.15},
		[3]float64{28,1.21,6286.6},
		[3]float64{28,1.9,6279.55},
		[3]float64{26,4.59,10447.39},
	},
	[][3]float64{ // R1
		[3]float64{103019.0,1.10749,6283.07585},
		[3]float64{1721.0,1.0644,12566.1517},
		[3]float64{702.0,3.142,0},
		[3]float64{32,1.02,18849.23},
		[3]float64{31,2.84,5507.55},
		[3]float64{25,1.32,5223.69},
		[3]float64{18,1.42,1577.34},
		[3]float64{10,5.91,10977.08},
		[3]float64{9,1.42,6275.96},
		[3]float64{9,0.27,5486.78},
	},
	[][3]float64{ // R2
		[3]float64{4359.0,5.7846,6283.0758},
		[3]float64{124.0,5.579,12566.152},
		[3]float64{12,3.14,0},
		[3]float64{9,3.63,77713.77},
		[3]float64{6,1.87,5573.14},
		[3]float64{3,5.47,18849.23},
	},
	[][3]float64{ // R3
		[3]float64{145.0,4.273,6283.076},
		[3]float64{7,3.92,12566.15},
	},
	[][3]float64{ // R4
		[3]float64{4,2.56,6283.08},
	},
}
