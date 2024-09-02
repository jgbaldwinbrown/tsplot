#!/usr/bin/env Rscript

library(ggplot2)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=TRUE))
	smalldata = data[data$gen == 0,]
	smalldata$true_minor_f = sapply(smalldata$minor_f, function(x){if (x>0.5) {return(1-x)}; return(x)})

	p = ggplot(smalldata, aes(true_minor_f), color="black") +
		geom_histogram(binwidth=0.1, aes(y=stat(width*density))) +
		labs(y = "Density", x = "Allele frequency") +
		xlim(-0.1, 0.6) +
		ylim(-0.1, 1.1) +
		theme_bw() +
		theme(text = element_text(size=18))

	pdf(args[2], height = 3, width = 4)
		print(p)
	dev.off()
}
		# xlim(0, 0.5) +
		# ylim(0, 1) +
		# geom_histogram(binwidth=0.1, aes(y=..density..)) +
		# geom_histogram(binwidth=0.1) +
		# xlim(-0.1, 0.6) +
		# ylim(-0.1, 1.1) +

main()
