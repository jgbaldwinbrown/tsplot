package tsplot

import (
	"bufio"
	"github.com/jgbaldwinbrown/accel/accel"
	"io"
	"github.com/jgbaldwinbrown/lscan/lscan"
	"strconv"
)

func PlottableLineToAf(line []string) (af accel.Af, err error) {
	af.Pop = line[8]
	gen, err := strconv.ParseInt(line[7], 0, 64)
	if err != nil {
		return
	}
	af.Gen = int(gen)

	freq, err := strconv.ParseFloat(line[6], 64)
	if err != nil {
		return
	}
	af.Af = freq
	return af, nil
}

func ReadPlottable(r io.Reader) (p [][]string) {
	s := bufio.NewScanner(r)
	s.Buffer(make([]byte, 0), 1e12)
	scanf := lscan.ByByte('\t')

	for s.Scan() {
		line := []string{}
		line = lscan.SplitByFunc(line, s.Text(), scanf)
		p = append(p, line)
	}
	return p
}

func PlottableToAfs(ps [][]string) ([]accel.Af, error) {
	var afs []accel.Af
	for _, p := range ps {
		af, err := PlottableLineToAf(p)
		if err != nil {
			return nil, err
		}
		afs = append(afs, af)
	}
	return afs, nil
}
