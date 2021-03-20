package main

import (
	"io"
	"testing"
)

func TestEmitStructuredRecords(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]string
		expected []Point
	}{
		{
			"Test with same datasets",
			[][]string{
				{"1", "37.966660", "23.728308", "1405594957"},
				{"1", "38.966660", "23.728308", "1405595957"},
				{"1", "39.966660", "23.728308", "1405596957"},
			},
			[]Point{
				{1, 37.966660, 23.728308, 1405594957},
				{1, 38.966660, 23.728308, 1405595957},
				{1, 39.966660, 23.728308, 1405596957},
			},
		},

	}
	for _, test := range tests {
		reader := &MockReader{
			curr:  0,
			lines: test.input,
		}
		points := emitStructuredRecords(reader)
		c := 0
		for p := range points {
			if p != test.expected[c]{
				t.Error("Serialization output is invalid:  ",p)
			}
			c++
		}
	}
}


func TestGroupUniqueRides(t *testing.T) {
	tests := []struct {
		name     string
		input    []Point
		expected [][]Point
	}{
		{
			"Test empty input",
			[]Point{
			},
			[][]Point{
			},
		},
		{
			"Test grouping points in a single ride",
			[]Point{
				{1, 37.937462, 23.634823, 1405596109},
				{1, 37.938193, 23.633550, 1405596117},
				{1, 37.936942, 23.628098, 1405596186},
				{1, 37.935490, 23.625655, 1405596220},
			},
			[][]Point{
				{
					{1, 37.937462, 23.634823, 1405596109},
					{1, 37.938193, 23.633550, 1405596117},
					{1, 37.936942, 23.628098, 1405596186},
					{1, 37.935490, 23.625655, 1405596220},
				},
			},
		},
		{
			"Test grouping points in multiple rides",
			[]Point{
				{1, 37.935597, 23.625688, 1405596212},
				{2, 37.946832, 23.755435, 1405591553},
				{2, 37.948515, 23.746622, 1405591645},
				{2, 37.965268, 23.733788, 1405591906},
				{3, 37.964710, 23.741237, 1405593257},
				{5, 37.922093, 23.931009, 1405591252},
			},
			[][]Point{
				{
					{1, 37.935597, 23.625688, 1405596212},
				},
				{
					{2, 37.946832, 23.755435, 1405591553},
					{2, 37.948515, 23.746622, 1405591645},
					{2, 37.965268, 23.733788, 1405591906},
				},
				{
					{3, 37.964710, 23.741237, 1405593257},
				},
				{
					{5, 37.922093, 23.931009, 1405591252}},
			},
		},
	}
	for _, test := range tests {
		points := make(chan Point)
		go func() {
			for _, p := range test.input {
				points <- p
			}
			close(points)
		}()
		rides := groupUniqueRides(points)
		i := 0
		for ride := range rides {
			j := 0
			for p := range ride {
				if p!= test.expected[i][j]{
					t.Error("Ride grouping output is invalid:  ",p)
				}
				j++
			}
			i++
		}
	}
}

func TestFilterInvalidPoints(t *testing.T) {
	tests := []struct {
		name     string
		input    []Point
		expected []Point
	}{
		{
			"Test empty input",
			[]Point{
			},
			[]Point{
			},
		},
		{
			"Test happy path with no filter",
			[]Point{
				{5,37.931423,23.939967,1405591315},
				{5,37.932744,23.941293,1405591325},
				{5,37.933906,23.942595,1405591336},
				{5,37.934419,23.944129,1405591346},
			},
			[]Point{
				{5,37.931423,23.939967,1405591315},
				{5,37.932744,23.941293,1405591325},
				{5,37.933906,23.942595,1405591336},
				{5,37.934419,23.944129,1405591346},
			},
		},
		{
			"Test inputs with high speed",
			[]Point{
				{1,37.944253,23.758287,1405591400},
				{2,37.944028,23.758845, 1405591401},
			},
			[]Point{
				{1,37.944253,23.758287,1405591400},
			},
		},
	}
	for _, test := range tests {
		points := make(chan Point)
		go func() {
			for _, p := range test.input {
				points <- p
			}
			close(points)
		}()
		filtered := filterInvalidPoints(points)
		i := 0
		for p := range filtered {
			if p != test.expected[i]{
				t.Error("Ride filtering output is invalid:  ",p)
			}
			i++
		}
	}
}
func TestEstimateFare(t *testing.T) {
	tests := []struct {
		name     string
		input    []Point
		expected Fare
	}{
		{
			"Test empty input",
			[]Point{
			},
			Fare{rideID: 0, fare: 0},
		},
		{
			"Test fare of a day ride",
			[]Point{
				{5,37.858163,23.840420,1405590552},
				{5,37.937010,23.947143,1405591490},
			},
			Fare{rideID: 5, fare: 10.792767516358326},
		},
		{
			"Test fare of a night ride",
			[]Point{
				{5,37.858163,23.840420,1614310652},
				{5,37.937010,23.947143,1614311552},
			},
			Fare{rideID: 5, fare: 17.976483474683548},
		},
		{
			"Test minimum fare",
			[]Point{
				{5,37.858163,23.840420,1405590552},
				{5,37.858163,23.840420,1405590652},
			},
			Fare{rideID: 5, fare: 3.47},
		},

	}
	for _, test := range tests {
		points := make(chan Point)
		go func() {
			for _, p := range test.input {
				points <- p
			}
			close(points)
		}()
		fare := estimateFare(points)
		for f := range fare{
			if test.expected != f{
				t.Error("Ride fare estimation is invalid:  ",f)
			}
		}
	}
}

type MockReader struct {
	curr     int
	lines [][]string
}

func (r *MockReader) Read() (record []string, err error) {
	if r.curr >= len(r.lines) {
		return nil, io.EOF
	}
	current := r.lines[r.curr]
	r.curr ++
	return current, nil
}