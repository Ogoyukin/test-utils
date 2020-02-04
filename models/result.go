package models

import (
	"time"
)

type Result struct {
	Completed int64
	Failed    int64
	Duration  time.Duration
}

func (m *Result) RequestDuration(threadsCount int) time.Duration {
	duration := m.Duration.Nanoseconds() / (m.Completed / int64(threadsCount))
	return time.Duration(duration)
}
