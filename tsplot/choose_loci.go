package tsplot

import (
	"fmt"
	"math/rand"
	"golang.org/x/exp/slices"
	"github.com/jgbaldwinbrown/gobedtools"
	"io"
)

func ChooseSitesCore(count int, sites bedtools.Bedder, avoid bedtools.Bedder, r *rand.Rand) ([]bedtools.BedEntry, error) {
	bc, err := bedtools.IntersectBeds(sites, []string{"-v"}, avoid)
	if err != nil {
		return nil, err
	}
	rsites := bedtools.CollectBed(bedtools.BedChan{bc})
	r.Shuffle(len(rsites), func(i, j int) {rsites[i], rsites[j] = rsites[j], rsites[i]})
	return rsites[:count], nil
}

type BedReader struct {
	r io.Reader
}

func (b BedReader) Bed() (io.Reader, error) {
	return b.r, nil
}

func SyncEToBedReader(ss []SyncE) BedReader {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()

		for _, s := range ss {
			fmt.Fprintf(pw, "%v\t%v\t%v\n", s.Chr, s.Pos, s.Pos+1)
		}
	}()
	return BedReader{pr}
}

func BedEToBedReader(ss []BedE) BedReader {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()

		for _, s := range ss {
			fmt.Fprintf(pw, "%v\t%v\t%v\n", s.Chr, s.Start, s.End)
		}
	}()
	return BedReader{pr}
}

func ChooseSitesOld(count int, sites []SyncE, avoid []BedE, r *rand.Rand) ([]BedE, error) {
	sitebr := SyncEToBedReader(sites)
	avoidbr := BedEToBedReader(avoid)
	spans, err := ChooseSitesCore(count, sitebr, avoidbr, r)
	if err != nil {
		return nil, err
	}
	var out []BedE
	for _, span := range spans {
		out = append(out, BedE{
			span.Chr,
			span.Start,
			span.End,
			0,
		})
	}
	return out, nil
}

func ChooseSites(count int, sites []SyncE, r *rand.Rand) ([]SyncE, error) {
	rsites := slices.Clone(sites)
	r.Shuffle(len(rsites), func(i, j int) {rsites[i], rsites[j] = rsites[j], rsites[i]})
	if len(rsites) < count {
		return nil, fmt.Errorf("ChooseSites: rsites %v has len %v shorter than count %v", rsites, len(rsites), count)
	}
	return rsites[:count], nil
}

func ChooseSitesFull(count int, mcfg MultiPlotCfg, r *rand.Rand) ([]SyncE, error) {
	var fullsync []SyncE
	for _, pcfg := range mcfg.Cfgs {
		expsync, _, _, err := ReadSBIInverse(pcfg.BeneficialExpSbi, true)
		if err != nil {
			return nil, err
		}
		sites, err := ChooseSites(count, expsync, r)
		if err != nil {
			return nil, err
		}
		fullsync = append(fullsync, sites...)
	}
	return fullsync, nil
}

func CountSites(plottable [][]string, toUse InfoSelection) int {
	counter := map[string]struct{}{}
	for _, line := range plottable {
		for _, rep := range toUse.Reps {
			if line[8] == rep {
				counter[line[9]] = struct{}{}
			}
		}
	}
	return len(counter)
}

// func RunChooseSites() {
// 	seedp := flag.Int("s", 0, "Random seed")
// 	flag.Parse()
// 	r := rand.New(rand.NewSource(int64(*seedp)))
// 
// 	cfgs, err := ReadMultiPlotCfgs(os.Stdin)
// 	if err != nil {
// 		panic(err)
// 	}
// 
// 	for _, cfg := range cfgs {
// 		err := ChooseSitesFull(cfg, r)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}
// }
