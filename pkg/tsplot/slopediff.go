package tsplot

import (
	"github.com/jgbaldwinbrown/pmap/pkg"
	"flag"
	"io"
	"bufio"
	"fmt"
	"os"
)

type SlopeDiffOpts struct {
	Threads int
}

type SlopeDiffCfg struct {
	PolarizerPaths []string
	SlopePaths []string
	OutPre string
}

type SlopeDiffArgs struct {
	Cfg SlopeDiffCfg
	Opts SlopeDiffOpts
	Polarizations map[string][]int
}

// func UniquePolarizerPathSets(cfgs ...SlopeDiffCfg) map[string][]string {
// 	upaths := make(map[string][]string)
// 	for _, cfg := range cfgs {
// 		paths := slices.Clone(cfg.PolarizerPaths)
// 		
// 	}
// }

// func Polarize(path string) []int, error {
// }

// func PolarizeAll(cfgs []SlopeDiffCfg, threads int) {
// 	topolarize := UniquePolarizerPaths(cfgs...)
// 	// MapE(function, cfgs, threads)
// }
// 
// func RunSlopeDiff() {
// 	opts := GetSlopeDiffOpts()
// 	cfgs := ReadSlopeDiffCfgs(os.Stdin)
// 
// 	polarizations, err := PolarizeAll(cfgs, opts.Threads)
// 	if err != nil {
// 		panic(err)
// 	}
// 
// 	args := CollectSlopeDiffArgs(cfgs, opts, polarizations)
// 
// 	err := AllSlopeDiffs(args, opts.Threads)
// 	if err != nil {
// 		panic(err)
// 	}
// 
// 	err := PlotAllSlopeDiffs(args, opts.Threads)
// 	if err != nil {
// 		panic(err)
// 	}
// }

type SlopeIter struct {
	SlopeExpSync Iter[SyncE]
	SlopeExpInfo []InfoE
	SlopeExpCols []int
	SlopeControlSync Iter[SyncE]
	SlopeControlInfo []InfoE
	SlopeControlCols []int
	PolarizeExpSync Iter[SyncE]
	PolarizeExpInfo []InfoE
	PolarizeExpCols []int
	PolarizeControlSync Iter[SyncE]
	PolarizeControlInfo []InfoE
	PolarizeControlCols []int
}

func NewSlopeIter(
	slopeExpSync Iter[SyncE], slopeExpInfo []InfoE,
	slopeControlSync Iter[SyncE], slopeControlInfo []InfoE,
	polarizeExpSync Iter[SyncE], polarizeExpInfo []InfoE,
	polarizeControlSync Iter[SyncE], polarizeControlInfo []InfoE,
	toUse InfoSelection,
) *SlopeIter {
	// SetBeneficialHigh36(s *SyncE, info []InfoE, expS SyncE, expInfo []InfoE, controlS SyncE, controlInfo []InfoE)
	s := new(SlopeIter)
	secols := GoodCols(slopeExpInfo, toUse)
	sccols := GoodCols(slopeControlInfo, toUse)
	pecols := GoodCols(polarizeExpInfo, toUse)
	pccols := GoodCols(polarizeControlInfo, toUse)
	*s = SlopeIter{
		slopeExpSync, SubsetInfoByCols(slopeExpInfo, secols), secols,
		slopeControlSync, SubsetInfoByCols(slopeControlInfo, sccols), sccols,
		polarizeExpSync, SubsetInfoByCols(polarizeExpInfo, pecols), pecols,
		polarizeControlSync, SubsetInfoByCols(polarizeControlInfo, pccols), pccols,
	}
	return s
}

type Slopes struct {
	Chr string
	Pos int64
	Experimental float64
	Control float64
}

func (s *SlopeIter) Next() (Slopes, bool) {
	slopeExp, ok := s.SlopeExpSync.Next()
	if !ok { return Slopes{}, false }
	slopeExp = SubsetSyncEByCols(slopeExp, s.SlopeExpCols)

	slopeControl, ok := s.SlopeControlSync.Next()
	if !ok { return Slopes{}, false }
	slopeControl = SubsetSyncEByCols(slopeControl, s.SlopeControlCols)

	polarizeExp, ok := s.PolarizeExpSync.Next()
	if !ok { return Slopes{}, false }
	polarizeExp = SubsetSyncEByCols(polarizeExp, s.PolarizeExpCols)

	polarizeControl, ok := s.PolarizeControlSync.Next()
	if !ok { return Slopes{}, false }
	polarizeControl = SubsetSyncEByCols(polarizeControl, s.PolarizeControlCols)

	beneficial, _ := GetBeneficialAndDeleteriousAfsByG36(polarizeExp.Afs, s.PolarizeExpInfo, polarizeControl.Afs, s.PolarizeControlInfo)

	var slopes Slopes
	slopes.Chr = slopeExp.Chr
	slopes.Pos = slopeExp.Pos
	slopes.Experimental = GetSyncESlope(slopeExp, s.SlopeExpInfo, beneficial)
	slopes.Control = GetSyncESlope(slopeControl, s.SlopeControlInfo, beneficial)

	return slopes, true
}

