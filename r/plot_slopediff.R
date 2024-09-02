#!/usr/bin/env Rscript

#source("plot_pretty_multiple_helpers.R")
sourcedir = Sys.getenv("RLIBS")
source(paste(sourcedir, "/plot_pretty_multiple_helpers.R", sep=""))


main <- function() {
	args = commandArgs(trailingOnly=TRUE)
	full_slope_path = args[1]
	out_path = args[2]
	full_slope = read_slope(full_slope_path)
	plot_slope(full_slope, out_path, 20, 8, 300)
}

main()
