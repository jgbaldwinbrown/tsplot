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

const slopecalcsuffix = "_slopecalc"
const slopecalcplotsuffix = "_slopecalc_plotted"

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

			plottable = append(plottable, ToSlopePlottable(sbi.Category, sync, info, expsync, expinfo, controlsync, controlinfo, pcfg.ToUses[i])...)
		}
	}
	return plottable, nil
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

func WriteSlopePlottableToPath(plottable [][]string, path string) error {
	w, err := os.Open(path)
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

func ProcessMultiPlotSlopeCalcs(cfg MultiSlopePlotCfg, o SbiOptions) error {
	plottablepath := cfg.Outpre + slopecalcsuffix

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

	plotpath := plottablepath + slopecalcplotsuffix
	if !o.NoPlot {
		err := PlotSlopeBox(plottablepath, plotpath)
		if err != nil {
			return err
		}
	}

	return nil
}

func ProcessMultiSlopePlotCfgs(cfgs []MultiSlopePlotCfg, o SbiOptions) error {
	f := func(cfg MultiSlopePlotCfg) error {
		return ProcessMultiPlotSlopeCalcs(cfg, o)
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

