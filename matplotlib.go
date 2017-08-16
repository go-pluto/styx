package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

type matplotlibFlags struct {
	Duration   time.Duration
	Prometheus string
	Title      string
}

var matplotlibFlag matplotlibFlags

func matplotlibAction(c *cli.Context) error {
	if !c.Args().Present() {
		return fmt.Errorf(color.RedString("need a query to run"))
	}

	end := time.Now()
	start := end.Add(-1 * matplotlibFlag.Duration)

	results, err := Query(matplotlibFlag.Prometheus, start, end, c.Args().First())
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
