package main

import (
	"sort"
	"time"
)

type Result struct {
	Code      uint16
	Timestamp time.Time
	Duration  time.Duration
	Error     string
}

type Results []Result

func (r Results) Sort() Results {
	sort.Sort(r)
	return r
}

func (r Results) Len() int           { return len(r) }
func (r Results) Less(i, j int) bool { return r[i].Duration < r[j].Duration }
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
