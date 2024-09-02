#!/usr/bin/env Rscript

library(ggplot2)
library(data.table)

main = function() {
	breaks = c("unbitted", "bitted", "unbitted_unsel")
	labels = c("Selected\nby preening", "Same alleles,\nimpaired preening", "Random alleles,\nnormal preening")

	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=TRUE))
	data = data[data$snpclass != "bitted_unsel",]

	data$snpclass = factor(data$snpclass, levels=c("unbitted", "bitted", "unbitted_unsel"))

	annotation = data.frame(annot = c("a", "b", "b"), Slope = c(0.017, 0.017, 0.017), snpclass = c("unbitted", "bitted", "unbitted_unsel"))

	toplot = data[data$snpclass != "bitted_unsel",]
	p = ggplot(data = toplot, aes(snpclass, Slope)) +
		geom_boxplot() +
		labs(x = "Treatment", y = "AF Slope") +
		lims(y=c(-0.012, 0.018)) +
		scale_x_discrete(name = "SNP class", breaks = breaks, labels = labels) +
		geom_text(data = annotation, aes(snpclass, Slope, label = annot)) +
		theme_bw() +
		theme(text = element_text(size=18), axis.text.x = element_text(size=8)) +

	pdf(args[2], height = 3, width = 4)
		print(p)
	dev.off()
}

main()

		# annotate("a", x = 0, y = 0.016) +
		# annotate("b", x = 1, y = 0.016) +
		# annotate("b", x = 2, y = 0.016)
