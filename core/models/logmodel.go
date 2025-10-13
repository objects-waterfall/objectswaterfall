package models

import (
	"time"

	"objectswaterfall.com/core/models/enums"
)

type WorkerJobLogModel struct {
	WorkerLog
	WorkerRequestLog
	StartTime            time.Time
	StopTime             time.Time
	RequestNumber        int64
	SuccessAttemptsCount int64
	FailedAttemptsCount  int64
}

type LogFunc func(WorkerJobLogModel)

type WorkerLog struct {
	WorkerName       string
	WorkerStopStatus enums.WorkerStopStatus
	TotalItemsToSend int64
	ItemsSended      int64
}

type WorkerRequestLog struct {
	CurrentRequestStatus     enums.WorkerRequestStatus
	RequestErrorMessage      string
	RequestDirationTime      float64
	MedianReuestDurationTime float64
}
