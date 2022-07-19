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

func ProcessSyncBedInfo(sbi SyncBedInfo) error {
	bed, err := ReadBed(sbi.Bed)
	if err != nil {
		panic(err)
	}

	sync, err := ReadSync(sbi.Sync, bed)
	if err != nil {
		panic(err)
	}

	info, err := ReadInfo(sbi.Info)
	if err != nil {
		panic(err)
	}

	plottables := ToSeparatePlottables(sync, bed, info, sbi.Out)
	err = PlotPlottables(plottables...)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, plottable := range plottables {
		paths = append(paths, plottable.Outprefix + plottablesuffix)
	}

	maxes, err := PathsToMaxes(paths)
	for i, max := range maxes {
		fmt.Printf("%v\t%v\t%v\n", paths[i], max, sbi.Category)
	}

	return nil
}

func ProcessSyncBedInfos(sbis []SyncBedInfo) error {
	for _, sbi := range sbis {
		err := ProcessSyncBedInfo(sbi)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProcessSyncBedInfoSets(sbiSets [][]SyncBedInfo) error {
	for _, set := range sbiSets {
		err := ProcessSyncBedInfos(set)
		if err != nil {
			return err
		}
	}
	return nil
}
