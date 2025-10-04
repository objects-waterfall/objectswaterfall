package models

type LogModel struct {
	Log                      string
	RequestDirationTime      float64
	MedianReuestDurationTime float64
	SuccessAttemptsCount     int
	FailedAttemptsCount      int
}

type LogFunc func(LogModel)
