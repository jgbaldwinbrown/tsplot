package tsplot

import (
	"os/exec"
	"io"
	"fmt"
	"os"
	"bufio"
	"math"
	"sort"
)

func GetFixations(freqs Iter[float64], gens Iter[float64]) FixationTime {
	fg := [][]float64{}
	for freq, ok := freqs.Next(); ok ; freq, ok = freqs.Next() {
		if gen, ok2 := gens.Next(); ok2 {
			fg = append(fg, []float64{freq, gen})
		} else {
			break
		}
	}

	sort.Slice(fg, func(i, j int) bool {return fg[i][1] > fg[j][1]})

	FixedFrom := -1
	LostFrom := -1
	var freqgen []float64
	hithresh := 0.9999999
	lothresh := 0.0000001
	for _, freqgen = range fg {
		igen := int(math.Round(freqgen[1]))
		if freqgen[0] > lothresh && freqgen[0] < hithresh {
			break
		}
		if freqgen[0] <= lothresh {
			if FixedFrom != -1 {
				break
			}
			LostFrom = igen
		}
		if freqgen[0] >= hithresh {
			if LostFrom != -1 {
				break
			}
			FixedFrom = igen
		}
	}
	if FixedFrom != -1 {
		return FixationTime{float64(FixedFrom), Fixed}
	}
	if LostFrom != -1 {
		return FixationTime{float64(LostFrom), Lost}
	}
	return FixationTime{-1.0, Unfixed}
}

type FixationState int

const (
	Unfixed FixationState = iota
	Fixed
	Lost
)

type FixationTime struct {
	Gen float64
	State FixationState
}

func ToFixationTime(sync SyncE, info []InfoE) FixationTime {
	freqs := CalcAfFreqs(sync.Afs)
	gens := GetGens(info)
	return GetFixations(NewTableColIter(freqs, sync.PlotcolPrimary), NewSliceIter(gens))
}

func ToFixationTimes(sync []SyncE, info []InfoE) []FixationTime {
	var times []FixationTime
	for _, synce := range sync {
		times = append(times, ToFixationTime(synce, info))
	}
	return times
}

type FixationCount struct {
	Name string
	Gen float64
	Total int
	Fixed int
	Lost int
	FixedOrLost int
	FixedFreq float64
	LostFreq float64
	FixedOrLostFreq float64
}

func (c *FixationCount) AddFixation(t FixationTime) {
	c.Total++
	if t.Gen == -1.0 || t.Gen > c.Gen {
		return
	}

	if t.State == Fixed || t.State == Lost {
		c.FixedOrLost++
	}

	if t.State == Fixed {
		c.Fixed++
	}

	if t.State == Lost {
		c.Lost++
	}
}

func (c *FixationCount) CalcFreqs() {
	if c.Total == 0 {
		c.FixedFreq = math.NaN()
		c.LostFreq = math.NaN()
		c.FixedOrLostFreq = math.NaN()
		return
	}
	c.FixedFreq = float64(c.Fixed) / float64(c.Total)
	c.LostFreq = float64(c.Lost) / float64(c.Total)
	c.FixedOrLostFreq = float64(c.FixedOrLost) / float64(c.Total)
}

func CountFixations(times []FixationTime) []FixationCount {
	var counts []FixationCount
	for i:=0.0; i<54.1; i+=6.0 {
		var count FixationCount
		count.Gen = i
		for _, time := range times {
			count.AddFixation(time)
		}
		count.CalcFreqs()
		counts = append(counts, count)
	}
	return counts
}

func NameFixationCounts(counts []FixationCount, name string) {
	for i, _ := range counts {
		counts[i].Name = name
	}
}

func ToFixationCounts(name string, sync []SyncE, info []InfoE) []FixationCount {
	times := ToFixationTimes(sync, info)
	counts := CountFixations(times)
	NameFixationCounts(counts, name)
	return counts
}

func FprintFixationCount(w io.Writer, count FixationCount) {
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\n",
		count.Name,
		count.Gen,
		count.Total,
		count.Fixed,
		count.Lost,
		count.FixedOrLost,
		count.FixedFreq,
		count.LostFreq,
		count.FixedOrLostFreq,
	)
}

func FprintFixationCounts(w io.Writer, counts ...FixationCount) {
	fmt.Fprintf(w, "Name\tTime\tTotal\tFixed\tLost\tFixedOrLost\tFixedFreq\tLostFreq\tFixedOrLostFreq\n")
	for _, count := range counts {
		FprintFixationCount(w, count)
	}
}

func WriteFixationCounts(counts []FixationCount, path string) error {
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	FprintFixationCounts(bw, counts...)
	return nil
}

func PlotFixCurve(inpath, outpath string) error {
	cmd := exec.Command("plotfixationcurve", inpath, outpath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

