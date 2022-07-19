package main

import (
	"fmt"
	"os"
	"github.com/jgbaldwinbrown/tsplot/tsplot"
	"github.com/jgbaldwinbrown/accel/accel"
)

func main() {
	p := tsplot.ReadPlottable(os.Stdin)

	afs, err := tsplot.PlottableToAfs(p)
	if err != nil {
		panic(err)
	}

	xs, ys := accel.SmoothedMeanSlope(afs)
	fmt.Print(xs)
	fmt.Print(ys)
}
