package model

import (
	"time"
)

type TimeEntry struct {
	// DATE is the date of the time entry
	DATE time.Time `json:"date"`

	// List of time entries
	ENTRY []time.Time `json:"entry"`
}

func (t *TimeEntry) AddTime(time time.Time) *TimeEntry {
	t.ENTRY = append(t.ENTRY, time)
	return t
}

func (t *TimeEntry) Size() int {
	return len(t.ENTRY)
}

func New() *TimeEntry {
	return &TimeEntry{
		DATE:  time.Now(),
		ENTRY: []time.Time{},
	}
}
