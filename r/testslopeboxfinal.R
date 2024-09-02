#!/usr/bin/env Rscript

library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	
	data = as.data.frame(fread(args[1], header=TRUE))
	data = data[data$snpclass != "bitted_unsel",]

	data$snpclass = factor(data$snpclass, levels=c("unbitted", "bitted", "unbitted_unsel"))

	model = aov(Slope~snpclass, data=data)
	print(summary(model))
	print(TukeyHSD(model))
	plot(TukeyHSD(model, conf.level=.95), las=2)
}

main()
