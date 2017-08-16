package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSteps(t *testing.T) {
	assert.Equal(t, 1, steps(time.Minute))
	assert.Equal(t, 1, steps(5*time.Minute))
	assert.Equal(t, 3, steps(15*time.Minute))
	assert.Equal(t, 7, steps(30*time.Minute))
	assert.Equal(t, 14, steps(time.Hour))
	assert.Equal(t, 28, steps(2*time.Hour))
	assert.Equal(t, 85, steps(6*time.Hour))
	assert.Equal(t, 171, steps(12*time.Hour))
	assert.Equal(t, 342, steps(24*time.Hour))
	assert.Equal(t, 685, steps(48*time.Hour))
	assert.Equal(t, 2400, steps(168*time.Hour))
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
