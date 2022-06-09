package bthreads

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/inancgumus/screen"
)

const (
	appName    = "BThreads"
	appVersion = "1.0"

	defaultName             = "An bthreads instance"
	defaultFuncsThreadCount = 1
	defaultThreadsDelay     = 100000
	defaultLogsDelay        = 30 * time.Millisecond
	defaultStartDelay       = 3 * time.Second
)

func New(param *Threads) *threads {
	return &threads{
		startTime:        time.Now().Local(),
		name:             param.Name,
		funcsThreadCount: param.FuncsThreadCount,
		threadsDelay:     param.ThreadsDelay,
		logsDelay:        param.LogsDelay,
		startDelay:       param.StartDelay,
	}
}

type Threads struct {
	Name             string
	FuncsThreadCount int
	ThreadsDelay     time.Duration
	LogsDelay        time.Duration
	StartDelay       time.Duration
}

type threads struct {
	startTime time.Time

	name             string
	funcsThreadCount int
	threadsDelay     time.Duration
	logsDelay        time.Duration
	startDelay       time.Duration

	numIter        int
	numIterSuccess int
	numIterFail    int

	funcs       []func() bool
	threadsData sync.Map
}

type threadData struct {
	num        int
	numSuccess int
	numFail    int
}

func (st *threads) AddFunc(f func() bool) {
	st.funcs = append(st.funcs, f)
}

func (st *threads) AddFuncs(funcs ...func() bool) {
	st.funcs = append(st.funcs, funcs...)
}

func (st *threads) Start() {
	fmt.Println(color.HiBlueString("[ %v ] ", time.Now().Local().Format(time.RFC1123)) +
		color.HiYellowString("[ %v v%v ][ BETA ]", appName, appVersion))
	fmt.Println(color.HiWhiteString("  Verifying instance ..."))
	if !st.verifyInstance() {
		return
	}
	fmt.Println(color.HiWhiteString("  Instance Verified !"))
	fmt.Println(color.HiWhiteString("  Starting instance ..."))
	<-time.After(st.startDelay)
	fmt.Println(color.HiWhiteString("  Instance Started !"))
	st.runThreads()
	st.runLogger()
}

func (st *threads) verifyInstance() bool {
	if st.name == "" {
		st.name = defaultName
	}
	if st.funcsThreadCount == 0 {
		st.funcsThreadCount = defaultFuncsThreadCount
	}
	if st.threadsDelay == 0 {
		st.threadsDelay = defaultThreadsDelay
	}
	if st.logsDelay == 0 {
		st.logsDelay = defaultLogsDelay
	}
	if st.startDelay == 0 {
		st.startDelay = defaultStartDelay
	}

	if st.funcsThreadCount < 0 {
		fmt.Println(color.RedString("[err] ") + color.HiWhiteString("Invalid 'FuncsThreadCount'"))
		return false
	}
	if st.threadsDelay < 0 {
		return false
	}
	if st.logsDelay < 0 {
		return false
	}
	if st.startDelay < 0 {
		return false
	}
	return true
}

func (st *threads) runThreads() {
	var threadID int
	for i := 0; i < st.funcsThreadCount; i++ {
		for _, f := range st.funcs {
			threadID++
			st.threadsData.Store(threadID, &threadData{})
			go st.runThread(threadID, f)
		}
	}
}

func (st *threads) runThread(threadID int, f func() bool) {
	for {
		st.sync(threadID, f())
		<-time.After(st.threadsDelay)
	}
}

func (st *threads) runLogger() {
	for {
		screen.Clear()

		// Title
		// color.HiGreenString("%v ", st.startTime.Format(time.RFC1123))
		fmt.Println(color.HiYellowString("[ %v v%v ][ BETA ]", appName, appVersion))

		fmt.Println(color.HiWhiteString("  " + st.name))

		fmt.Println()

		// Instance
		diffTime := time.Time{}.Add(time.Now().Local().Sub(st.startTime))
		ths := float64(st.numIter) / float64(diffTime.Second())
		sr := (float64(st.numIterSuccess) / float64(st.numIter)) * 100

		fmt.Println(color.HiYellowString("[ Instance ]"))

		fmt.Println(color.HiMagentaString("  Time Elapsed\t") +
			color.HiWhiteString(diffTime.Format("15:04:05.00")))

		fmt.Println()

		fmt.Println(color.HiMagentaString("  Iter Speed\t") +
			color.HiWhiteString(fmt.Sprintf("%.2f it/s", ths)))

		fmt.Println(color.HiMagentaString("  Success Rate\t") +
			color.HiGreenString(fmt.Sprintf("%.2f", sr)+" %"))

		fmt.Println()

		fmt.Println(color.HiMagentaString("  Success\t") +
			color.HiGreenString(fmt.Sprint(st.numIterSuccess)+" it"))

		fmt.Println(color.HiMagentaString("  Fail\t\t") +
			color.RedString(fmt.Sprint(st.numIterFail)+" it"))

		fmt.Println(color.HiMagentaString("  Total\t \t") +
			color.HiWhiteString(fmt.Sprint(st.numIter)+" it"))

		fmt.Println()

		// Workers
		fmt.Println(color.HiYellowString("[ Worker ]"))

		mlog := make(map[int]string)
		fmt.Println(color.HiMagentaString("  Thread ID\tSuccess\t\tFail\t\tTotal"))
		st.threadsData.Range(func(key, value interface{}) bool {
			threadID, _ := key.(int)
			td, _ := value.(*threadData)
			r := color.HiWhiteString("  Thread-%v\t", threadID) +
				color.HiGreenString("%v\t\t", td.numSuccess) +
				color.RedString("%v\t\t", td.numFail) +
				color.HiWhiteString("%v", td.num)
			mlog[key.(int)] = r
			return true
		})
		for i := 0; i < len(mlog); i++ {
			fmt.Println(mlog[i+1])
		}
		fmt.Println()

		<-time.After(st.logsDelay)
	}
}

func (st *threads) sync(threadID int, isSuccess bool) {
	st.syncCount(threadID, isSuccess)
	st.syncData(threadID, isSuccess)
}

func (st *threads) syncCount(threadID int, isSuccess bool) {
	st.numIter++
	if isSuccess {
		st.numIterSuccess++
	} else {
		st.numIterFail++
	}
}

func (st *threads) syncData(threadID int, isSuccess bool) {
	rawtd, _ := st.threadsData.LoadAndDelete(threadID)
	td := rawtd.(*threadData)
	td.num++
	if isSuccess {
		td.numSuccess++
	} else {
		td.numFail++
	}
	st.threadsData.Store(threadID, td)
}
