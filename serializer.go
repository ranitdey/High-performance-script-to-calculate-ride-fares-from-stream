package main

import (
	"errors"
	"strconv"
)

func serialize(record []string) (*Point, error) {
	if len(record) != 4 {
		return nil, errors.New("A record doesn't contain exactly 4 fields")
	}
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return nil, err
	}
	lat, err := strconv.ParseFloat(record[1], 64)
	if err != nil {
		return nil, err
	}
	lng, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return nil, err
	}
	timestamp, err := strconv.ParseInt(record[3], 10, 64)
	if err != nil {
		return nil, err
	}
	return &Point{id, lat, lng, timestamp}, nil
}
