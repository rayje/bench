package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_NewMetrics(t *testing.T) {
	results := []Result{
		Result{Duration: time.Millisecond * 1},
		Result{Duration: time.Microsecond * 536},
	}

	totalRtt := time.Duration(0)
	for _, result := range results {
		totalRtt += result.Duration
	}
	mean := float64(totalRtt) / 2

	metrics := NewMetrics(results)
	assert.Equal(t, 2, metrics.Total, fmt.Sprintf("Expected 2, got %f", metrics.Total))
	assert.Equal(t, mean, metrics.Mean, fmt.Sprintf("Expected %f, got %f", mean, metrics.Mean))
}

func Test_GetSingleResultPercentile(t *testing.T) {
	results := []Result{
		Result{Duration: time.Microsecond * 100},
	}
	metrics := NewMetrics(results)

	expected := time.Microsecond * 100
	_25m := metrics.GetPercentile(0.25)
	_75m := metrics.GetPercentile(0.75)
	_95m := metrics.GetPercentile(0.95)
	_99m := metrics.GetPercentile(0.99)
	_999m := metrics.GetPercentile(0.999)
	assert.Equal(t, expected, _25m, fmt.Sprintf("Got an error: %+v", _25m))
	assert.Equal(t, expected, _75m, fmt.Sprintf("Got an error: %+v", _75m))
	assert.Equal(t, expected, _95m, fmt.Sprintf("Got an error: %+v", _95m))
	assert.Equal(t, expected, _99m, fmt.Sprintf("Got an error: %+v", _99m))
	assert.Equal(t, expected, _999m, fmt.Sprintf("Got an error: %+v", _999m))
}

func Test_GetPercentile(t *testing.T) {
	results := []Result{
		Result{Duration: time.Microsecond * 250},
		Result{Duration: time.Microsecond * 500},
		Result{Duration: time.Microsecond * 750},
		Result{Duration: time.Microsecond * 1000},
	}
	metrics := NewMetrics(results)

	_25m := metrics.GetPercentile(0.25)
	_75m := metrics.GetPercentile(0.75)
	_95m := metrics.GetPercentile(0.95)
	_99m := metrics.GetPercentile(0.99)
	_999m := metrics.GetPercentile(0.999)

	assert.Equal(t, time.Nanosecond*312500, _25m, fmt.Sprintf("Got an error _25m: %+v", _25m))
	assert.Equal(t, time.Nanosecond*937500, _75m, fmt.Sprintf("Got an error _75m: %+v", _75m))
	assert.Equal(t, time.Nanosecond*1000000, _95m, fmt.Sprintf("Got an error _95m: %+v", _95m))
	assert.Equal(t, time.Nanosecond*1000000, _99m, fmt.Sprintf("Got an error _99m: %+v", _99m))
	assert.Equal(t, time.Nanosecond*1000000, _999m, fmt.Sprintf("Got an error _999m: %+v", _999m))
}
