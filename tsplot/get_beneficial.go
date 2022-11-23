package tsplot

import (
)

func GetGens(info []InfoE) []float64 {
	out := make([]float64, len(info))
	for i, infoe := range info {
		out[i] = float64(infoe.Gen)
	}
	return out
}

func GetBeneficialAndDeleteriousAFs(afs [][]int64, info []InfoE) (beneficial, deleterious int) {
	const maxgen float64 = 49.0

	freqs := CalcAfFreqs(afs)

	major, minor := GetMajorTwoAFs(afs)
	gens := GetGens(info)

	majorslope := GetSlope(NewTableColIter(freqs, major), NewSliceIter(gens), maxgen)
	// minorslope := GetSlope(NewTableColIter(afs, minor), NewSliceIter(gens))

	if majorslope >= 0 {
		return major, minor
	}
	return minor, major
}

func SetBeneficial(s *SyncE, info []InfoE) {
	beneficial, deleterious := GetBeneficialAndDeleteriousAFs(s.Afs, info)
	s.PlotcolPrimary = beneficial
	s.PlotcolSecondary = deleterious
}

func SetBeneficials(s []SyncE, info []InfoE) {
	for i, _ := range s {
		SetBeneficial(&s[i], info)
	}
}
