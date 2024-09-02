#!/usr/bin/env Rscript

library(ggplot2)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=TRUE))

	p = ggplot(data = data, aes(Time, FixedFreq, color = factor(Name))) +
		geom_line() +
		labs(title = "Frequency of allele fixation", x = "Time (months)", y = "Proportion fixed") +
		theme_bw() +
		scale_color_discrete(name="Treatment") +
		scale_x_continuous(breaks = seq(0, max(data$Time), 6))

	pdf(args[2], height = 3, width = 4)
		print(p)
	dev.off()

	if (length(args) < 3) {
		print("only 2 args")
		exit()
	}

	p = ggplot(data = data, aes(Time, LostFreq, color = factor(Name))) +
		geom_line() +
		labs(title = "Frequency of allele loss", x = "Time (months)", y = "Proportion lost") +
		theme_bw() +
		scale_color_discrete(name="Treatment") +
		scale_x_continuous(breaks = seq(0, max(data$Time), 6))

	pdf(args[3], height = 3, width = 4)
		print(p)
	dev.off()

	if (length(args) < 4) {
		print("only 3 args")
		exit()
	}

	p = ggplot(data = data, aes(Time, FixedOrLostFreq, color = factor(Name))) +
		geom_line() +
		labs(title = "Frequency of allele fixation or loss", x = "Time (months)", y = "Proportion fixed or lost") +
		theme_bw() +
		scale_color_discrete(name="Treatment") +
		scale_x_continuous(breaks = seq(0, max(data$Time), 6))

	pdf(args[4], height = 3, width = 4)
		print(p)
	dev.off()
}

main()
