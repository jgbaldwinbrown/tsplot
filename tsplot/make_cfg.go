package tsplot

import (
	"flag"
	"encoding/json"
	"fmt"
)

func SyncString(breed, bit string) string {
	return fmt.Sprintf("split_all_pools/sync/sync2/%v_tall_names_%v_repall.split_goods_f0.1_c10.sync.gz", breed, bit)
}

func Reformat(in string) string {
	out, ok := map[string]string {
		"black": "BlackHomer",
		"white": "WhiteHomer",
		"figurita": "Figurita",
		"runt": "Runt",
		"feral": "Feral",
		"bitted": "Bitted",
		"unbitted": "Unbitted",
	}[in]

	if !ok {
		return in
	}
	return out
}

func RepReformat(in int) string {
	out, ok := map[int]string {
		0: "All",
		1: "R1",
		2: "R2",
		3: "R3",
		4: "R4",
	}[in]

	if !ok {
		return ""
	}
	return out
}

func Treatment(breed string) string {
	out, ok := map[string]string {
		"black": "Color",
		"white": "Color",
		"figurita": "Size",
		"runt": "Size",
	}[breed]

	if !ok {
		return breed
	}
	return out
}

func BedString(breed, bit string, rep int) string {
	breedrefmt := Reformat(breed)
	bitrefmt := Reformat(bit)
	reprefmt := RepReformat(rep)
	breedtreat := Treatment(breed)
	return fmt.Sprintf(
		"/media/jgbaldwinbrown/3564-3063/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/vcftools/plot_pfst/run4_10k/_breed_%v_time_36_bit_%v_replicate_%v_breed_Feral_time_36_bit_%v_replicate_%v_%v_Full_Full__win10000_1000_multiplot_subtractionsbed.bed_pFst_tips.bed",
		breedrefmt, bitrefmt, reprefmt, bitrefmt, reprefmt, breedtreat,
	)
}
// /media/jgbaldwinbrown/3564-3063/jgbaldwinbrown/Documents/work_stuff/louse/poolfstat/vcftools/plot_pfst/run4_10k/_breed_BlackHomer_time_36_bit_Unbitted_replicate_All_breed_Feral_time_36_bit_Unbitted_replicate_All_Color_Full_Full__win10000_1000_multiplot_subtractionsbed.bed_pFst_tips.bed

func InfoString(breed, bit string) string {
	return fmt.Sprintf("split_all_pools/sync/sync2/%v_pooled_info_ne_%v.txt", breed, bit)
}

func OutprefixString(breed, bit string) string {
	return fmt.Sprintf("%v_%v_separated", breed, bit)
}

func BuildCfg(breed, bit string, reps []int, beneficial, beneficial36, ugly bool) MultiPlotCfg {
	var m MultiPlotCfg
	m.Outpre = fmt.Sprintf("%v_pfst_%v_peaks_%v_%v_multiplot", breed, breed, bit, "all4reps")
	if beneficial {
		m.Outpre = m.Outpre + "_beneficial"
	}
	if beneficial36 {
		m.Outpre = m.Outpre + "_beneficial36"
	}
	if ugly {
		m.Outpre = m.Outpre + "_ugly"
	}
	for _, rep := range reps {
		cfg := PlotCfg {
			Sbi: SyncBedInfo {
				Sync: SyncString(breed, bit),
				Bed: BedString(breed, "unbitted", rep),
				Info: InfoString(breed, bit),
				Out: "",
				Category: "",
			},
			ToUse: InfoSelection{[]string{fmt.Sprintf("%d", rep)}},

			BeneficialSbi: SyncBedInfo {
				Sync: SyncString(breed, "unbitted"),
				Bed: BedString(breed, "unbitted", rep),
				Info: InfoString(breed, "unbitted"),
				Out: "",
				Category: "",
			},
			BeneficialToUse: InfoSelection{[]string{fmt.Sprintf("%d", rep)}},

			BeneficialExpSbi: SyncBedInfo {
				Sync: SyncString(breed, "unbitted"),
				Bed: BedString(breed, "unbitted", rep),
				Info: InfoString(breed, "unbitted"),
				Out: "",
				Category: "",
			},
			BeneficialExpToUse: InfoSelection{[]string{fmt.Sprintf("%d", rep)}},

			BeneficialControlSbi: SyncBedInfo {
				Sync: SyncString("feral", "unbitted"),
				Bed: BedString(breed, "unbitted", rep),
				Info: InfoString("feral", "unbitted"),
				Out: "",
				Category: "",
			},
			BeneficialControlToUse: InfoSelection{[]string{fmt.Sprintf("%d", rep)}},
		}
		m.Cfgs = append(m.Cfgs, cfg)
	}
	return m
}

func BuildControl(breed, bit string, reps []int, beneficial, beneficial36, ugly bool) MultiPlotCfg {
	var m MultiPlotCfg
	m.Outpre = fmt.Sprintf("%v_pfst_%v_peaks_%v_%v_multiplot", "feral", breed, bit, "all4reps")
	if beneficial {
		m.Outpre = m.Outpre + "_beneficial"
	}
	if beneficial36 {
		m.Outpre = m.Outpre + "_beneficial36"
	}
	if ugly {
		m.Outpre = m.Outpre + "_ugly"
	}
	for _, rep := range reps {
		cfg := PlotCfg {
			Sbi: SyncBedInfo {
				Sync: SyncString("feral", bit),
				Bed: BedString(breed, "unbitted", rep),
				Info: InfoString("feral", bit),
				Out: "",
				Category: "",
			},
			ToUse: InfoSelection{[]string{fmt.Sprintf("%d", rep)}},
			BeneficialSbi: SyncBedInfo {
				Sync: SyncString(breed, "unbitted"),
				Bed: BedString(breed, "unbitted", rep),
				Info: InfoString(breed, "unbitted"),
				Out: "",
				Category: "",
			},
			BeneficialToUse: InfoSelection{[]string{fmt.Sprintf("%d", rep)}},
		}
		m.Cfgs = append(m.Cfgs, cfg)
	}
	return m
}

func BuildCfgAndControl(breed, bit string, reps []int, beneficial, beneficial36, ugly bool) []MultiPlotCfg {
	return []MultiPlotCfg {
		BuildCfg(breed, bit, reps, beneficial, beneficial36, ugly),
		BuildControl(breed, bit, reps, beneficial, beneficial36, ugly),
	}
}

func MakeCfg() {
	beneficialp := flag.Bool("b", false, "Plot beneficial alleles")
	beneficial36p := flag.Bool("B", false, "Plot beneficial alleles based on generation 36 pFST comparisons")
	uglyp := flag.Bool("u", false, "Make an ugly (lots of unnecessary lines) plot")
	flag.Parse()

	breeds := []string{"black", "white", "figurita", "runt"}
	bits := []string{"bitted", "unbitted"}
	reps := []int{1,2,3,4}
	cfgs := []MultiPlotCfg{}

	for _, breed := range breeds {
		for _, bit := range bits {
			cfgs = append(cfgs, BuildCfgAndControl(breed, bit, reps, *beneficialp, *beneficial36p, *uglyp)...)
		}
	}

	out, err := json.Marshal(cfgs)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))
}


// black_unbitted_pFst_subtractedalts_subfulls.bed_plfmt.bed

