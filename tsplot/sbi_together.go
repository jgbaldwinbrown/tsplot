package tsplot

import (
	"flag"
	"io"
	"os"
	"github.com/jgbaldwinbrown/pmap/pkg"
	"fmt"
	"golang.org/x/exp/slices"
	"encoding/json"
)

func GetOpts() SbiOptions {
	var opts SbiOptions
	flag.BoolVar(&opts.NoWritePlottables, "w", false, "Do not write to plottable files (they already exist)")
	flag.BoolVar(&opts.NoPlot, "p", false, "Do not plot plottables")
	flag.BoolVar(&opts.PlotBeneficial, "b", false, "Plot beneficial alleles, not minor alleles")
	flag.BoolVar(&opts.PlotBeneficialG36, "B", false, "Plot beneficial alleles based on generation 36 pFST comparisons, not minor alleles")
	flag.BoolVar(&opts.PlotUgly, "u", false, "Plot ugly-style, i.e. lots of extra lines")
	flag.IntVar(&opts.Threads, "t", 1, "Threads to use")
	flag.Parse()

	return opts
}

type InfoSelection struct {
	Reps []string
}

type PlotCfg struct {
	Sbi SyncBedInfo
	ToUse InfoSelection
	BeneficialSbi SyncBedInfo
	BeneficialToUse InfoSelection
	BeneficialExpSbi SyncBedInfo
	BeneficialExpToUse InfoSelection
	BeneficialControlSbi SyncBedInfo
	BeneficialControlToUse InfoSelection
}

type MultiPlotCfg struct {
	Cfgs []PlotCfg
	Outpre string
}

func SubsetSyncEByCols(s SyncE, cols []int) SyncE {
	var out SyncE
	out.Chr = s.Chr
	out.Pos = s.Pos
	for _, col := range cols {
		out.Afs = append(out.Afs, slices.Clone(s.Afs[col]))
		out.Line = append(out.Line, s.Line[col])
	}
	return out
}

func SubsetSyncByCols(s []SyncE, cols []int) []SyncE {
	var out []SyncE
	for _, se := range s {
		out = append(out, SubsetSyncEByCols(se, cols))
	}
	return out
}

func SubsetInfoByCols(info []InfoE, cols []int) []InfoE {
	var out []InfoE
	for _, col := range cols {
		out = append(out, info[col])
	}
	return out
}

func GoodCols(info []InfoE, toUse InfoSelection) []int {
	goodreps := map[string]struct{}{}
	for _, rep := range toUse.Reps {
		goodreps[rep] = struct{}{}
	}

	var cols []int
	for i, infoe := range info {
		if _, ok := goodreps[infoe.Repl]; ok {
			cols = append(cols, i)
		}
	}

	return cols
}

func ToPlottableSubset(sync []SyncE, info []InfoE, toUse InfoSelection) [][]string {
	cols := GoodCols(info, toUse)
	newsync := SubsetSyncByCols(sync, cols)
	newinfo := SubsetInfoByCols(info, cols)
	return ToPlottable(newsync, newinfo)
}

func ToPlottableBeneficialSubset(sync []SyncE, info []InfoE, beneSync []SyncE, beneInfo []InfoE, toUse InfoSelection) [][]string {
	cols := GoodCols(info, toUse)
	newsync := SubsetSyncByCols(sync, cols)
	newinfo := SubsetInfoByCols(info, cols)

	beneCols := GoodCols(beneInfo, toUse)
	newBeneSync := SubsetSyncByCols(beneSync, beneCols)
	newBeneInfo := SubsetInfoByCols(beneInfo, beneCols)

	SetBeneficials(newsync, newinfo, newBeneSync, newBeneInfo)
	return ToPlottablePlotcol(newsync, newinfo)
}

