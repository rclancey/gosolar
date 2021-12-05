package solar

import (
	"log"
	"time"
)

const (
	JulianDayOffset = float64(1721424.5) // add to datetime.datetime.toordinal() to get Julian day number
	GregorianDayOffset = float64(719163) // number of days to add to datetime.datetime.timestamp() / seconds_per_day to agree with datetime.datetime.toordinal()
	TtOffset = float64(32.184) // seconds to add to TAI to get TT
	LeapSecondsBaseYear = int(1972)
	DeltaTBaseYear = int(1973)
	DeltaTBaseMonth = time.February
)

var LeapSecondsAdjustments = [][2]int{
	// two entries per year starting from 1972, first for 23:59:59 June 30,
	// second for 23:59:59 December 31. +1 indicates that 23:59:60 follows,
	// -1 indicates that 23:59:59 does not exist, not that the latter has ever occurred.
	// source: https://www.nist.gov/pml/time-and-frequency-division/atomic-standards/leap-second-and-ut1-utc-information
	[2]int{1, 1}, // 1972
	[2]int{0, 1}, // 1973
	[2]int{0, 1}, // 1974
	[2]int{0, 1}, // 1975
	[2]int{0, 1}, // 1976
	[2]int{0, 1}, // 1977
	[2]int{0, 1}, // 1978
	[2]int{0, 1}, // 1979
	[2]int{0, 0}, // 1980
	[2]int{1, 0}, // 1981
	[2]int{1, 0}, // 1982
	[2]int{1, 0}, // 1983
	[2]int{0, 0}, // 1984
	[2]int{1, 0}, // 1985
	[2]int{0, 0}, // 1986
	[2]int{0, 1}, // 1987
	[2]int{0, 0}, // 1988
	[2]int{0, 1}, // 1989
	[2]int{0, 1}, // 1990
	[2]int{0, 0}, // 1991
	[2]int{1, 0}, // 1992
	[2]int{1, 0}, // 1993
	[2]int{1, 0}, // 1994
	[2]int{0, 1}, // 1995
	[2]int{0, 0}, // 1996
	[2]int{1, 0}, // 1997
	[2]int{0, 1}, // 1998
	[2]int{0, 0}, // 1999
	[2]int{0, 0}, // 2000
	[2]int{0, 0}, // 2001
	[2]int{0, 0}, // 2002
	[2]int{0, 0}, // 2003
	[2]int{0, 0}, // 2004
	[2]int{0, 1}, // 2005
	[2]int{0, 0}, // 2006
	[2]int{0, 0}, // 2007
	[2]int{0, 1}, // 2008
	[2]int{0, 0}, // 2009
	[2]int{0, 0}, // 2010
	[2]int{0, 0}, // 2011
	[2]int{1, 0}, // 2012
	[2]int{0, 0}, // 2013
	[2]int{0, 0}, // 2014
	[2]int{1, 0}, // 2015
	[2]int{0, 1}, // 2016
	[2]int{0, 0}, // 2017
	[2]int{0, 0}, // 2018
	[2]int{0, 0}, // 2019
	[2]int{0, 0}, // 2020
	[2]int{0, 0}, // 2021
}

func GetLeapSeconds(when time.Time) int {
	adj := 10
	year := LeapSecondsBaseYear
	for {
		if year > when.Year() {
			break
		}
		if year - LeapSecondsBaseYear >= len(LeapSecondsAdjustments) {
			if when.Year() - LeapSecondsBaseYear > len(LeapSecondsAdjustments) || when.Year() - LeapSecondsBaseYear == len(LeapSecondsAdjustments) && when.Month() > time.June {
				log.Printf("I don't know about leap seconds after %d", LeapSecondsBaseYear + len(LeapSecondsAdjustments))
				break
			}
		}
		entry := LeapSecondsAdjustments[year - LeapSecondsBaseYear]
		if year == when.Year() {
			if when.Month() > time.June {
				adj += entry[0]
			}
			break
		}
		adj += entry[0] + entry[1]
		year += 1
	}
	return adj
}

