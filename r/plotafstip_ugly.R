#!/usr/bin/env Rscript

library(ggplot2)
library(data.table)

main = function() {
	args = commandArgs(trailingOnly=TRUE)
	data = as.data.frame(fread(args[1], header=TRUE))
	smalldata = data[data$gen < 49,]

	p = ggplot(data, aes(gen, minor_f)) +
		geom_line(alpha=0.5, aes(groups = factor(chrposrepl))) +
		geom_smooth(data = smalldata, aes(groups = factor(chrposrepl)), alpha=0.25, color="blue", se = FALSE) +
		geom_smooth(data = smalldata, color="red", se = FALSE) +
		geom_smooth(data = smalldata, method="lm", color="dark green") +
		labs(title = "Allele frequencies over time", x = "Time (months)", y = "Minor allele frequency") +
		ylim(0, 1) +
		scale_color_discrete(name = "Replicate") +
		theme_bw()

	pdf(args[2], height = 3, width = 4)
		print(p)
	dev.off()
}

main()


# #!/usr/bin/env Rscript
# 
# suppressWarnings(library(ggplot2))
# 
# smallmtcars = mtcars[mtcars$disp <= 300,]
# 
# g = ggplot(data=mtcars, aes(disp, mpg))
# g2 = g + geom_line(aes(groups=factor(cyl)))
# g3 = g2 + geom_smooth(aes(groups=factor(cyl)), color="blue")
# g4 = g3 + geom_smooth(color="red")
# g5 = g2 + geom_smooth(data = smallmtcars, aes(groups=factor(cyl)), method='lm', color="blue")
# g6 = g5 + geom_smooth(data = smallmtcars, method='lm', color="red")
# print(g6)
# 
