package tsplot

import (
	"io"
	"encoding/json"
	"github.com/jgbaldwinbrown/pmap/pkg"
	"os/exec"
	"os"
	"bufio"
	"math/rand"
	"fmt"
)


func CalcSlope(sync SyncE, info []InfoE) float64 {
	maxgen := 49.0
	freqs := CalcAfFreqs(sync.Afs)
	gens := GetGens(info)
	return GetSlope(NewTableColIter(freqs, sync.PlotcolPrimary), NewSliceIter(gens), maxgen)
}

func CalcSlopes(sync []SyncE, info []InfoE) []float64 {
	var out []float64
	for _, synce := range sync {
		out = append(out, CalcSlope(synce, info))
	}
	return out
}

func PlottableSlopes(name string, slopes []float64) [][]string {
	var out [][]string
	for _, slope := range slopes {
		out = append(out, []string{name, fmt.Sprintf("%v", slope)})
	}
	return out
}

func ToSlopePlottable(name string, sync []SyncE, info []InfoE, expsync []SyncE, expinfo []InfoE, controlsync []SyncE, controlinfo []InfoE, toUse InfoSelection) [][]string {
	cols := GoodCols(info, toUse)

	newsync := SubsetSyncByCols(sync, cols)
	newinfo := SubsetInfoByCols(info, cols)

	expCols := GoodCols(expinfo, toUse)
	newExpSync := SubsetSyncByCols(expsync, expCols)
	newExpInfo := SubsetInfoByCols(expinfo, expCols)

	controlCols := GoodCols(controlinfo, toUse)
	newControlSync := SubsetSyncByCols(controlsync, controlCols)
	newControlInfo := SubsetInfoByCols(controlinfo, controlCols)

	SetBeneficialsHigh36(newsync, newinfo, newExpSync, newExpInfo, newControlSync, newControlInfo)
	return PlottableSlopes(name, CalcSlopes(newsync, newinfo))
}

func ToSlopePlottableAndFixationCounts(name string, sync []SyncE, info []InfoE, expsync []SyncE, expinfo []InfoE, controlsync []SyncE, controlinfo []InfoE, toUse InfoSelection) ([][]string, []FixationCount) {
	cols := GoodCols(info, toUse)

	newsync := SubsetSyncByCols(sync, cols)
	newinfo := SubsetInfoByCols(info, cols)

	expCols := GoodCols(expinfo, toUse)
	newExpSync := SubsetSyncByCols(expsync, expCols)
	newExpInfo := SubsetInfoByCols(expinfo, expCols)

	controlCols := GoodCols(controlinfo, toUse)
	newControlSync := SubsetSyncByCols(controlsync, controlCols)
	newControlInfo := SubsetInfoByCols(controlinfo, controlCols)

	SetBeneficialsHigh36(newsync, newinfo, newExpSync, newExpInfo, newControlSync, newControlInfo)
	slopes := PlottableSlopes(name, CalcSlopes(newsync, newinfo))
	counts := ToFixationCounts(name, newsync, newinfo)
	return slopes, counts
}

func GetPcfgSlopePlottables(pcfg SlopePlotCfg) (slopePlottable [][]string, err error) {
	expsync, _, expinfo, err := ReadSBI(pcfg.PolarizerExpSbi)
	if err != nil {
		return nil, err
	}

	controlsync, _, controlinfo, err := ReadSBI(pcfg.PolarizerControlSbi)
	if err != nil {
		return nil, err
	}

	for i, sbi := range pcfg.Sbis {
		sync, _, info, err := ReadSBI(sbi)
		if err != nil {
			return nil, err
		}
		slopePlottable = append(slopePlottable, ToSlopePlottable(sbi.Category, sync, info, expsync, expinfo, controlsync, controlinfo, pcfg.ToUses[i])...)
	}
	return slopePlottable, nil
}

func GetPcfgSlopePlottablesAndFixationCounts(pcfg SlopePlotCfg) (slopePlottable [][]string, fixationCounts []FixationCount, err error) {
	expsync, _, expinfo, err := ReadSBI(pcfg.PolarizerExpSbi)
	if err != nil {
		return nil, nil, err
	}

	controlsync, _, controlinfo, err := ReadSBI(pcfg.PolarizerControlSbi)
	if err != nil {
		return nil, nil, err
	}

	for i, sbi := range pcfg.Sbis {
		sync, _, info, err := ReadSBI(sbi)
		if err != nil {
			return nil, nil, err
		}
		oneSlopePlottable, oneCount := ToSlopePlottableAndFixationCounts(sbi.Category, sync, info, expsync, expinfo, controlsync, controlinfo, pcfg.ToUses[i])
		slopePlottable = append(slopePlottable, oneSlopePlottable...)
		fixationCounts = append(fixationCounts, oneCount...)
	}
	return slopePlottable, fixationCounts, nil
}

