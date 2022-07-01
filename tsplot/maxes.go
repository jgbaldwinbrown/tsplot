package main

import (
	"fmt"
	"os"
	"bufio"
	"github.com/jgbaldwinbrown/tsplot/tsplot"
	"github.com/jgbaldwinbrown/accel/accel"
)

func PathsToMaxes(paths []string) ([]float64, error) {
	var afsets []accel.Afs
	for _, path := range paths
		conn, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		p := tsplot.ReadPlottable(conn)

		afs, err := tsplot.PlottableToAfs(p)
		if err != nil {
			return nil, err
		}
		afsets = append(afsets, afs)
	}

	maxes := accel.MaxSlopeTimes(afsets)
	return maxes, nil
}
