package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

type gnuplotFlags struct {
	Duration   time.Duration
	Prometheus string
	Title      string
}

var gnuplotFlag gnuplotFlags

func gnuplotAction(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf(color.RedString("need a query to run"))
	}

	end := time.Now()
	start := end.Add(-1 * gnuplotFlag.Duration)

	results, err := Query(gnuplotFlag.Prometheus, start, end, c.Args().First(), 0)
	if err != nil {
		return err
	}

	header := "set grid\n" +
		"set key left top\n" +
		"set xdata time\n" +
		"set timefmt '%s'\n" +
		"set datafile separator ','\n"

	buf := bytes.NewBufferString(header)

	for i, result := range results {
		plot := fmt.Sprintf("plot '-' using 1:%d with lines lw 1 title '%s'\n", i+2, escapeMetricName(result.Metric))
		buf.WriteString(plot)
	}

	if err := csvWriter(buf, results); err != nil {
		return err
	}

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
