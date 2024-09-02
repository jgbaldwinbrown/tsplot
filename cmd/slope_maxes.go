package main

import (
	"fmt"
	"os"
	"bufio"
	"github.com/jgbaldwinbrown/tsplot/pkg/tsplot"
)

func main() {
	var paths []string
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		paths = append(paths, s.Text())
	}

	maxes, err := tsplot.PathsToMaxes(paths)
	if err != nil {
		panic(err)
	}
	fmt.Println(maxes)
}
