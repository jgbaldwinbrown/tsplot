package main

import (
	"os"
	"github.com/jgbaldwinbrown/tsplot/tsplot"
)

func main() {
	syncBedInfo, err := tsplot.ReadSyncBedInfoCategories(os.Stdin)
	if err != nil {
		panic(err)
	}

	sbiSets := tsplot.SplitSbiByCategory(syncBedInfo)
	err = tsplot.ProcessSyncBedInfoSets(sbiSets)
	if err != nil {
		panic(err)
	}
}
