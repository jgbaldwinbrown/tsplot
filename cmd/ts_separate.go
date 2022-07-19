package main

import (
	"flag"
	"../tsplot"
)

type Flags struct {
	SyncPath string
	BedPath string
	InfoPath string
	OutputPrefix string
}

func GetFlags() Flags {
	f := Flags{}
	flag.StringVar(&f.SyncPath, "s", "", "Path to sync file with allele frequencies")
	flag.StringVar(&f.BedPath, "b", "", "Path to bed file with regions to include")
	flag.StringVar(&f.InfoPath, "i", "", "Path to info file with identities of sync columns")
	flag.StringVar(&f.OutputPrefix, "o", "separated", "Prefix for output files")
	flag.Parse()
	return f
}

func main() {
	f := GetFlags()

	info, err := tsplot.ReadInfo(f.InfoPath)
	if err != nil {
		panic(err)
	}

	bed, err := tsplot.ReadBed(f.BedPath)
	if err != nil {
		panic(err)
	}

	sync, err := tsplot.ReadSync(f.SyncPath, bed)
	if err != nil {
		panic(err)
	}

	plottables := tsplot.ToSeparatePlottables(sync, bed, info, f.OutputPrefix)

	err = tsplot.PlotPlottables(plottables...)
	if err != nil {
		panic(err)
	}
}
