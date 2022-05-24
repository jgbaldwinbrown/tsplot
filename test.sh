#!/bin/bash
set -e

(cd scripts && go build ts.go)
./scripts/ts -s testdata/minisync.sync -b testdata/toplot.bed -i testdata/black_pooled_info_ne_bitted.txt > plottable.txt
./r/plotafs.R plottable.txt plotted.pdf
evince plotted.pdf
