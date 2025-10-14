package dtos

import (
	"time"

	"objectswaterfall.com/core/models"
)

type WorkerResultsDto struct {
	WorkerName               string
	StartTime                time.Time
	StopTime                 time.Time
	MedianReuestDurationTime float64
	ItemsSended              int64
	SuccessAttemptsCount     int64
	FailedAttemptsCount      int64
}

func ToLogResult(log models.WorkerJobLogModel) WorkerResultsDto {
	return WorkerResultsDto{
		WorkerName:               log.WorkerName,
		StartTime:                log.StartTime,
		StopTime:                 log.StopTime,
		MedianReuestDurationTime: log.MedianReuestDurationTime,
		ItemsSended:              log.ItemsSended,
		SuccessAttemptsCount:     log.SuccessAttemptsCount,
		FailedAttemptsCount:      log.FailedAttemptsCount,
	}
}
