package db

import "time"

type Vacation struct {
	ID        int64
	StartDate time.Time
	EndDate   time.Time
}
