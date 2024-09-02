#!/usr/bin/env Rscript

library(ggplot2)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=TRUE))
	smalldata = data[data$gen < 49,]

	p = ggplot(smalldata, aes(gen, minor_f)) +
		geom_smooth(data = smalldata, aes(groups = factor(chrposrepl)), alpha=0.25, color="black", se = FALSE, size = 0.5) +
		geom_smooth(data = smalldata, color="red", se = FALSE, size = 1.5) +
		labs(x = "Time (months)", y = "Allele frequency") +
		ylim(0, 1) +
		scale_color_discrete(name = "Replicate") +
		scale_x_continuous(breaks = seq(0, max(smalldata$gen), 6)) +
		theme_bw() +
		theme(text = element_text(size=18))

	pdf(args[2], height = 3, width = 4)
		print(p)
	dev.off()
}

main()
