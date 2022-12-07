package tsplot

import (
	"encoding/json"
	"fmt"
)

type SlopePlotCfg struct {
	Sbis []SyncBedInfo
	ToUses []InfoSelection
	PolarizerExpSbi SyncBedInfo
	PolarizerExpToUse InfoSelection
	PolarizerControlSbi SyncBedInfo
	PolarizerControlToUse InfoSelection
}

type MultiSlopePlotCfg struct {
	Cfgs []SlopePlotCfg
	Outpre string
}

func BuildSlopeCfg(breed string, reps []int) MultiSlopePlotCfg {
	var m MultiSlopePlotCfg
	m.Outpre = fmt.Sprintf("%v_pfst_%v_peaks_%v_multiplot_slopecalc", breed, breed, "all4reps")

	for _, rep := range reps {
		cfg := SlopePlotCfg {
			Sbis: []SyncBedInfo {
				SyncBedInfo {
					Sync: SyncString(breed, "unbitted"),
					Bed: BedString(breed, "unbitted", rep),
					Info: InfoString(breed, "unbitted"),
					Category: fmt.Sprintf("%v_unbitted", breed),
				},
				SyncBedInfo {
					Sync: SyncString(breed, "bitted"),
					Bed: BedString(breed, "unbitted", rep),
					Info: InfoString(breed, "bitted"),
					Category: fmt.Sprintf("%v_bitted", breed),
				},
				SyncBedInfo {
					Sync: SyncString("feral", "unbitted"),
					Bed: BedString(breed, "unbitted", rep),
					Info: InfoString("feral", "unbitted"),
					Category: "feral_unbitted",
				},
				SyncBedInfo {
					Sync: SyncString("feral", "bitted"),
					Bed: BedString(breed, "unbitted", rep),
					Info: InfoString("feral", "bitted"),
					Category: "feral_bitted",
				},
			},
			ToUses: []InfoSelection {
				InfoSelection{[]string{fmt.Sprintf("%d", rep)}},
				InfoSelection{[]string{fmt.Sprintf("%d", rep)}},
				InfoSelection{[]string{fmt.Sprintf("%d", rep)}},
				InfoSelection{[]string{fmt.Sprintf("%d", rep)}},
			},

			PolarizerExpSbi: SyncBedInfo {
				Sync: SyncString(breed, "unbitted"),
				Bed: BedString(breed, "unbitted", rep),
				Info: InfoString(breed, "unbitted"),
				Category: fmt.Sprintf("%v_unbitted", breed),
			},
			PolarizerExpToUse: InfoSelection{[]string{fmt.Sprintf("%d", rep)}},

			PolarizerControlSbi: SyncBedInfo {
				Sync: SyncString("feral", "unbitted"),
				Bed: BedString(breed, "unbitted", rep),
				Info: InfoString("feral", "unbitted"),
				Category: "feral_unbitted",
			},
			PolarizerControlToUse: InfoSelection{[]string{fmt.Sprintf("%d", rep)}},
		}
		m.Cfgs = append(m.Cfgs, cfg)
	}
	return m
}
func MakeSlopeCfg() {
	breeds := []string{"black", "white", "figurita", "runt"}
	reps := []int{1,2,3,4}
	cfgs := []MultiSlopePlotCfg{}

	for _, breed := range breeds {
		cfgs = append(cfgs, BuildSlopeCfg(breed, reps))
	}

	out, err := json.Marshal(cfgs)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}