func FprintSlopes(w io.Writer, s Slopes) error {
	fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\t%v\n", s.Chr, s.Pos, s.Pos+1, s.Experimental, s.Control, s.Experimental - s.Control)
	return nil
}

func FprintSlopesIter(w io.Writer, si Iter[Slopes]) error {
	for s, ok := si.Next(); ok; s, ok = si.Next() {
		err := FprintSlopes(w, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func WriteSlopes(si Iter[Slopes], outpath string) error {
	w, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer w.Close()

	bw := bufio.NewWriter(w)
	defer bw.Flush()

	err = FprintSlopesIter(bw, si)
	if err != nil {
		return err
	}

	return nil
}

func SlopeDiffFull(pcfg PlotCfg, outpath string) error {
	polarizeExpsync, polarizeExpinfo, err := ReadSyncIterInfo(pcfg.BeneficialExpSbi)
	if err != nil { return err }
	defer polarizeExpsync.Close()

	polarizeControlsync, polarizeControlinfo, err := ReadSyncIterInfo(pcfg.BeneficialControlSbi)
	if err != nil { return err }
	defer polarizeControlsync.Close()

	slopeExpsync, slopeExpInfo, err := ReadSyncIterInfo(pcfg.Sbi)
	if err != nil { return err }
	defer slopeExpsync.Close()

	slopeControlsync, slopeControlInfo, err := ReadSyncIterInfo(pcfg.SlopeControlSbi)
	if err != nil { return err }
	defer slopeControlsync.Close()

	slopes := NewSlopeIter(
		slopeExpsync, slopeExpInfo,
		slopeControlsync, slopeControlInfo,
		polarizeExpsync, polarizeExpinfo,
		polarizeControlsync, polarizeControlinfo,
		pcfg.ToUse,
	)

	err = WriteSlopes(slopes, outpath)
	if err != nil { return err }

	return nil
}

func ReadSyncIterInfo(sbi SyncBedInfo) (*SyncIter, []InfoE, error) {
	sync, err := OpenSyncIter(sbi.Sync)
	if err != nil {
		return nil, nil, err
	}

	info, err := ReadInfo(sbi.Info)
	if err != nil {
		return nil, nil, err
	}
	return sync, info, nil
}

func GetSyncESlope(se SyncE, info []InfoE, allele int) float64 {
	const maxgen float64 = 49.0

	freqs := CalcAfFreqs(se.Afs)
	gens := GetGens(info)

	slope := GetSlope(NewTableColIter(freqs, allele), NewSliceIter(gens), maxgen)
	return slope
}

func SlopeDiffFullMcfg(mcfg MultiPlotCfg) error {
	for i, pcfg := range mcfg.Cfgs {
		err := SlopeDiffFull(pcfg, fmt.Sprintf("%v_%v_slopes.bed", mcfg.Outpre, i))
		if err != nil {
			return err
		}
	}
	return nil
}

func ProcessMultiPlotCfgsSlopeDiff(cfgs []MultiPlotCfg, o SlopeDiffOptions) error {
	errs := pmap.Map(SlopeDiffFullMcfg, cfgs, o.Threads)
	for _, err := range errs {
		if err != nil {
			fmt.Fprintln(os.Stderr, errs)
			return err
		}
	}
	return nil
}

type SlopeDiffOptions struct {
	Threads int
}

func GetOptsSlopeDiff() SlopeDiffOptions {
	var o SlopeDiffOptions
	flag.IntVar(&o.Threads, "t", 1, "Threads")
	flag.Parse()
	return o
}

func RunSlopeDiffs() {
	opts := GetOptsSlopeDiff()
	mcfgs, err := ReadMultiPlotCfgs(os.Stdin)
	if err != nil {
		panic(err)
	}

	err = ProcessMultiPlotCfgsSlopeDiff(mcfgs, opts)
	if err != nil {
		panic(err)
	}
}

