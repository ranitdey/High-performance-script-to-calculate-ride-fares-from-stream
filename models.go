package main

type Point struct {
	rideID    int
	lat       float64
	lng       float64
	timestamp int64
}

type Trip struct {
	kmsNight   float64
	kmsDay    float64
	idleHours  float64
}

type Fare struct {
	rideID int
	fare   float64
}
