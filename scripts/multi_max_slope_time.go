package main

import (
	"os"
	"github.com/jgbaldwinbrown/tsplot/tsplot"
)

func main() {
	syncBedInfo, err := tsplot.ReadSyncBedInfo(os.Stdin)
	if err != nil {
		panic(err)
	}

	for _, sbi := range syncBedInfo {
		err := tsplot.ProcessSyncBedInfo(sbi)
		if err != nil {
			panic(err)
		}
	}
}
