#!/usr/bin/env Rscript

library(ggplot2)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=TRUE))
	smalldata = data[data$gen < 49,]

	p = ggplot(smalldata, aes(Treatment, Slope)) +
		geom_boxplot(data = smalldata,
		labs(title = "Allele frequency slopes of related treatments", x = "Treatment", y = "Allele frequency Slope") +
		theme_bw()

	pdf(args[2], height = 3, width = 4)
		print(p)
	dev.off()
}

main()
