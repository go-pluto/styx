package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSteps(t *testing.T) {
	assert.Equal(t, 1, steps(time.Second, time.Second))
	assert.Equal(t, 1, steps(time.Second, time.Millisecond))
	assert.Equal(t, 1, steps(time.Second, time.Microsecond))
	assert.Equal(t, 1, steps(time.Minute, time.Second))
	assert.Equal(t, 1, steps(time.Minute, time.Millisecond))
	assert.Equal(t, 1, steps(time.Minute, time.Microsecond))
	assert.Equal(t, 1, steps(time.Hour, time.Second))
	assert.Equal(t, 2, steps(time.Duration(12800)*time.Second*3, 2*time.Second))
	assert.Equal(t, 1, steps(time.Duration(12800)*time.Second*3, 2*time.Millisecond))
	assert.Equal(t, 1, steps(time.Duration(12800)*time.Second*3, 2*time.Microsecond))
	assert.Equal(t, 1, steps(time.Duration(12800)*time.Second*3, 2*time.Nanosecond))
	assert.Equal(t, 1, steps(time.Duration(12800)*time.Second*3, 100*time.Nanosecond))

	// test on the fly 2023-06-28T13:19:16' + 15m with various resolutions
	assert.Equal(t, 60, steps(
		time.Date(2023, 6, 28, 13, 34, 16, 0, time.UTC).Sub(
			time.Date(2023, 6, 28, 13, 19, 16, 0, time.UTC)), 60*time.Second))

}

func TestMetricName(t *testing.T) {
	metric := make(map[string]string)
	assert.Equal(t, `{}`, metricName(metric))

	metric["__name__"] = "go_goroutines"
	assert.Equal(t, `go_goroutines`, metricName(metric))

	metric["job"] = "prometheus"
	assert.Equal(t, `go_goroutines{job="prometheus"}`, metricName(metric))

	metric["instance"] = "localhost:9090"
	assert.Equal(t, `go_goroutines{instance="localhost:9090",job="prometheus"}`, metricName(metric))
}
