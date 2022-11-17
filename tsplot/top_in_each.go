package tsplot

/*
import (
	"os"
	"flag"
)

func GetPeakmaxOpts() SbiOptions {
	var opts tsplot.SbiOptions
	flag.BoolVar(&opts.WritePlottables, "w", true, "Do not write to plottable files (they already exist)")
	flag.BoolVar(&opts.Plot, "p", true, "Do not plot plottables")
	flag.IntVar(&opts.Threads, "t", 1, "Threads to use")
	flag.Parse()

	return opts
}

func RunPeakmax() {
	opts := GetPeakmaxOpts()
	syncBedInfo, err := tsplot.ReadSyncBedInfoCategories(os.Stdin)
	if err != nil {
		panic(err)
	}

	sbiSets := SplitSbiByCategory(syncBedInfo)

	err := PlotSbiPeaksParallel(sbiSets, opts)
	if err != nil {
		panic(err)
	}

	err = PlotSbiPeaksJoin(sbiSets, opts)
	if err != nil {
		panic(err)
	}
}
*/

func ReadSBI(sbi SyncBedInfo) ([]SyncE, []BedE, []InfoE, error) {
	bed, err := ReadBed(sbi.Bed)
	if err != nil {
		return nil, nil, nil, err
	}

	sync, err := ReadSync(sbi.Sync, bed)
	if err != nil {
		return nil, nil, nil, err
	}

	info, err := ReadInfo(sbi.Info)
	if err != nil {
		return nil, nil, nil, err
	}
	return sync, bed, info, nil
}

/*

func MaxInPlottables(plottables []Plottable, ) []Plottable {
}

func PlotSbiPeaks(sbi SyncBedInfo, o SbiOptions) error {
	sync, bed, info, err := ReadSBI(sbi)
	if err != nil {
		return err
	}

	maxed_bed := 
	plottables := ToSeparatePlottables(sync, bed, info, sbi.Out)
	plottables_maxes := MaxInPlottables(plottables, bed)
	plottables_joined := JoinPlottables(plottables...)

	if o.WritePlottables {
		if err := WritePlottablesToFiles(plottables_joined...); err != nil {
			return err
		}
	}

	if o.Plot {
		if err := PlotPlottables(plottables_joined...); err != nil {
			return err
		}
	}

	path := plottables_joined.Outprefix + plottablesuffix)
	maxes, err := PathsToMaxes([]string{paths})
	if err != nil {
		return err
	}

	for i, max := range maxes {
		fmt.Printf("%v\t%v\t%v\n", paths[i], max, sbi.Category)
	}

	return nil
}

*/
