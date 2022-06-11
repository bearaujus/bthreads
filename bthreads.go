package bthreads

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Bearaujus/bthreads/pkg/util"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

// App Const
const (
	appName    = "BThreads"
	appVersion = "1.6"
)

// Default Const
const (
	dName                = "An bthreads instance"
	dFuncGoroutinesCount = 1
	dGoroutinesDelay     = 0
	dLogDelay            = 60 * time.Millisecond
	dStartDelay          = 3 * time.Second
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

	numIter        int64
	numIterSuccess int64
	numIterFail    int64

	funcs        []func() bool
	instanceData sync.Map
}

type instanceData struct {
	num        int64
	numSuccess int64
	numFail    int64
	fIndex     int
}

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

func validateInstance(param *Config) (*Config, error) {
	// Verify Instance
	if param.FuncGoroutinesCount < 0 {
		return nil, errors.New("'Config.FuncGoroutinesCount' cannot < 0")
	}
	if param.GoroutinesDelay < 0 {
		return nil, errors.New("'Config.GoroutinesDelay' cannot < 0")
	}
	if param.LogDelay < 0 {
		return nil, errors.New("'Config.LogDelay' cannot < 0")
	}
	if param.StartDelay < 0 {
		return nil, errors.New("'Config.StartDelay' cannot < 0")
	}

	// Set Default
	if param.Name == "" {
		param.Name = dName
	}
	if param.FuncGoroutinesCount == 0 {
		param.FuncGoroutinesCount = dFuncGoroutinesCount
	}
	if param.GoroutinesDelay == 0 {
		param.GoroutinesDelay = dGoroutinesDelay
	}
	if param.LogDelay == 0 {
		param.LogDelay = dLogDelay
	}
	if param.StartDelay == 0 {
		param.StartDelay = dStartDelay
	}

	return param, nil
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

func (st *instance) runAllFunc() {
	var grID int

	// Iterate over gouroutines count
	for i := 0; i < st.funcGoroutinesCount; i++ {

		// Iterate over input func
		for fIndex, f := range st.funcs {
			grID++

			// Init instance data with goroutineID
			st.instanceData.Store(grID, &instanceData{
				fIndex: fIndex,
			})

			// Run and put the func into goruoutines
			go st.runFunc(grID, f)
		}
	}
}

func (st *instance) runFunc(grID int, f func() bool) {
	for {
		// Run inputed func
		res := f()

		// Sync count data for monitoring
		st.syncCount(grID, res)

		// If worker data is hided
		if !st.hideWorkerData {

			// Sync worker data for monitoring
			st.syncData(grID, res)
		}

		// Make delay relative to goroutinesDelay
		<-time.After(st.goroutinesDelay)
	}
}

func (st *instance) log() {
	for {
		// Clear screen
		util.ClearScreen()

		// Print simple log if hideWorkerData is true
		st.printSimpleLog()
		if st.hideWorkerData {
			continue
		}

		// Print advanced log i hideWorkerData is false
		st.printWorkerLog()

		// Make delay relative to logDelay
		<-time.After(st.logDelay)
	}
}

func (st *instance) runStartingLog() {
	// Print condition when startDelay < time.Second
	if st.startDelay < time.Second {
		util.ClearPrint(color.HiWhiteString("Starting ") + color.HiYellowString(appName) + color.HiWhiteString(" instance..."))
		<-time.After(st.startDelay)
		return
	}

	// Print condition when startDelay > time.second
	for i := int(st.startDelay / time.Second); i > 0; i-- {
		util.ClearPrint(color.HiYellowString(appName) + color.HiWhiteString(fmt.Sprintf(" instance will start on %d", i)))
		<-time.After(time.Second)
	}
}

func (st *instance) getHeader() string {
	return strings.Join([]string{
		// Header
		color.HiYellowString("[ %v ", appName) + color.HiWhiteString("v%v", appVersion) + color.HiYellowString(" ]"),
		color.HiWhiteString("  " + st.name),
		"",
	}, "\n")
}

func (st *instance) printSimpleLog() {
	// Calculate time elapsed
	td := time.Now().Local().Sub(st.startTime)

	// Calculate iter speed
	ths := float64(st.numIter) / float64(td.Seconds())

	// Calculate success rate
	sr := (float64(st.numIterSuccess) / float64(st.numIter)) * float64(100)

	// Calculate fail rate
	fr := (float64(st.numIterFail) / float64(st.numIter)) * float64(100)

	fmt.Println(strings.Join([]string{
		st.getHeader(),

		// Instance
		color.HiYellowString("[ Instance ]"),
		color.HiWhiteString("  Time Elapsed\t") + color.HiCyanString(td.String()),
		"",
		color.HiWhiteString("  Total Iter\t") + color.HiWhiteString("%v it", st.numIter),
		color.HiWhiteString("  Iter Speed\t") + color.HiCyanString(fmt.Sprintf("%.2f it/s", ths)),
		"",
		color.HiWhiteString("  Success Rate\t") + color.HiGreenString("%.2f %v", sr, "%") + color.HiWhiteString(" %v it", st.numIterSuccess),
		color.HiWhiteString("  Fail Rate\t") + color.RedString("%.2f %v", fr, "%") + color.HiWhiteString(" %v it", st.numIterFail),
		"",
	}, "\n"))
}

func (st *instance) printWorkerLog() {
	// Initialize output for worker data
	wd := make(map[int]map[int][]string)

	// Iterate over sync.Map
	st.instanceData.Range(func(key, value interface{}) bool {

		// Casting sync.Map interfaces
		grID, _ := key.(int)
		td, _ := value.(*instanceData)

		// Add data to worker data
		wd[grID] = make(map[int][]string)
		wd[grID][td.fIndex] = []string{
			color.HiCyanString("gr-%v", grID),
			color.HiWhiteString("%v", td.fIndex),
			color.HiGreenString("%v", td.numSuccess),
			color.RedString("%v", td.numFail),
			color.HiWhiteString("%v", td.num),
		}

		return true
	})

	// Initialize output table
	ot := tablewriter.NewWriter(os.Stdout)
	ot.SetHeader([]string{"Goroutine", "FIndex", "Success", "Fail", "Total"})
	ot.SetAlignment(tablewriter.ALIGN_LEFT)
	ot.SetAutoFormatHeaders(false)

	// Sort workers data
	for i := 0; i < len(st.funcs); i++ {
		for j := 0; j < len(wd); j++ {
			// Add worker data to output table
			ot.Append(wd[j+1][i])
		}
	}

	// Workers
	fmt.Println(color.HiYellowString("[ Workers ]"))
	ot.Render()
	fmt.Println()
}

func (st *instance) syncCount(grID int, isSuccess bool) {
	// Each syncCount() called, insert num iter
	st.numIter++

	// If isSuccess is true, add numIterSuccess
	if isSuccess {
		st.numIterSuccess++
		return
	}

	// If isSuccess is false, add numIterFail
	st.numIterFail++
}

func (st *instance) syncData(grID int, isSuccess bool) {
	// Load and delete gouroutine data by their ID
	rawgr, _ := st.instanceData.LoadAndDelete(grID)

	// Casting raw goroutine data to instance data
	gr := rawgr.(*instanceData)

	// Update selected goroutine total called num
	gr.num++

	// If this iterate indicates success, increase this goroutine log numSuccess otherwise increase numFail
	if isSuccess {
		gr.numSuccess++
	} else {
		gr.numFail++
	}

	// Store the updated data
	st.instanceData.Store(grID, gr)
}
