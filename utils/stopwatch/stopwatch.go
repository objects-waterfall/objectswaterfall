package stopwatch

import "time"

type StopWatch struct {
	start time.Time
}

func NewStopWatch() *StopWatch {
	return &StopWatch{}
}

func (s *StopWatch) Start() {
	s.start = time.Now()
}

func (s *StopWatch) Elapsed(unit time.Duration) float64 {
	return float64(time.Since(s.start)) / float64(unit)
}
