#!/bin/bash
set -e

OUTPRE="${1}"

(cd cmd && (
	go build ts.go
	go build ts_separate.go
	go build resize_ranges.go
	go build multi_max_slope_categories.go
	go build multi_max_together.go
	go build make_multiplot_cfg.go
	go build make_slopebox_cfg.go
	go build slope_diffs.go
	go build slope_box.go
	go build gen0sfs.go
))

cp cmd/gen0sfs "${OUTPRE}gen0sfs"
cp cmd/ts "${OUTPRE}ts"
cp cmd/ts_separate "${OUTPRE}ts_separate"
cp cmd/resize_ranges "${OUTPRE}resize_ranges"
cp r/plotslopebox.R "${OUTPRE}plotslopebox"
cp r/plotslopeboxfinal.R "${OUTPRE}plotslopeboxfinal"
cp r/plotslopeboxfinal2.R "${OUTPRE}plotslopeboxfinal2"
cp r/plotslopeboxfinal3.R "${OUTPRE}plotslopeboxfinal3"
cp r/testslopeboxfinal.R "${OUTPRE}testslopeboxfinal"
cp r/plotfixationcurve.R "${OUTPRE}plotfixationcurve"
cp r/plotfixationcurvefinal.R "${OUTPRE}plotfixationcurvefinal"
cp r/plotfixationcurvefinal2.R "${OUTPRE}plotfixationcurvefinal2"
cp r/plotafs.R "${OUTPRE}plotafs"
cp r/plotafstip.R "${OUTPRE}plotafstip"
cp r/plotafstipfinal.R "${OUTPRE}plotafstipfinal"
cp r/plotgen0sfs.R "${OUTPRE}plotgen0sfs"
cp r/plotafstip_ugly.R "${OUTPRE}plotafstip_ugly"
cp cmd/multi_max_slope_categories ~/mybin
cp cmd/multi_max_together ~/mybin
cp cmd/make_multiplot_cfg ~/mybin
cp cmd/make_slopebox_cfg ~/mybin
cp cmd/slope_diffs ~/mybin
cp cmd/slope_box ~/mybin
