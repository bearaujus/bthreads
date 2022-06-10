package bthreads

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
)

func (st *instance) log() {
	for {
		if st.hideWorkerData {
			print(st.getSimpleLog())
		} else {
			print(st.getAdvancedLog())
		}
		<-time.After(st.logDelay)
	}
}

func (st *instance) runStartingLog() {
	if st.startDelay < time.Second {
		print(color.HiWhiteString("Starting ") + color.HiYellowString(appName) + color.HiWhiteString(" instance..."))
		<-time.After(st.startDelay)
	} else {
		for i := int(st.startDelay / time.Second); i > 0; i-- {
			print(color.HiYellowString(appName) + color.HiWhiteString(fmt.Sprintf(" instance will start on %d", i)))
			<-time.After(time.Second)
		}
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

func (st *instance) getSimpleLog() string {
	diffTime := time.Time{}.Add(time.Now().Local().Sub(st.startTime))
	ths := float64(st.numIter) / float64(diffTime.Second())
	sr := (float64(st.numIterSuccess) / float64(st.numIter)) * 100
	fr := 100.0 - sr

	return strings.Join([]string{
		st.getHeader(),

		// Instance
		color.HiYellowString("[ Instance ]"),
		color.HiWhiteString("  Time Elapsed\t") + color.HiWhiteString(diffTime.Format("15:04:05.00")),
		color.HiWhiteString("  Iter Speed\t") + color.HiWhiteString(fmt.Sprintf("%.2f it/s", ths)),
		"",
		color.HiWhiteString("  Success Rate\t") + color.HiGreenString(fmt.Sprintf("%.2f", sr)+" %"),
		color.HiWhiteString("  Fail Rate\t") + color.MagentaString(fmt.Sprintf("%.2f", fr)+" %"),
		"",
		color.HiWhiteString("  Success\t") + color.HiGreenString(fmt.Sprint(st.numIterSuccess)+" it"),
		color.HiWhiteString("  Fail\t\t") + color.RedString(fmt.Sprint(st.numIterFail)+" it"),
		color.HiWhiteString("  Total\t \t") + color.HiWhiteString(fmt.Sprint(st.numIter)+" it"),
		"",
	}, "\n")
}

func (st *instance) getAdvancedLog() string {
	mlog := make(map[int]string)
	st.instanceData.Range(func(key, value interface{}) bool {
		grID, _ := key.(int)
		td, _ := value.(*instanceData)
		r := color.HiCyanString("  gr-%v\t\t", grID) +
			color.HiGreenString("%v\t\t", td.numSuccess) +
			color.RedString("%v\t\t", td.numFail) +
			color.HiWhiteString("%v", td.num)
		mlog[key.(int)] = r
		return true
	})
	workersTableData := []string{}
	for i := 0; i < len(mlog); i++ {
		workersTableData = append(workersTableData, mlog[i+1])
	}
	workersTable := strings.Join(workersTableData, "\n")
	return strings.Join([]string{
		st.getSimpleLog(),

		// Workers
		color.HiYellowString("[ Workers ]"),
		color.HiWhiteString("  Goroutine\tSuccess\t\tFail\t\tTotal"),
		workersTable,
		"",
	}, "\n")
}

var (
	cmdFunc = map[string]func(){
		"linux": func() {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
		"darwin": func() {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
		"windows": func() {
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		},
	}
)

func clear() {
	if f, ok := cmdFunc[runtime.GOOS]; !ok {
		panic("Your platform is unsupported!")
	} else {
		f()
	}
}

func print(param string) {
	clear()
	fmt.Printf("\n%v\n", param)
}
