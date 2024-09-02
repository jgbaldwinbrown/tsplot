package tsplot

import (
	"github.com/sajari/regression"
)

func TrainTable(table [][]float64, depcol int, depname string, indepcols []int, indepnames []string) *regression.Regression {
	totrain := make(regression.DataPoints, len(table))
	for i, row := range table {
		indeps := make([]float64, len(indepcols))
		for j, col := range indepcols {
			indeps[j] = row[col]
		}
		totrain[i] = regression.DataPoint(row[depcol], indeps)
	}

	r := new(regression.Regression)
	r.SetObserved(depname)
	for i, name := range indepnames {
		r.SetVar(i, name)
	}
	r.Train(totrain...)
	r.Run()
	return r
}

type Iter[T any] interface {
	Next() (T, bool)
}

type SliceIter[T any] struct {
	Slice []T
	Idx int
}

func (s *SliceIter[T]) Next() (T, bool) {
	s.Idx++
	if s.Idx >= len(s.Slice) {
		var t T
		return t, false
	}
	return s.Slice[s.Idx], true
}

func NewSliceIter[T any](slice []T) *SliceIter[T] {
	s := new(SliceIter[T])
	s.Idx = -1
	s.Slice = slice
	return s
}

type TableColIter[T any] struct {
	Table [][]T
	Col int
	Idx int
}

func (s *TableColIter[T]) Next() (T, bool) {
	s.Idx++
	if s.Idx >= len(s.Table) {
		var t T
		return t, false
	}
	return s.Table[s.Idx][s.Col], true
}

func NewTableColIter[T any](table [][]T, col int) *TableColIter[T] {
	s := new(TableColIter[T])
	s.Idx = -1
	s.Table = table
	s.Col = col
	return s
}

func GetSlope(xs Iter[float64], ys Iter[float64], maxgen float64) float64 {
	var totrain regression.DataPoints
	for {
		x, xok := xs.Next()
		y, yok := ys.Next()
		if !xok || !yok {
			break
		}

		if x < maxgen {
			totrain = append(totrain, regression.DataPoint(x, []float64{y}))
		}
	}

	r := new(regression.Regression)
	r.Train(totrain...)
	r.Run()
	return r.Coeff(1)
}
