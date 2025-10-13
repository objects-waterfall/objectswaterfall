package enums

type WorkerStopStatus int

const (
	_                  WorkerStopStatus = iota
	StoppedByTimer     WorkerStopStatus = iota
	StoppedByUser      WorkerStopStatus = iota
	StoppedByCondition WorkerStopStatus = iota
	StoppeByError      WorkerStopStatus = iota
)
