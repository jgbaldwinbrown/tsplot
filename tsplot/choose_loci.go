package tsplot

import (
	"fmt"
	"rand"
	"golang.org/x/exp/slices"
	"github.com/jgbaldwinbrown/gobedtools"
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

func (b BedReader) Bed() io.Reader {
	return b.r
}

func ChooseSites(sites []SyncE, avoid []BedE, r *rand.Rand) ([]BedE, error) {
	
}
