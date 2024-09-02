#!/usr/bin/env Rscript

library(reshape2)
library(ggplot2)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=TRUE))
	data = data[data$Name != "bitted_unsel",]
	data$Treatment = factor(data$Treatment, levels=c("White", "Black", "Runt", "Figurita"))
	data$Name = factor(data$Name, levels=c("unbitted", "bitted", "unbitted_unsel"))

	breaks = c("unbitted", "bitted", "unbitted_unsel")
	labels = c("Unbitted", "Bitted", "Unselected")

	p = ggplot(data = data, aes(Time, FixedFreq, linetype = factor(Name))) +
		geom_line() +
		labs(x = "Time (months)", y = "Proportion fixed") +
		theme_bw() +
		scale_linetype_discrete(name="SNP class", breaks = breaks, labels = labels) +
		facet_grid(.~Treatment) +
		scale_x_continuous(breaks = seq(0, max(data$Time), 6))

	pdf(args[2], height = 3, width = 12)
		print(p)
	dev.off()

	if (length(args) < 3) {
		print("only 2 args")
		exit()
	}

	p = ggplot(data = data, aes(Time, LostFreq, linetype = factor(Name))) +
		geom_line() +
		labs(x = "Time (months)", y = "Proportion lost") +
		theme_bw() +
		scale_linetype_discrete(name="SNP class", breaks = breaks, labels = labels) +
		facet_grid(.~Treatment) +
		scale_x_continuous(breaks = seq(0, max(data$Time), 6))

	pdf(args[3], height = 3, width = 12)
		print(p)
	dev.off()

	if (length(args) < 4) {
		print("only 3 args")
		exit()
	}

	p = ggplot(data = data, aes(Time, FixedOrLostFreq, linetype = factor(Name))) +
		geom_line() +
		labs( x = "Time (months)", y = "Proportion fixed or lost") +
		theme_bw() +
		scale_linetype_discrete(name="SNP class", breaks = breaks, labels = labels) +
		facet_grid(.~Treatment) +
		scale_x_continuous(breaks = seq(0, max(data$Time), 6))

	pdf(args[4], height = 3, width = 12)
		print(p)
	dev.off()

	mdata = melt(data, na.rm = FALSE, id = c("Treatment", "Name", "Time"))
	smdata = mdata[mdata$variable == "FixedFreq" | mdata$variable == "LostFreq",]
	p = ggplot(data = smdata, aes(Time, value, linetype = factor(Name))) +
		geom_line() +
		labs( x = "Time (months)", y = "Proportion fixed or lost") +
		theme_bw() +
		scale_linetype_discrete(name="SNP class", breaks = breaks, labels = labels) +
		facet_grid(variable~Treatment) +
		scale_x_continuous(breaks = seq(0, max(data$Time), 6))

	pdf(args[5], height = 3, width = 12)
		print(p)
	dev.off()
}

main()