func ToPlottableBeneficialG36Subset(sync []SyncE, info []InfoE, beneExpSync []SyncE, beneExpInfo []InfoE, beneControlSync []SyncE, beneControlInfo []InfoE, toUse InfoSelection) [][]string {
	cols := GoodCols(info, toUse)

	newsync := SubsetSyncByCols(sync, cols)
	newinfo := SubsetInfoByCols(info, cols)

	beneExpCols := GoodCols(beneExpInfo, toUse)
	newBeneExpSync := SubsetSyncByCols(beneExpSync, beneExpCols)
	newBeneExpInfo := SubsetInfoByCols(beneExpInfo, beneExpCols)

	beneControlCols := GoodCols(beneControlInfo, toUse)
	newBeneControlSync := SubsetSyncByCols(beneControlSync, beneControlCols)
	newBeneControlInfo := SubsetInfoByCols(beneControlInfo, beneControlCols)

	SetBeneficialsHigh36(newsync, newinfo, newBeneExpSync, newBeneExpInfo, newBeneControlSync, newBeneControlInfo)
	return ToPlottablePlotcol(newsync, newinfo)
}

func ProcessMultiPlotCfg(cfg MultiPlotCfg, o SbiOptions) error {
	var plottable [][]string
	plottablepath := cfg.Outpre + plottablesuffix

	if !o.NoWritePlottables {
		for _, pcfg := range cfg.Cfgs {
			sync, _, info, err := ReadSBI(pcfg.Sbi)
			if err != nil {
				return err
			}
			if o.PlotBeneficial {
				benesync, _, beneinfo, err := ReadSBI(pcfg.BeneficialSbi)
				if err != nil {
					return err
				}
				plottable = append(plottable, ToPlottableBeneficialSubset(sync, info, benesync, beneinfo, pcfg.ToUse)...)
			} else if o.PlotBeneficialG36 {
				expsync, _, expinfo, err := ReadSBI(pcfg.BeneficialExpSbi)
				if err != nil {
					return err
				}
				controlsync, _, controlinfo, err := ReadSBI(pcfg.BeneficialControlSbi)
				if err != nil {
					return err
				}
				plottable = append(plottable, ToPlottableBeneficialG36Subset(sync, info, expsync, expinfo, controlsync, controlinfo, pcfg.ToUse)...)
			} else {
				plottable = append(plottable, ToPlottableSubset(sync, info, pcfg.ToUse)...)
			}
		}


		err := WritePlottableToFile(plottable, cfg.Outpre)
		if err != nil {
			return err
		}
	}

	plotpath := cfg.Outpre + plottedsuffix
	if !o.NoPlot {
		err := PlotPlottableFileTip(plottablepath, plotpath, o.PlotUgly)
		if err != nil {
			return err
		}
	}

	paths := []string{plottablepath}

	maxes, err := PathsToMaxes(paths)
	if err != nil {
		return err
	}

	for i, max := range maxes {
		fmt.Printf("%v\t%v\t%v\n", paths[i], max, i)
	}

	return nil
}

func ProcessMultiPlotCfgs(cfgs []MultiPlotCfg, o SbiOptions) error {
	f := func(cfg MultiPlotCfg) error {
		return ProcessMultiPlotCfg(cfg, o)
	}

	errs := pmap.Map(f, cfgs, o.Threads)
	for _, err := range errs {
		if err != nil {
			fmt.Fprintln(os.Stderr, errs)
			return err
		}
	}
	return nil
}

func ReadMultiPlotCfgs(r io.Reader) ([]MultiPlotCfg, error) {
	var out []MultiPlotCfg
	dec := json.NewDecoder(r)
	err := dec.Decode(&out)
	return out, err
}

func RunMultiPlotCfgs() {
	opts := GetOpts()
	mcfgs, err := ReadMultiPlotCfgs(os.Stdin)
	if err != nil {
		panic(err)
	}

	err = ProcessMultiPlotCfgs(mcfgs, opts)
	if err != nil {
		panic(err)
	}
}
