package main

import (
	"github.com/umahmood/haversine"
	"math"
)

func u(a Point, b Point) float64 {
	return ds(a, b) / dt(a, b)
}

func ds(a Point, b Point) float64 {
	c1 := haversine.Coord{Lat: b.lat, Lon: b.lng}
	c2 := haversine.Coord{Lat: a.lat, Lon: a.lng}
	_, ds := haversine.Distance(c1, c2)
	return ds
}

func dt(a Point, b Point) float64 {
	return math.Abs(float64(a.timestamp-b.timestamp)) / 3600
}

func fare(trip Trip) float64 {
	fare := flagCost + (trip.idleHours * idleHourCost) + (trip.kmsNight * kmNightCost) + (trip.kmsDay * kmDefaultCost)
	return max(fare, minCost)
}

func max(a float64, b float64) float64 {
	if a < b {
		return b
	} else {
		return a
	}
}
