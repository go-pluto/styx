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
			Name:        "duration,d",
			Usage:       "The duration to get timeseries from",
			Value:       time.Hour,
			Destination: &flag.Duration,
		},
		&cli.TimestampFlag{
			Name:        "start,s",
			Usage:       "The start time to get timeseries from",
			Layout: 	 "2006-01-02T15:04:05",
			Destination: &flag.Start,
		},
		&cli.BoolFlag{
			Name:        "header",
			Usage:       "Include a header into the csv file",
			Destination: &flag.Header,
		},
		&cli.StringFlag{
			Name:        "prometheus",
			Value:       "http://localhost:9090",
			Destination: &flag.Prometheus,
		},
	}

	app.Commands = []*cli.Command{{
		Name:   "gnuplot",
		Usage:  "Directly plot a graph with gnuplot",
		Action: gnuplotAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "prometheus",
				Value:       "http://localhost:9090",
				Destination: &gnuplotFlag.Prometheus,
			},
			&cli.DurationFlag{
				Name:        "duration,d",
				Usage:       "The duration to get timeseries from",
				Value:       time.Hour,
				Destination: &gnuplotFlag.Duration,
			},
			&cli.TimestampFlag{
				Name:        "start,s",
				Usage:       "The start time to get timeseries from",
				Layout: 	 "2006-01-02T15:04:05",
				Destination: &gnuplotFlag.Start,
			},
			&cli.StringFlag{
				Name:        "title",
				Usage:       "Give the gnuplot graph a title",
				Destination: &gnuplotFlag.Title,
			},
		},
	}, {
		Name:   "matplotlib",
		Usage:  "Generate a file that uses matplotlib",
		Action: matplotlibAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "prometheus",
				Value:       "http://localhost:9090",
				Destination: &matplotlibFlag.Prometheus,
			},
			&cli.DurationFlag{
				Name:        "duration,d",
				Usage:       "The duration to get timeseries from",
				Value:       time.Hour,
				Destination: &matplotlibFlag.Duration,
			},
			&cli.TimestampFlag{
				Name:        "start,s",
				Usage:       "The start time to get timeseries from",
				Layout: 	 "2006-01-02T15:04:05",
				Destination: &matplotlibFlag.Start,
			},
			&cli.StringFlag{
				Name:        "title",
				Usage:       "Give the gnuplot graph a title",
				Destination: &matplotlibFlag.Title,
			},
		},
	}}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

type flags struct {
	Duration   time.Duration
	Start      cli.Timestamp
	Header     bool
	Prometheus string
}

var flag flags

func exportAction(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf(color.RedString("need a query to run"))
	}
	
	start, err := time.Parse("2006-01-02T15:04:05", flag.Start.String())
	if err != nil {
		return err
	}
	end := start.Add(flag.Duration)

	results, err := Query(flag.Prometheus, start, end, c.Args().First())
	if err != nil {
		return err
	}

	// Only add a line as header when the flag is true, which is the default
	if flag.Header {
		if err := csvHeaderWriter(os.Stdout, results); err != nil {
			return err
		}
	}

	return csvWriter(os.Stdout, results)
}
