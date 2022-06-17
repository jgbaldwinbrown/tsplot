#!/bin/bash
set -e

(cd scripts && go build ts_separate.go)
cp r/plotafs.R .

./scripts/ts_separate -s testdata/minisync.sync -b testdata/toplot_sep.bed -i testdata/black_pooled_info_ne_bitted.txt > plottables_separate.txt

rm plotafs.R
# ./r/plotafs.R plottable.txt plotted.pdf
# evince plotted.pdf
