package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "styx"
	app.Usage = "Export metrics from prometheus"

	app.Action = exportAction
	app.Flags = []cli.Flag{
		&cli.DurationFlag{
			Name:  "duration,d",
			Usage: "The duration to get timeseries from",
			Value: time.Hour,
		},
		&cli.TimestampFlag{
			Name:   "start,s",
			Usage:  "The start time to get timeseries from",
			Layout: "2006-01-02T15:04:05",
		},
		&cli.DurationFlag{
			Name:  "resolution,r",
			Usage: "The requested resolution of the timeseries in seconds (default 1s)",
			Value: time.Second,
		},
		&cli.BoolFlag{
			Name:  "header",
			Usage: "Include a header into the csv file",
		},
		&cli.StringFlag{
			Name:  "prometheus",
			Value: "http://localhost:9090",
		},
	}

	app.Commands = []*cli.Command{{
		Name:   "gnuplot",
		Usage:  "Directly plot a graph with gnuplot",
		Action: gnuplotAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "prometheus",
				Value: "http://localhost:9090",
			},
			&cli.DurationFlag{
				Name:  "duration,d",
				Usage: "The duration to get timeseries from",
				Value: time.Hour,
			},
			&cli.TimestampFlag{
				Name:   "start,s",
				Usage:  "The start time to get timeseries from",
				Layout: "2006-01-02T15:04:05",
			},
			&cli.DurationFlag{
				Name:  "resolution,r",
				Usage: "The requested resolution of the timeseries in seconds (default 1s)",
				Value: time.Second,
			},
			&cli.StringFlag{
				Name:  "title",
				Usage: "Give the gnuplot graph a title",
			},
		},
	}, {
		Name:   "matplotlib",
		Usage:  "Generate a file that uses matplotlib",
		Action: matplotlibAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "prometheus",
				Value: "http://localhost:9090",
			},
			&cli.DurationFlag{
				Name:  "duration,d",
				Usage: "The duration to get timeseries from",
				Value: time.Hour,
			},
			&cli.TimestampFlag{
				Name:   "start,s",
				Usage:  "The start time to get timeseries from",
				Layout: "2006-01-02T15:04:05",
			},
			&cli.DurationFlag{
				Name:  "resolution,r",
				Usage: "The requested resolution of the timeseries in seconds (default 1s)",
				Value: time.Second,
			},
			&cli.StringFlag{
				Name:  "title",
				Usage: "Give the gnuplot graph a title",
			},
		},
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func exportAction(c *cli.Context) error {
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
	// fmt.Printf("Querying %s from %s to %s\n", c.Args().First(), *start, end)
	results, err := Query(prometheus, *start, end, resolution, c.Args().First())
	if err != nil {
		return err
	}

	// Only add a line as header when the flag is true
	header := c.Bool("header")
	if header {
		if err := csvHeaderWriter(os.Stdout, results); err != nil {
			return err
		}
	}

	return csvWriter(os.Stdout, results)
}
