package tsplot

import (
	"os/exec"
	"os"
	"flag"
)

func PlotPlottableFileSfs(inpath, outpath string) error {
	cmd := exec.Command("plotgen0sfs", inpath, outpath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func PlotSfs(unsel bool, mcfg MultiPlotCfg) error {
	outpre := mcfg.Outpre
	if unsel {
		outpre = mcfg.Outpre + "_unsel"
	}
	inpath := outpre + plottablesuffix
	outpath := outpre + "_sfs.pdf"
	return PlotPlottableFileSfs(inpath, outpath)
}

func PlotEverySfs(unsel bool, mcfgs ...MultiPlotCfg) error {
	for _, mcfg := range mcfgs {
		err := PlotSfs(unsel, mcfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunPlotSfs() {
	unselptr := flag.Bool("P", false, "unselected sites")
	flag.Parse()

	mcfgs, err := ReadMultiPlotCfgs(os.Stdin)
	if err != nil {
		panic(err)
	}

	err = PlotEverySfs(*unselptr, mcfgs...)
	if err != nil {
		panic(err)
	}
}
