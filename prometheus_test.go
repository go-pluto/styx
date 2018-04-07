package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSteps(t *testing.T) {
	assert.InEpsilon(t, 1*time.Second, steps(time.Minute), 0.01)
	assert.InEpsilon(t, 1*time.Second, steps(5*time.Minute), 0.01)
	assert.InEpsilon(t, 3*time.Second, steps(15*time.Minute), 0.01)
	assert.InEpsilon(t, 7*time.Second, steps(30*time.Minute), 0.02)
	assert.InEpsilon(t, 14*time.Second, steps(time.Hour), 0.02)
	assert.InEpsilon(t, 28*time.Second, steps(2*time.Hour), 0.02)
	assert.InEpsilon(t, 85*time.Second, steps(6*time.Hour), 0.01)
	assert.InEpsilon(t, 171*time.Second, steps(12*time.Hour), 0.01)
	assert.InEpsilon(t, 342*time.Second, steps(24*time.Hour), 0.01)
	assert.InEpsilon(t, 685*time.Second, steps(48*time.Hour), 0.01)
	assert.InEpsilon(t, 2400*time.Second, steps(168*time.Hour), 0.01)
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
