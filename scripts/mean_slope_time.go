package main

import (
	"os"
	"bufio"
	"flag"
	"github.com/jgbaldwinbrown/tsplot/tsplot"
)

func main() {
	p, err := ReadPlottable(os.Stdin)
	if err != nil {
		panic(err)
	}

	afs := PlottableToAfs(p)
	s := accel.SmoothedMeanSlope(afs)
	fmt.Print(s)
}
