package main

import (
	"io"
	"time"
)

type Reader interface {
	Read() (record []string, err error)
}

func emitStructuredRecords(r Reader) <-chan Point {
	out := make(chan Point, channelBufferSize)
	go func() {
		for {
			data, err := r.Read()
			if err == io.EOF {
				break
			}
			checkError("Error in reading line ", err)
			point, err := serialize(data)
			if err == nil {
				out <- *point
			} else {
				checkError("Error while serializing line: ", err)
			}
		}
		close(out)
	}()
	return out
}

func groupUniqueRides(points <-chan Point) <-chan chan Point {
	rides := make(chan chan Point, channelBufferSize)

	go func() {
		previous, ok := <-points
		if !ok {
			close(rides)
			return
		}
		ride := make(chan Point, channelBufferSize)
		rides <- ride
		ride <- previous
		for current := range points {
			if previous.rideID != current.rideID {
				close(ride)
				ride = make(chan Point, channelBufferSize)
				rides <- ride
			}
			ride <- current
			previous = current
		}
		close(ride)
		close(rides)
	}()
	return rides
}

func filterInvalidPoints(positions <-chan Point) <-chan Point {
	out := make(chan Point, channelBufferSize)

	go func() {
		previous, err := <-positions
		if !err {
			close(out)
			return
		}
		out <- previous
		for current := range positions {
			if u(current, previous) <= 100 {
				out <- current
			}
			previous = current
		}
		close(out)
	}()
	return out
}

func estimateFare(positions <-chan Point) <-chan Fare {
	trip := Trip{
		idleHours: 0,
		kmsNight:  0,
		kmsDay:    0,
	}

	out := make(chan Fare, 1)
	go func() {
		previous, ok := <-positions
		if !ok {
			close(out)
			return
		}
		for current := range positions {
			if u(previous, current) <= maxSpeedIdle {
				trip.idleHours += dt(current, previous)
			} else {
				d := time.Unix(current.timestamp, 0).UTC()
				h := d.Hour()
				if h > nightHourStart && h <= nightHourEnd {
					trip.kmsNight += ds(current, previous)
				} else {
					trip.kmsDay += ds(current, previous)
				}
			}
			previous = current
		}

		out <- Fare{
			rideID: previous.rideID,
			fare:   fare(trip),
		}
		close(out)
	}()

	return out
}

