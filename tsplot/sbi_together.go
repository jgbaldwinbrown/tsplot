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

func ToPlottableBeneficialSubset(sync []SyncE, info []InfoE, toUse InfoSelection) [][]string {
	cols := GoodCols(info, toUse)
	newsync := SubsetSyncByCols(sync, cols)
	newinfo := SubsetInfoByCols(info, cols)
	SetBeneficials(newsync, info)
	return ToPlottablePlotcol(newsync, newinfo)
}


func ProcessMultiPlotCfg(cfg MultiPlotCfg, o SbiOptions) error {
	var plottable [][]string
	plottablepath := cfg.Outpre + plottablesuffix

	if !o.NoWritePlottables {
		fmt.Fprintln(os.Stderr, "o.NoWritePlottables is false")
		for _, pcfg := range cfg.Cfgs {
			sync, _, info, err := ReadSBI(pcfg.Sbi)
			if err != nil {
				return err
			}
			if o.PlotBeneficial {
				plottable = append(plottable, ToPlottableBeneficialSubset(sync, info, pcfg.ToUse)...)
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
