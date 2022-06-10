package bthreads

import (
	"sync"
	"time"
)

type Config struct {
	Name                string
	FuncGoroutinesCount int
	GoroutinesDelay     time.Duration
	LogDelay            time.Duration
	StartDelay          time.Duration
	HideWorkersData     bool
}

type instance struct {
	startTime time.Time

	name                string
	funcGoroutinesCount int
	goroutinesDelay     time.Duration
	logDelay            time.Duration
	startDelay          time.Duration
	hideWorkerData      bool

	numIter        int
	numIterSuccess int
	numIterFail    int

	funcs        []func() bool
	instanceData sync.Map
}

type instanceData struct {
	num        int
	numSuccess int
	numFail    int
}
