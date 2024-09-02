package tsplot

import (
	"testing"
	"fmt"
)

func TestGetFixations(t *testing.T) {
	inf := []float64{.5, .5, .6, 1.0, 1.0, 1.0}
	ingen := []float64{0, 6, 12, 18, 24, 36}
	fmt.Println("fixation:", GetFixations(NewSliceIter(inf), NewSliceIter(ingen)))
}
	// func GetFixations(freqs Iter[float64], gens Iter[float64]) FixationTime {
