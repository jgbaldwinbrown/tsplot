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

type AfStats struct {
	Freqs [][]float64
	MajorAllele int
	MinorAllele int
	Gens []float64
}

func CalcAfStats(afs [][]int64, info []InfoE) AfStats {
	out := AfStats{}
	out.Freqs = CalcAfFreqs(afs)
	out.MajorAllele, out.MinorAllele = GetMajorTwoAFs(afs)
	out.Gens = GetGens(info)
	return out
}

func GetHigher36(expStats, controlStats AfStats) (beneficial, deleterious int) {
	var exp36Major, control36Major float64
	for i:=0; i<len(expStats.Freqs); i++ {
		if expStats.Gens[i] >= 35.5 && expStats.Gens[i] <= 36.5 {
			exp36Major = expStats.Freqs[i][expStats.MajorAllele]
		}
		if controlStats.Gens[i] >= 35.5 && controlStats.Gens[i] <= 36.5 {
			control36Major = controlStats.Freqs[i][controlStats.MajorAllele]
		}
	}
	if exp36Major >= control36Major {
		return expStats.MajorAllele, expStats.MinorAllele
	}
	return expStats.MinorAllele, expStats.MajorAllele
}

func GetBeneficialAndDeleteriousAfsByG36(afs [][]int64, info []InfoE, controlafs [][]int64, controlinfo []InfoE) (beneficial, deleterious int) {
	expStats := CalcAfStats(afs, info)
	controlStats := CalcAfStats(controlafs, controlinfo)
	controlStats.MajorAllele = expStats.MajorAllele
	controlStats.MinorAllele = expStats.MinorAllele

	beneficial, deleterious = GetHigher36(expStats, controlStats)
	return beneficial, deleterious
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

func SetBeneficialHigh36(s *SyncE, info []InfoE, expS SyncE, expInfo []InfoE, controlS SyncE, controlInfo []InfoE) {
	beneficial, deleterious := GetBeneficialAndDeleteriousAfsByG36(expS.Afs, expInfo, controlS.Afs, controlInfo)
	s.PlotcolPrimary = beneficial
	s.PlotcolSecondary = deleterious
}

func SetBeneficialsHigh36(s []SyncE, info []InfoE, expS []SyncE, expInfo []InfoE, controlS []SyncE, controlInfo []InfoE) {
	for i, _ := range s {
		SetBeneficialHigh36(&s[i], info, expS[i], expInfo, controlS[i], controlInfo)
	}
}

func SetBeneficialOld(s *SyncE, info []InfoE) {
	beneficial, deleterious := GetBeneficialAndDeleteriousAFs(s.Afs, info)
	s.PlotcolPrimary = beneficial
	s.PlotcolSecondary = deleterious
}

func SetBeneficialsOld(s []SyncE, info []InfoE) {
	for i, _ := range s {
		SetBeneficialOld(&s[i], info)
	}
}

func SetBeneficial(s *SyncE, info []InfoE, beneS SyncE, beneInfo []InfoE) {
	beneficial, deleterious := GetBeneficialAndDeleteriousAFs(beneS.Afs, beneInfo)
	s.PlotcolPrimary = beneficial
	s.PlotcolSecondary = deleterious
}

func SetBeneficials(s []SyncE, info []InfoE, beneS []SyncE, beneInfo []InfoE) {
	for i, _ := range s {
		SetBeneficial(&s[i], info, beneS[i], beneInfo)
	}
}
