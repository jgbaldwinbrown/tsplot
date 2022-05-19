#!/bin/bash
set -e

go build ts.go
./ts -s minisync.sync -b toplot.bed -i black_pooled_info_ne_bitted.txt > plottable.txt
./plotafs.R plottable.txt plotted.pdf
evince plotted.pdf
