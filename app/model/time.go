package model

import (
	"time"
)

type TimeEntry struct {
	// DATE is the date of the time entry
	DATE time.Time `json:"date"`

	// List of time entries
	ENTRIES []time.Time `json:"entries"`
}

func (t *TimeEntry) AddTime(time time.Time) *TimeEntry {
	t.ENTRIES = append(t.ENTRIES, time)
	return t
}

func (t *TimeEntry) Size() int {
	return len(t.ENTRIES)
}

func New() *TimeEntry {
	return &TimeEntry{
		DATE:    time.Now(),
		ENTRIES: []time.Time{},
	}
}