var DeltaT = [][]float64{
	[]float64{43.4724, 43.5648, 43.6737, 43.7782, 43.8763, 43.9562, 44.0315, 44.1132, 44.1982, 44.2952, 44.3936}, // 1973, starting from feb
	[]float64{44.4841, 44.5646, 44.6425, 44.7386, 44.8370, 44.9302, 44.9986, 45.0584, 45.1284, 45.2064, 45.2980, 45.3897}, // 1974
	[]float64{45.4761, 45.5633, 45.6450, 45.7375, 45.8284, 45.9133, 45.9820, 46.0408, 46.1067, 46.1825, 46.2789, 46.3713}, // 1975
	[]float64{46.4567, 46.5445, 46.6311, 46.7302, 46.8284, 46.9247, 46.9970, 47.0709, 47.1451, 47.2362, 47.3413, 47.4319}, // 1976
	[]float64{47.5214, 47.6049, 47.6837, 47.7781, 47.8771, 47.9687, 48.0348, 48.0942, 48.1608, 48.2460, 48.3439, 48.4355}, // 1977
	[]float64{48.5344, 48.6325, 48.7294, 48.8365, 48.9353, 49.0319, 49.1013, 49.1591, 49.2286, 49.3070, 49.4018, 49.4945}, // 1978
	[]float64{49.5862, 49.6805, 49.7602, 49.8556, 49.9489, 50.0347, 50.1019, 50.1622, 50.2260, 50.2968, 50.3831, 50.4599}, // 1979
	[]float64{50.5387, 50.6161, 50.6866, 50.7658, 50.8454, 50.9187, 50.9761, 51.0278, 51.0843, 51.1538, 51.2319, 51.3063}, // 1980
	[]float64{51.3808, 51.4526, 51.5160, 51.5985, 51.6809, 51.7573, 51.8133, 51.8532, 51.9014, 51.9603, 52.0328, 52.0985}, // 1981
	[]float64{52.1668, 52.2316, 52.2938, 52.3680, 52.4465, 52.5180, 52.5752, 52.6178, 52.6668, 52.7340, 52.8056, 52.8792}, // 1982
	[]float64{52.9565, 53.0445, 53.1268, 53.2197, 53.3024, 53.3747, 53.4335, 53.4778, 53.5300, 53.5845, 53.6523, 53.7256}, // 1983
	[]float64{53.7882, 53.8367, 53.8830, 53.9443, 54.0042, 54.0536, 54.0856, 54.1084, 54.1463, 54.1914, 54.2452, 54.2958}, // 1984
	[]float64{54.3427, 54.3911, 54.4320, 54.4898, 54.5456, 54.5977, 54.6355, 54.6532, 54.6776, 54.7174, 54.7741, 54.8253}, // 1985
	[]float64{54.8713, 54.9161, 54.9581, 54.9997, 55.0476, 55.0912, 55.1132, 55.1328, 55.1532, 55.1898, 55.2416, 55.2838}, // 1986
	[]float64{55.3222, 55.3613, 55.4063, 55.4629, 55.5111, 55.5524, 55.5812, 55.6004, 55.6262, 55.6656, 55.7168, 55.7698}, // 1987
	[]float64{55.8197, 55.8615, 55.9130, 55.9663, 56.0220, 56.0700, 56.0939, 56.1105, 56.1314, 56.1611, 56.2068, 56.2583}, // 1988
	[]float64{56.3000, 56.3399, 56.3790, 56.4283, 56.4804, 56.5352, 56.5697, 56.5983, 56.6328, 56.6739, 56.7332, 56.7972}, // 1989
	[]float64{56.8553, 56.9111, 56.9755, 57.0471, 57.1136, 57.1738, 57.2226, 57.2597, 57.3073, 57.3643, 57.4334, 57.5016}, // 1990
	[]float64{57.5653, 57.6333, 57.6973, 57.7711, 57.8407, 57.9058, 57.9576, 57.9975, 58.0426, 58.1043, 58.1679, 58.2389}, // 1991
	[]float64{58.3092, 58.3833, 58.4537, 58.5401, 58.6228, 58.6917, 58.7410, 58.7836, 58.8406, 58.8986, 58.9714, 59.0438}, // 1992
	[]float64{59.1218, 59.2003, 59.2747, 59.3574, 59.4434, 59.5242, 59.5850, 59.6344, 59.6928, 59.7588, 59.8386, 59.9111}, // 1993
	[]float64{59.9845, 60.0564, 60.1231, 60.2042, 60.2804, 60.3530, 60.4012, 60.4440, 60.4900, 60.5578, 60.6324, 60.7059}, // 1994
	[]float64{60.7853, 60.8664, 60.9387, 61.0277, 61.1103, 61.1870, 61.2454, 61.2881, 61.3378, 61.4036, 61.4760, 61.5525}, // 1995
	[]float64{61.6287, 61.6846, 61.7433, 61.8132, 61.8823, 61.9497, 61.9969, 62.0343, 62.0714, 62.1202, 62.1810, 62.2382}, // 1996
	[]float64{62.2950, 62.3506, 62.3995, 62.4754, 62.5463, 62.6136, 62.6571, 62.6942, 62.7383, 62.7926, 62.8567, 62.9146}, // 1997
	[]float64{62.9659, 63.0217, 63.0807, 63.1462, 63.2053, 63.2599, 63.2844, 63.2961, 63.3126, 63.3422, 63.3871, 63.4339}, // 1998
	[]float64{63.4673, 63.4979, 63.5319, 63.5679, 63.6104, 63.6444, 63.6642, 63.6739, 63.6926, 63.7147, 63.7518, 63.7927}, // 1999
	[]float64{63.8285, 63.8557, 63.8804, 63.9075, 63.9393, 63.9691, 63.9799, 63.9833, 63.9938, 64.0093, 64.0400, 64.0670}, // 2000
	[]float64{64.0908, 64.1068, 64.1282, 64.1584, 64.1833, 64.2094, 64.2117, 64.2073, 64.2116, 64.2223, 64.2500, 64.2761}, // 2001
	[]float64{64.2998, 64.3192, 64.3450, 64.3735, 64.3943, 64.4151, 64.4132, 64.4118, 64.4097, 64.4168, 64.4329, 64.4511}, // 2002
	[]float64{64.4734, 64.4893, 64.5053, 64.5269, 64.5471, 64.5597, 64.5512, 64.5371, 64.5359, 64.5415, 64.5544, 64.5654}, // 2003
	[]float64{64.5736, 64.5891, 64.6015, 64.6176, 64.6374, 64.6549, 64.6530, 64.6379, 64.6372, 64.6400, 64.6543, 64.6723}, // 2004
	[]float64{64.6876, 64.7052, 64.7313, 64.7575, 64.7811, 64.8001, 64.7995, 64.7876, 64.7831, 64.7921, 64.8096, 64.8311}, // 2005
	[]float64{64.8452, 64.8597, 64.8850, 64.9175, 64.9480, 64.9794, 64.9895, 65.0028, 65.0138, 65.0371, 65.0773, 65.1122}, // 2006
	[]float64{65.1464, 65.1833, 65.2145, 65.2494, 65.2921, 65.3279, 65.3413, 65.3452, 65.3496, 65.3711, 65.3972, 65.4296}, // 2007
	[]float64{65.4573, 65.4868, 65.5152, 65.5450, 65.5781, 65.6127, 65.6288, 65.6370, 65.6493, 65.6760, 65.7097, 65.7461}, // 2008
	[]float64{65.7768, 65.8025, 65.8237, 65.8595, 65.8973, 65.9323, 65.9509, 65.9534, 65.9628, 65.9839, 66.0147, 66.0420}, // 2009
	[]float64{66.0699, 66.0961, 66.1310, 66.1683, 66.2072, 66.2356, 66.2409, 66.2335, 66.2349, 66.2441, 66.2751, 66.3054}, // 2010
	[]float64{66.3246, 66.3406, 66.3624, 66.3957, 66.4289, 66.4619, 66.4749, 66.4751, 66.4829, 66.5056, 66.5383, 66.5706}, // 2011
	[]float64{66.6030, 66.6340, 66.6569, 66.6925, 66.7289, 66.7579, 66.7708, 66.7740, 66.7846, 66.8103, 66.8400, 66.8779}, // 2012
	[]float64{66.9069, 66.9443, 66.9763, 67.0258, 67.0716, 67.1100, 67.1266, 67.1331, 67.1458, 67.1718, 67.2091, 67.2460}, // 2013
	[]float64{67.2810, 67.3136, 67.3457, 67.3890}, // 2014
}

