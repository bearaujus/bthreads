package bthreads

func (st *instance) sync(grID int, isSuccess bool) {
	st.syncCount(grID, isSuccess)
	st.syncData(grID, isSuccess)
}

func (st *instance) syncCount(grID int, isSuccess bool) {
	st.numIter++
	if isSuccess {
		st.numIterSuccess++
	} else {
		st.numIterFail++
	}
}

func (st *instance) syncData(grID int, isSuccess bool) {
	rawtd, _ := st.instanceData.LoadAndDelete(grID)
	td := rawtd.(*instanceData)
	td.num++
	if isSuccess {
		td.numSuccess++
	} else {
		td.numFail++
	}
	st.instanceData.Store(grID, td)
}
