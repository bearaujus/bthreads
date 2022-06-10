package bthreads

import (
	"math"
	"time"
)

func NewInstance(param *Config) (*instance, error) {
	// Validate the instance
	if config, err := validateInstance(param); err != nil {
		return nil, err
	} else {
		param = config
	}

	// Return the instance
	return &instance{
		name:                param.Name,
		funcGoroutinesCount: param.FuncGoroutinesCount,
		goroutinesDelay:     param.GoroutinesDelay,
		logDelay:            param.LogDelay,
		startDelay:          param.StartDelay,
		hideWorkerData:      param.HideWorkersData,
	}, nil
}

func (st *instance) AddFunc(f func() bool) {
	// Add a func to the instance
	st.funcs = append(st.funcs, f)
}

func (st *instance) AddFuncs(funcs ...func() bool) {
	// Add funcs to the instance
	st.funcs = append(st.funcs, funcs...)
}

func (st *instance) Start() {
	// Run starting message
	st.runStartingLog()

	// Set start time
	st.startTime = time.Now().Local()

	// Run all funcs into goroutines
	st.runAllFunc()

	// Run instance logger
	go st.log()

	// Run this instance forever (until ^c)
	time.Sleep(time.Duration(math.MaxInt64))
}
