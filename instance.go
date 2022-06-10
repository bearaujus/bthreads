package bthreads

import "time"

func (st *instance) runAllFunc() {
	var grID int
	for i := 0; i < st.funcGoroutinesCount; i++ {
		for _, f := range st.funcs {
			grID++
			st.instanceData.Store(grID, &instanceData{})
			go st.runFunc(grID, f)
		}
	}
}

func (st *instance) runFunc(grID int, f func() bool) {
	for {
		st.sync(grID, f())
		<-time.After(st.goroutinesDelay)
	}
}
