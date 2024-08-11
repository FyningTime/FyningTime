package db

import "time"

type Worktime struct {
	ID        int64
	Type      string
	Time      time.Time
	Breaktime time.Time
	Workday   Workday
}
