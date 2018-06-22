package main

import (
	"math"
	"fmt"
)

func main() {
    res := algr(43.86169731617,125.311807394028,10,6371)
    fmt.Println(res)
}

type CoorBounding struct {
	min_lat  float64
	max_lat  float64
	min_lon  float64
	max_lon  float64
}

func (cb *CoorBounding) isunderlimitation(lat float64, lon float64) bool {
	if (cb.min_lat < lat && lat < cb.max_lat) && (cb.min_lon < lon && lon < cb.max_lon) {
		return true
	} else {
		return false
	}
}


func algr(la float64, lo float64, distance float64, radius float64) (cb *CoorBounding) {
	lat := la * math.Pi / 180
	lon := lo * math.Pi / 180
	rad_dist := distance / radius
	lat_min := lat - rad_dist
	lat_max := lat + rad_dist

	var lon_min float64
	var lon_max float64
	var lon_t float64
	if (lat_min > -math.Pi/2 && lat_max < math.Pi/2) {
		lon_t = math.Asin(math.Sin(rad_dist) / math.Cos(lat))
		lon_min = lon - lon_t

		if ( lon_min < -math.Pi) {
			lon_min = lon_min + 2 * math.Pi
			lon_max = lon + lon_t
		}
		if ( lon_max > math.Pi) {
			lon_max = lon_max -  2 * math.Pi
		}
	} else {
		lat_min = math.Max(lat_min, -math.Pi/2)
		lat_max = math.Min(lat_max, math.Pi/2)
		lon_min = -math.Pi
		lon_max = math.Pi
	}
	fin_lat_min := lat_min * 180 / math.Pi
	fin_lat_max := lat_max * 180 / math.Pi
	fin_lon_min := lon_min * 180 / math.Pi
	fin_lon_max := lon_max * 180 / math.Pi
	//fin_lon_max := (lo - fin_lon_min) + lo
	//fmt.Println(lon_max)
	obj := CoorBounding{min_lat:fin_lat_min, max_lat:fin_lat_max, min_lon:fin_lon_min, max_lon:fin_lon_max}
	return &obj
}
