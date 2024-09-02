#!/usr/bin/env Rscript

library(ggplot2)
library(data.table)

main = function() {
	breaks = c("unbitted", "bitted", "unbitted_unsel")
	labels = c("Unbitted", "Bitted", "Unselected")

	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=TRUE))
	data = data[data$snpclass != "bitted_unsel",]

	data$snpclass = factor(data$snpclass, levels=c("unbitted", "bitted", "unbitted_unsel"))

	p = ggplot(data = data[data$snpclass != "bitted_unsel",], aes(snpclass, Slope)) +
		geom_boxplot() +
		labs(x = "Treatment", y = "AF Slope") +
		lims(y=c(-0.012, 0.018)) +
		scale_x_discrete(name = "SNP class", breaks = breaks, labels = labels) +
		theme_bw() +
		theme(text = element_text(size=18))

	pdf(args[2], height = 3, width = 4)
		print(p)
	dev.off()
}

main()