func GetPcfgPlottables(pcfg SlopePlotCfg) (plottable [][]string, err error) {
	sync, _, info, err := ReadSBI(pcfg.Sbis[0])
	if err != nil {
		return nil, err
	}

	expsync, _, expinfo, err := ReadSBI(pcfg.PolarizerExpSbi)
	if err != nil {
		return nil, err
	}
	controlsync, _, controlinfo, err := ReadSBI(pcfg.PolarizerControlSbi)
	if err != nil {
		return nil, err
	}

	plottable = append(plottable, ToPlottableBeneficialG36Subset(sync, info, expsync, expinfo, controlsync, controlinfo, pcfg.ToUses[0])...)
	return plottable, nil
}

func GetAllPcfgPlottables(cfg MultiSlopePlotCfg) (slopePlottable [][]string, plottable [][]string, err error) {
	for _, pcfg := range cfg.Cfgs {
		onePlottable, err := GetPcfgPlottables(pcfg)
		if err != nil {
			return nil, nil, err
		}
		oneSlopePlottable, err := GetPcfgSlopePlottables(pcfg)
		if err != nil {
			return nil, nil, err
		}

		slopePlottable = append(slopePlottable, oneSlopePlottable...)
		plottable = append(plottable, onePlottable...)
	}

	return slopePlottable, plottable, nil
}

func GetAllPcfgPlottablesAndFixationCounts(cfg MultiSlopePlotCfg) (slopePlottable [][]string, plottable [][]string, fixCounts []FixationCount, err error) {
	for _, pcfg := range cfg.Cfgs {
		onePlottable, err := GetPcfgPlottables(pcfg)
		if err != nil {
			return nil, nil, nil, err
		}
		oneSlopePlottable, oneFixCount, err := GetPcfgSlopePlottablesAndFixationCounts(pcfg)
		if err != nil {
			return nil, nil, nil, err
		}

		slopePlottable = append(slopePlottable, oneSlopePlottable...)
		plottable = append(plottable, onePlottable...)
		fixCounts = append(fixCounts, oneFixCount...)
	}

	return slopePlottable, plottable, fixCounts, nil
}

func GetInverseSitesBed(selectedPlottable [][]string, pcfg SlopePlotCfg, seed int64) ([]BedE, error) {
	nsites := CountSites(selectedPlottable, pcfg.ToUses[0])

	r := rand.New(rand.NewSource(seed))
	sitessync, err := ChooseSitesFullGeneric(nsites, pcfg.PolarizerExpSbi, r)

	if err != nil {
		return nil, err
	}
	sitesbed := SyncToBed(sitessync)
	return sitesbed, nil
}

const unselsuffix = "_unsel"

func GetInPopUnselSlopePlottable(selectedPlottable [][]string, cfg MultiSlopePlotCfg, seed int64) ([][]string, error) {

	var plottable [][]string
	for _, pcfg := range cfg.Cfgs {
		sitesbed, err := GetInverseSitesBed(selectedPlottable, pcfg, seed);

		expsync, expinfo, err := ReadSyncInfo(pcfg.PolarizerExpSbi.Sync, pcfg.PolarizerExpSbi.Info, sitesbed)
		if err != nil {
			return nil, err
		}

		controlsync, controlinfo, err := ReadSyncInfo(pcfg.PolarizerControlSbi.Sync, pcfg.PolarizerControlSbi.Info, sitesbed)
		if err != nil {
			return nil, err
		}

		for i, sbi := range pcfg.Sbis {
			sync, info, err := ReadSyncInfo(sbi.Sync, sbi.Info, sitesbed)
			if err != nil {
				return nil, err
			}

			plottable = append(plottable, ToSlopePlottable(sbi.Category + unselsuffix, sync, info, expsync, expinfo, controlsync, controlinfo, pcfg.ToUses[i])...)
		}
	}
	return plottable, nil
}

func GetInPopUnselSlopePlottableAndFixationCounts(selectedPlottable [][]string, cfg MultiSlopePlotCfg, seed int64) ([][]string, []FixationCount, error) {

	var plottable [][]string
	var counts []FixationCount
	for _, pcfg := range cfg.Cfgs {
		sitesbed, err := GetInverseSitesBed(selectedPlottable, pcfg, seed);

		expsync, expinfo, err := ReadSyncInfo(pcfg.PolarizerExpSbi.Sync, pcfg.PolarizerExpSbi.Info, sitesbed)
		if err != nil {
			return nil, nil, err
		}

		controlsync, controlinfo, err := ReadSyncInfo(pcfg.PolarizerControlSbi.Sync, pcfg.PolarizerControlSbi.Info, sitesbed)
		if err != nil {
			return nil, nil, err
		}

		for i, sbi := range pcfg.Sbis {
			sync, info, err := ReadSyncInfo(sbi.Sync, sbi.Info, sitesbed)
			if err != nil {
				return nil, nil, err
			}

			onePlottable, oneCount := ToSlopePlottableAndFixationCounts(sbi.Category + unselsuffix, sync, info, expsync, expinfo, controlsync, controlinfo, pcfg.ToUses[i])
			plottable = append(plottable, onePlottable...)
			counts = append(counts, oneCount...)
		}
	}
	return plottable, counts, nil
}


