#!/usr/bin/env Rscript

library(ggplot2)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=TRUE))

	p = ggplot(data = data, aes(Treatment, Slope)) +
		geom_boxplot() +
		labs(title = "Allele frequency slopes of related treatments", x = "Treatment", y = "Allele frequency Slope") +
		theme_bw()

	pdf(args[2], height = 3, width = 4)
		print(p)
	dev.off()
}

main()
