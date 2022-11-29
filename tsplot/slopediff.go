package tsplot

import (
	// "os"
	// "github.com/jgbaldwinbrown/pmap/pkg"
	// "golang.org/x/exp/slices"
)

// type SlopeDiffOpts struct {
// 	Threads int
// }
// 
// type SlopeDiffCfg struct {
// 	PolarizerPaths []string
// 	SlopePaths []string
// 	OutPre string
// }
// 
// type SlopeDiffArgs struct {
// 	Cfg SlopeDiffCfg
// 	Opts SlopeDiffOpts
// 	Polarizations map[string][]int
// }
// 
// // func UniquePolarizerPathSets(cfgs ...SlopeDiffCfg) map[string][]string {
// // 	upaths := make(map[string][]string)
// // 	for _, cfg := range cfgs {
// // 		paths := slices.Clone(cfg.PolarizerPaths)
// // 		
// // 	}
// // }
// 
// // func Polarize(path string) []int, error {
// // }
// 
// func PolarizeAll(cfgs []SlopeDiffCfg, threads int) {
// 	topolarize := UniquePolarizerPaths(cfgs...)
// 	// MapE(function, cfgs, threads)
// }
// 
// func RunSlopeDiff() {
// 	opts := GetSlopeDiffOpts()
// 	cfgs := ReadSlopeDiffCfgs(os.Stdin)
// 
// 	polarizations, err := PolarizeAll(cfgs, opts.Threads)
// 	if err != nil {
// 		panic(err)
// 	}
// 
// 	args := CollectSlopeDiffArgs(cfgs, opts, polarizations)
// 
// 	err := AllSlopeDiffs(args, opts.Threads)
// 	if err != nil {
// 		panic(err)
// 	}
// 
// 	err := PlotAllSlopeDiffs(args, opts.Threads)
// 	if err != nil {
// 		panic(err)
// 	}
// }
