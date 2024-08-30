package db

import "time"

type Workday struct {
	ID        int64
	Date      time.Time
	Time      string
	Breaktime string
}
