package main

import "time"

type Bandwidth float32

type SpeedTestResult struct {
	Timestamp time.Time     `bson:"timestamp"` // timestamp at which the bandwidth has been monitored
	Ping      time.Duration `bson:"ping"`      // response time (latency)
	Download  Bandwidth     `bson:"download"`  // downstream in Mbit/s
	Upload    Bandwidth     `bson:"upload"`  // upstream in Mbit/s
}

type SpeedTester interface {
	MeasureBandwidth() (*SpeedTestResult, error)
}