func GetSlopePlottables(cfg MultiSlopePlotCfg, seed int64) (slopePlottable [][]string, err error) {
	slopePlottable, plottable, err := GetAllPcfgPlottables(cfg)
	if err != nil {
		return nil, err
	}

	unselSlopePlottable, err := GetInPopUnselSlopePlottable(plottable, cfg, seed)
	if err != nil {
		return nil, err
	}

	slopePlottable = append(slopePlottable, unselSlopePlottable...)
	return slopePlottable, nil
}

func GetSlopeAndFixPlottables(cfg MultiSlopePlotCfg, seed int64) (slopePlottable [][]string, fixCounts []FixationCount, err error) {
	slopePlottable, plottable, fixCounts, err := GetAllPcfgPlottablesAndFixationCounts(cfg)
	if err != nil {
		return nil, nil, err
	}

	unselSlopePlottable, unselFixCounts, err := GetInPopUnselSlopePlottableAndFixationCounts(plottable, cfg, seed)
	if err != nil {
		return nil, nil, err
	}

	slopePlottable = append(slopePlottable, unselSlopePlottable...)
	fixCounts = append(fixCounts, unselFixCounts...)
	return slopePlottable, fixCounts, nil
}

func WriteSlopePlottableToPath(plottable [][]string, path string) error {
	w, err := os.Create(path)
	if err != nil {
		return err
	}
	defer w.Close()
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	fmt.Fprintf(bw, "Treatment\tSlope\n")
	for _, line := range plottable {
		if len(line) > 1 {
			fmt.Fprintf(bw, "%v\t%v\n", line[0], line[1])
		}
	}

	return nil
}

func PlotSlopeBox(inpath, outpath string) error {
	cmd := exec.Command("plotslopebox", inpath, outpath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

const slopecalcsuffix = "_slopecalc"
const fixsuffix = "_fixations"
const plotsuffix = "_plotted"

func ProcessMultiPlotSlopeCalcsOld(cfg MultiSlopePlotCfg, o SbiOptions) error {

	plottablepath := cfg.Outpre + slopecalcsuffix + ".bed"
	plotpath := cfg.Outpre + slopecalcsuffix + plotsuffix + ".pdf"
	if !o.NoWritePlottables {
		slopePlottable, err := GetSlopePlottables(cfg, int64(o.Seed))
		if err != nil {
			return err
		}

		err = WriteSlopePlottableToPath(slopePlottable, plottablepath)
		if err != nil {
			return err
		}
	}

	if !o.NoPlot {
		err := PlotSlopeBox(plottablepath, plotpath)
		if err != nil {
			return err
		}
	}

	return nil
}

func ProcessMultiPlotSlopeAndFixCalcs(cfg MultiSlopePlotCfg, o SbiOptions) error {
	plottablepath := cfg.Outpre + slopecalcsuffix + ".bed"
	plotpath := cfg.Outpre + slopecalcsuffix + plotsuffix + ".pdf"
	fixpath := cfg.Outpre + fixsuffix + ".bed"
	fixplotpath := cfg.Outpre + fixsuffix + plotsuffix + ".bed"


	if !o.NoWritePlottables {
		slopePlottable, fixcounts, err := GetSlopeAndFixPlottables(cfg, int64(o.Seed))
		if err != nil {
			return err
		}

		err = WriteSlopePlottableToPath(slopePlottable, plottablepath)
		if err != nil {
			return err
		}

		err = WriteFixationCounts(fixcounts, fixpath)
		if err != nil {
			return err
		}
	}

	if !o.NoPlot {
		err := PlotSlopeBox(plottablepath, plotpath + ".pdf")
		if err != nil {
			return err
		}

		err = PlotFixCurve(fixpath, fixplotpath + ".pdf")
		if err != nil {
			return err
		}
	}

	return nil
}

func ProcessMultiSlopePlotCfgs(cfgs []MultiSlopePlotCfg, o SbiOptions) error {
	f := func(cfg MultiSlopePlotCfg) error {
		return ProcessMultiPlotSlopeAndFixCalcs(cfg, o)
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

func ReadMultiSlopePlotCfgs(r io.Reader) ([]MultiSlopePlotCfg, error) {
	var out []MultiSlopePlotCfg
	dec := json.NewDecoder(r)
	err := dec.Decode(&out)
	return out, err
}

func RunMultiSlopePlotCfgs() {
	opts := GetOpts()
	mcfgs, err := ReadMultiSlopePlotCfgs(os.Stdin)
	if err != nil {
		panic(err)
	}

	err = ProcessMultiSlopePlotCfgs(mcfgs, opts)
	if err != nil {
		panic(err)
	}
}

