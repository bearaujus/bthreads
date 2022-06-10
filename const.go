package bthreads

import "time"

const (
	appName    = "BThreads"
	appVersion = "1.2"
)

const (
	dName                = "An bthreads instance"
	dFuncGoroutinesCount = 1
	dGoroutinesDelay     = 0
	dLogDelay            = 60 * time.Millisecond
	dStartDelay          = 3 * time.Second
)
