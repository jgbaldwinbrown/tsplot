package main

import (
	"os"
	"github.com/jgbaldwinbrown/tsplot/tsplot"
	"flag"
)

func GetOpts() tsplot.SbiOptions {
	var opts tsplot.SbiOptions
	flag.BoolVar(&opts.WritePlottables, "w", true, "Do not write to plottable files (they already exist)")
	flag.BoolVar(&opts.Plot, "p", true, "Do not plot plottables")
	flag.Parse()
	if !opts.WritePlottables {
		opts.Plot = false
	}
	return opts
}

func main() {
	opts := GetOpts()
	syncBedInfo, err := tsplot.ReadSyncBedInfoCategories(os.Stdin)
	if err != nil {
		panic(err)
	}

	sbiSets := tsplot.SplitSbiByCategory(syncBedInfo)
	err = tsplot.ProcessSyncBedInfoSets(sbiSets, opts)
	if err != nil {
		panic(err)
	}
}
