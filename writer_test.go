package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCSVWriter(t *testing.T) {
	// No results
	buf := bytes.NewBuffer(nil)
	assert.NoError(t, csvWriter(buf, nil))
	assert.Equal(t, "", buf.String())

	// Result with one entry
	res := []Result{{
		Metric: "foobar",
		Values: map[string]string{
			"1502749393": "42",
		},
	}}
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, csvWriter(buf, res))
	assert.Equal(t, "1502749393,42\n", buf.String())

	// One result with multiple time series
	res = []Result{{
		Metric: "foobar",
		Values: map[string]string{
			"1502749391": "1",
			"1502749392": "2",
			"1502749393": "3",
			"1502749394": "4",
			"1502749395": "5",
		},
	}}
	expected := "1502749391,1\n1502749392,2\n1502749393,3\n1502749394,4\n1502749395,5\n"
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, csvWriter(buf, res))
	assert.Equal(t, expected, buf.String())

	// Two results with multiple time series
	res = []Result{{
		Metric: "foobar",
		Values: map[string]string{
			"1502749390": "0",
			"1502749391": "1",
			"1502749392": "2",
			"1502749393": "3",
			"1502749394": "4",
		},
	}, {
		Metric: "foobaz",
		Values: map[string]string{
			"1502749390": "5",
			"1502749391": "6",
			"1502749392": "7",
			"1502749393": "8",
			"1502749394": "9",
		},
	}}
	expected = "1502749390,0,5\n1502749391,1,6\n1502749392,2,7\n1502749393,3,8\n1502749394,4,9\n"
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, csvWriter(buf, res))
	assert.Equal(t, expected, buf.String())

	// Two results with multiple time series
	res = []Result{{
		Metric: "foobar",
		Values: map[string]string{
			"1502749390": "0",
			"1502749392": "2",
			"1502749393": "3",
			"1502749394": "4",
			"1502749396": "10",
		},
	}, {
		Metric: "foobaz",
		Values: map[string]string{
			"1502749390": "5",
			"1502749391": "6",
			"1502749392": "7",
			"1502749393": "8",
			"1502749394": "9",
		},
	}}
	expected = "1502749390,0,5\n1502749391,,6\n1502749392,2,7\n1502749393,3,8\n1502749394,4,9\n1502749396,10,\n"
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, csvWriter(buf, res))
	assert.Equal(t, expected, buf.String())
}

func TestCSVHeaderWriter(t *testing.T) {
	// No results
	buf := bytes.NewBuffer(nil)
	assert.NoError(t, csvHeaderWriter(buf, nil))
	assert.Equal(t, "", buf.String())

	// Result with one entry
	res := []Result{{
		Metric: "foobar",
		Values: map[string]string{
			"1502749393": "42",
		},
	}}
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, csvHeaderWriter(buf, res))
	assert.Equal(t, "Time,foobar\n", buf.String())

	// Two results with multiple time series
	res = []Result{{
		Metric: "foobar",
		Values: map[string]string{
			"1502749390": "0",
		},
	}, {
		Metric: "foobaz",
		Values: map[string]string{
			"1502749390": "5",
		},
	}}
	expected := "Time,foobar,foobaz\n"
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, csvHeaderWriter(buf, res))
	assert.Equal(t, expected, buf.String())

}

func TestMatplotlibWriter(t *testing.T) {
	// No results
	buf := bytes.NewBuffer(nil)
	assert.NoError(t, matplotlibWriter(buf, nil))
	assert.Equal(t, "", buf.String())

	// Result with one entry
	res := []Result{{
		Metric: "foobar",
		Values: map[string]string{
			"1502749393": "42",
		},
	}}
	expected := "t = [1502749393]\ns0 = [42]\nplot.plot(t, s0)\n"
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, matplotlibWriter(buf, res))
	assert.Equal(t, expected, buf.String())

	// One result with multiple time series
	res = []Result{{
		Metric: "foobar",
		Values: map[string]string{
			"1502749391": "1",
			"1502749392": "2",
			"1502749393": "3",
			"1502749394": "4",
			"1502749395": "5",
		},
	}}
	expected = "t = [1502749391, 1502749392, 1502749393, 1502749394, 1502749395]\n" +
		"s0 = [1, 2, 3, 4, 5]\n" +
		"plot.plot(t, s0)\n"
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, matplotlibWriter(buf, res))
	assert.Equal(t, expected, buf.String())

	// Two results with multiple time series
	res = []Result{{
		Metric: "foobar",
		Values: map[string]string{
			"1502749390": "0",
			"1502749391": "1",
			"1502749392": "2",
			"1502749393": "3",
			"1502749394": "4",
		},
	}, {
		Metric: "foobaz",
		Values: map[string]string{
			"1502749390": "5",
			"1502749391": "6",
			"1502749392": "7",
			"1502749393": "8",
			"1502749394": "9",
		},
	}}
	expected = "t = [1502749390, 1502749391, 1502749392, 1502749393, 1502749394]\n" +
		"s0 = [0, 1, 2, 3, 4]\n" +
		"plot.plot(t, s0)\n" +
		"s1 = [5, 6, 7, 8, 9]\n" +
		"plot.plot(t, s1)\n"
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, matplotlibWriter(buf, res))
	assert.Equal(t, expected, buf.String())

	// Two results with multiple time series
	res = []Result{{
		Metric: "foobar",
		Values: map[string]string{
			"1502749390": "0",
			"1502749392": "2",
			"1502749393": "3",
			"1502749394": "4",
			"1502749396": "10",
		},
	}, {
		Metric: "foobaz",
		Values: map[string]string{
			"1502749390": "5",
			"1502749391": "6",
			"1502749392": "7",
			"1502749393": "8",
			"1502749394": "9",
		},
	}}
	expected = "t = [1502749390, 1502749391, 1502749392, 1502749393, 1502749394, 1502749396]\n" +
		"s0 = [0, None, 2, 3, 4, 10]\n" +
		"plot.plot(t, s0)\n" +
		"s1 = [5, 6, 7, 8, 9, None]\n" +
		"plot.plot(t, s1)\n"
	buf = bytes.NewBuffer(nil)
	assert.NoError(t, matplotlibWriter(buf, res))
	assert.Equal(t, expected, buf.String())
}
