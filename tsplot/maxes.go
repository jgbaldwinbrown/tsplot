package tsplot

import (
	"os"
	"github.com/jgbaldwinbrown/accel/accel"
)

func PathsToMaxes(paths []string) ([]float64, error) {
	var afsets []accel.Afs
	for _, path := range paths {
		conn, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		p := ReadPlottable(conn)

		afs, err := PlottableToAfs(p)
		if err != nil {
			return nil, err
		}
		afsets = append(afsets, afs)
	}

	maxes, err := accel.MaxSlopeTimes(afsets)
	return maxes, err
}
