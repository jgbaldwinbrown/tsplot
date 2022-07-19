package main

import (
	"os"
	"bufio"
	"flag"
	"github.com/jgbaldwinbrown/tsplot/tsplot"
)

type Flags struct {
	SyncPath string
	BedPath string
	InfoPath string
}

func GetFlags() Flags {
	f := Flags{}
	flag.StringVar(&f.SyncPath, "s", "", "Path to sync file with allele frequencies")
	flag.StringVar(&f.BedPath, "b", "", "Path to bed file with regions to include")
	flag.StringVar(&f.InfoPath, "i", "", "Path to info file with identities of sync columns")
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

	plottable := tsplot.ToPlottable(sync, info)

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	tsplot.WritePlottable(w, plottable)
}
