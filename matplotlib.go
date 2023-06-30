package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func matplotlibAction(c *cli.Context) error {
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

	header := "import matplotlib.pyplot as plot\n\n"
	buf := bytes.NewBufferString(header)

	if err := matplotlibWriter(buf, results); err != nil {
		return err
	}

	if err := matplotlibLegendWriter(buf, results); err != nil {
		return err
	}

	footer := "plot.grid(True)\n" +
		fmt.Sprintf("plot.title('%s')\n", c.Args().First()) +
		"plot.show()\n"

	buf.WriteString(footer)

	_, err = fmt.Fprint(os.Stdout, buf.String())
	return err
}

//import matplotlib.pyplot as plt
//
//t = [1502573433, ...]
//s = [231, ...]
//plt.plot(t, s)
//plt.plot(t, u)
//
//plt.legend(['y = x', 'y = 2x', 'y = 3x', 'y = 4x'], loc='upper left')
//
//plt.xlabel('time (s)')
//plt.ylabel('voltage (mV)')
//plt.title('go_goroutines{asdf="asdf"}')
//plt.grid(True)
//plt.show()