// returns a suitable value for delta_t for the given datetime.
func GetDeltaT(when time.Time) float64 {
	year := when.Year()
	month := int(when.Month())
	if year < DeltaTBaseYear {
		year = DeltaTBaseYear
		month = 1
	} else if year == DeltaTBaseYear {
		month = month - int(DeltaTBaseMonth)
		if month < 0 {
			month = 0
		}
		month += 1
	} else if year >= DeltaTBaseYear + len(DeltaT) {
		year = DeltaTBaseYear + len(DeltaT) - 1
	}
	if year == DeltaTBaseYear + len(DeltaT) - 1 {
		xmonth := len(DeltaT[year - DeltaTBaseYear])
		if month > xmonth {
			month = xmonth
		}
	}
	return DeltaT[year - DeltaTBaseYear][month - 1]
}

/*
returns the UT Julian day number (including fraction of a day) corresponding to
the specified date/time. This version assumes the proleptic Gregorian calendar;
trying to adjust for pre-Gregorian dates/times seems pointless when the changeover
happened over such wildly varying times in different regions.
*/
func GetJulianSolarDay(when time.Time) float64 {
	t := float64(when.UnixMicro()) / 1e6
	t += float64(GetLeapSeconds(when))
	t += float64(TtOffset)
	t -= GetDeltaT(when)
	return t / 86400.0 + GregorianDayOffset + JulianDayOffset
}

/*
returns the TT Julian day number (including fraction of a day) corresponding to
the specified date/time. This version assumes the proleptic Gregorian calendar;
trying to adjust for pre-Gregorian dates/times seems pointless when the changeover
happened over such wildly varying times in different regions.
*/
func GetJulianEphemerisDay(when time.Time) float64 {
	t := float64(when.UnixMicro()) / 1e6
	t += float64(GetLeapSeconds(when))
	t += float64(TtOffset)
	return t / 86400.0 + GregorianDayOffset + JulianDayOffset
}

func GetJulianCentury(julianDay float64) float64 {
	return (julianDay - 2451545.0) / 36525.0
}

func GetJulianEphemerisCentury(julianEphemerisDay float64) float64 {
	return (julianEphemerisDay - 2451545.0) / 36525.0
}

func GetJulianEphemerisMillenium(julianEphemerisCentury float64) float64 {
	return julianEphemerisCentury / 10
}
