package main

import (
	"testing"
)

func TestReadDataAndCalculateFares(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Fare
	}{
		{
			"Test empty file",
			"./resource/test_data/empty.csv",
			[]Fare{
			},
		},
		{
			"Test happy file",
			"./resource/paths.csv",
			[]Fare{
				{1, 11.5252795303018},
				{2, 13.099227812467959},
				{3, 35.07100192476959},
				{4, 3.47},
				{5, 22.783577234160713},
				{6, 9.416193510354596},
				{7, 31.42153572253091},
				{8, 9.251912312410514},
				{9, 6.347124165556508},
			},
		},
		{
			"Test invalid file extension",
			"./resource/test_data/input.yaml",
			[]Fare{
			},
		},
		{
			"Test corrupted file",
			"./resource/test_data/corrupted.csv",
			[]Fare{
				{5, 3.47},
			},
		},

	}
	for _, tt := range tests {
		fares := readDataAndCalculateFares(tt.input)
		i := 0
		for f := range fares {
			if f!=tt.expected[i]{
				t.Error("Ride fare estimation is invalid:  ",f)
			}
			i++
		}
	}
}