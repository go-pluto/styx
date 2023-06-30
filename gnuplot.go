package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func gnuplotAction(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf(color.RedString("need a query to run"))
	}

	start := c.Timestamp("start")
	end := start.Add(c.Duration("duration"))
	prometheus := c.String("prometheus")

	// if resolution is not set, use 1s as default
	resolution := c.Duration("resolution")
	if resolution == 0 || resolution == time.Duration(0) {
		resolution = time.Second
	}
	if resolution < time.Second {
		return fmt.Errorf(color.RedString("resolution must be >= 1s"))
	}

	results, err := Query(prometheus, *start, end, resolution, c.Args().First())
	if err != nil {
		return err
	}

	header := "set grid\n" +
		"set key left top\n" +
		"set xdata time\n" +
		"set timefmt '%s'\n" +
		"set datafile separator ','\n"

	buf := bytes.NewBufferString(header)

	buf.WriteString("$DATA << EOD\n")
	if err := csvWriter(buf, results); err != nil {
		return err
	}
	buf.WriteString("EOD\n")

	for i, result := range results {
		if i == 0 {
			buf.WriteString("plot ")
		} else {
			buf.WriteString(", ")
		}
		plot := fmt.Sprintf("$DATA using 1:%d with lines lw 1 title '%s'", i+2, escapeMetricName(result.Metric))
		buf.WriteString(plot)
	}
	buf.WriteString("\n")

	fmt.Printf("%s\n", buf.String())

	return nil
}

func escapeMetricName(name string) string {
	// Escape: { } = _
	name = strings.Replace(name, `{`, `\{`, -1)
	name = strings.Replace(name, `}`, `\}`, -1)
	name = strings.Replace(name, `=`, `\=`, -1)
	name = strings.Replace(name, `_`, `\_`, -1)
	return name
}
