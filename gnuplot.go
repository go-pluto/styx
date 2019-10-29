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

	results, err := Query(gnuplotFlag.Prometheus, start, end, c.Args().First())
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
