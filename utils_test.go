package main

import (
	"math"
	"testing"
)

func TestU(t *testing.T) {
	tests := []struct {
		desc     string
		a        Point
		b        Point
		expected float64
	}{
		{
			"This dataset will produce high speed",
			Point{lat: 37.944253, lng: 23.758287, timestamp: 1405591400},
			Point{lat: 37.944028, lng: 23.758845, timestamp: 1405591401},
			197.8415568911867,
		},
		{
			"This dataset will produce low speed",
			Point{lat: 37.946545, lng: 23.754918, timestamp: 1405591065},
			Point{lat: 37.946413, lng: 23.754767, timestamp: 1405591094},
			2.4538892024620274,
		},

	}
	for _, test := range tests {
		spd := u(test.a,test.b)
		if spd != test.expected{
			t.Error("Calculated speed is not expected: ",spd)
		}
	}
}

func TestDs(t *testing.T) {
	tests := []struct {
		name     string
		a        Point
		b        Point
		expected float64
	}{
		{
			"Calculate distance between two points",
			Point{lat: 37.944440, lng: 23.758473},
			Point{lat: 37.944130, lng: 23.758447},
			0.03454574361591293,
		},
	}
	for _, test := range tests {
		d := ds(test.a,test.b)
		if d != test.expected{
			t.Error("Calculated distance is not expected: ",d)
		}
	}

}

func TestDt(t *testing.T) {
	tests := []struct {
		name     string
		a        Point
		b        Point
		expected float64
	}{
		{
			"Calculate dt when first timestamp >  second timestamp",
			Point{timestamp: 1405594957},
			Point{timestamp: 1405595203},
			0.06833333333333333,
		},
		{
			"Calculate dt when second timestamp >  first timestamp",
			Point{timestamp: 1405595203},
			Point{timestamp: 1405594957},
			0.06833333333333333,
		},
	}
	for _, test := range tests {
		time := dt(test.a,test.b)
		if time != test.expected{
			t.Error("Calculated time difference is not expected: ",time)
		}
	}
}

func TestFare(t *testing.T) {
	tests := []struct {
		name     string
		trip     Trip
		expected float64
	}{
		{
			"Test minimum fare",
			Trip{idleHours: 0, kmsDay: 0, kmsNight: 0},
			3.47,
		},
		{
			"Test a day trip",
			Trip{idleHours: 2, kmsDay: 20, kmsNight: 0},
			1.30 + 11.90*2 + 0.74*20 + 1.30*0,
		},
		{
			"Test a night trip",
			Trip{idleHours: 4, kmsDay: 0, kmsNight: 4},
			1.30 + 11.90*4 + 0.74*0 + 1.30*4,
		},
		{
			"Test a day-night trip",
			Trip{idleHours: 0, kmsDay: 5, kmsNight: 5},
			1.30 + 11.90*0 + 0.74*5 + 1.30*5,
		},
		{
			"Test a idle trip",
			Trip{idleHours: 9, kmsDay: 0, kmsNight: 0},
			1.30 + 11.90*9 + 0.74*0 + 1.30*0,
		},

	}
	for _, test := range tests {
		f := fare(test.trip)
		threshold := 0.0001
		diff := math.Abs(test.expected - f)
		if diff>threshold{
			t.Error("Calculated fare differed more than the threshold: ",f, test.expected)
		}
	}

}
func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		a        float64
		b        float64
		expected float64
	}{
		{
			"Test a>b",
			20,
			1,
			20,
		},
		{
			"Test b>a",
			12,
			112,
			112,
		},
		{
			"Test a=b",
			-87,
			-87,
			-87,
		},
		{
			"Test positive and negative value",
			-87,
			2,
			2,
		},
	}
	for _, test := range tests {
		m := max(test.a,test.b)
		if m != test.expected{
			t.Error("Calculated max is not expected: ",m)
		}
	}
}

