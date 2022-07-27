package tsplot

import (
	"io"
	"github.com/jgbaldwinbrown/lscan/lscan"
	"bufio"
	"fmt"
)

type SyncBedInfo struct {
	Sync string
	Bed string
	Info string
	Out string
	Category string
}

type SbiOptions struct {
	WritePlottables bool
	Plot bool
	Threads int
}

func ReadSyncBedInfo(r io.Reader) ([]SyncBedInfo, error) {
	var out []SyncBedInfo

	s := bufio.NewScanner(r)
	scanf := lscan.ByByte('\t')
	var line []string

	for s.Scan() {
		line = lscan.SplitByFunc(line, s.Text(), scanf)
		if len(line) < 4 {
			return nil, fmt.Errorf("sbi line too short")
		}
		out = append(out, SyncBedInfo{line[0], line[1], line[2], line[3], ""})
	}
	return out, nil
}

func ReadSyncBedInfoCategories(r io.Reader) ([]SyncBedInfo, error) {
	var out []SyncBedInfo

	s := bufio.NewScanner(r)
	scanf := lscan.ByByte('\t')
	var line []string

	for s.Scan() {
		line = lscan.SplitByFunc(line, s.Text(), scanf)
		if len(line) < 5 {
			return nil, fmt.Errorf("sbi line too short")
		}
		out = append(out, SyncBedInfo{line[0], line[1], line[2], line[3], line[4]})
	}
	return out, nil
}

func SplitSbiByCategory(sbis []SyncBedInfo) [][]SyncBedInfo {
	var out [][]SyncBedInfo
	catsMap := make(map[string][]SyncBedInfo)
	for _, sbi := range sbis {
		list, _ := catsMap[sbi.Category]
		catsMap[sbi.Category] = append(list, sbi)
	}
	for _, sbis := range catsMap {
		out = append(out, sbis)
	}
	return out
}

func ProcessSyncBedInfo(sbi SyncBedInfo, o SbiOptions) error {
	bed, err := ReadBed(sbi.Bed)
	if err != nil {
		return err
	}

	sync, err := ReadSync(sbi.Sync, bed)
	if err != nil {
		return err
	}

	info, err := ReadInfo(sbi.Info)
	if err != nil {
		return err
	}

	plottables := ToSeparatePlottables(sync, bed, info, sbi.Out)

	if o.WritePlottables {
		err := WritePlottablesToFiles(plottables...)
		if err != nil {
			return err
		}
	}

	if o.Plot {
		err = PlotPlottables(plottables...)
		if err != nil {
			return err
		}
	}

	var paths []string
	for _, plottable := range plottables {
		paths = append(paths, plottable.Outprefix + plottablesuffix)
	}

	maxes, err := PathsToMaxes(paths)
	if err != nil {
		return err
	}

	for i, max := range maxes {
		fmt.Printf("%v\t%v\t%v\n", paths[i], max, sbi.Category)
	}

	return nil
}

func ProcessSyncBedInfos(sbis []SyncBedInfo, o SbiOptions) error {
	for _, sbi := range sbis {
		err := ProcessSyncBedInfo(sbi, o)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProcessSyncBedInfoSets(sbiSets [][]SyncBedInfo, o SbiOptions) error {
	for _, set := range sbiSets {
		err := ProcessSyncBedInfos(set, o)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProcessSyncBedInfoSetsParallel(sbiSets [][]SyncBedInfo, o SbiOptions) error {
	njobs := len(sbiSets)
	jobs := make(chan []SyncBedInfo, njobs)
	errs := make(chan error, njobs)
	nworkers := o.Threads
	for i:=0; i<nworkers; i++ {
		go func() {
			for job := range jobs {
				errs <- ProcessSyncBedInfos(job, o)
			}
		}()
	}
	for _, set := range sbiSets {
		jobs <- set
	}
	close(jobs)
	for i := 0; i<njobs; i++ {
		err := <-errs
		if err != nil {
			return err
		}
	}
	return nil
}
