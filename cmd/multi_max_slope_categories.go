package main

import (
	"os"
	"github.com/jgbaldwinbrown/tsplot/tsplot"
	"flag"
)

func GetOpts() tsplot.SbiOptions {
	var opts tsplot.SbiOptions
	flag.BoolVar(&opts.NoWritePlottables, "w", false, "Do not write to plottable files (they already exist)")
	flag.BoolVar(&opts.NoPlot, "p", false, "Do not plot plottables")
	flag.IntVar(&opts.Threads, "t", 1, "Threads to use")
	flag.Parse()

	return opts
}

func main() {
	opts := GetOpts()
	syncBedInfo, err := tsplot.ReadSyncBedInfoCategories(os.Stdin)
	if err != nil {
		panic(err)
	}

	sbiSets := tsplot.SplitSbiByCategory(syncBedInfo)
	err = tsplot.ProcessSyncBedInfoSetsParallel(sbiSets, opts)
	if err != nil {
		panic(err)
	}
}
