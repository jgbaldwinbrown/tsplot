package main

import (
	"github.com/jgbaldwinbrown/tsplot/tsplot"
	"flag"
	"bufio"
	"os"
)

func main() {
	size := flag.Int("s", 1, "Size of new windows")
	flag.Parse()

	s := bufio.NewScanner(os.Stdin)
	s.Buffer(make([]byte, 0), 1e12)
	bed, err := tsplot.ReadBedScanner(s)
	if err != nil {
		panic(err)
	}

	resized := tsplot.ResizeBed(bed, int64(*size))

	tsplot.WriteBed(os.Stdout, resized)
}
