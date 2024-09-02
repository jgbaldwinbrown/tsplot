#!/bin/bash
set -e

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

cp cmd/gen0sfs ~/mybin/gen0sfs
cp cmd/ts ~/mybin/ts
cp cmd/ts_separate ~/mybin/ts_separate
cp cmd/resize_ranges ~/mybin/resize_ranges
cp r/plotslopebox.R ~/mybin/plotslopebox
cp r/plotslopeboxfinal.R ~/mybin/plotslopeboxfinal
cp r/plotslopeboxfinal2.R ~/mybin/plotslopeboxfinal2
cp r/plotslopeboxfinal3.R ~/mybin/plotslopeboxfinal3
cp r/testslopeboxfinal.R ~/mybin/testslopeboxfinal
cp r/plotfixationcurve.R ~/mybin/plotfixationcurve
cp r/plotfixationcurvefinal.R ~/mybin/plotfixationcurvefinal
cp r/plotfixationcurvefinal2.R ~/mybin/plotfixationcurvefinal2
cp r/plotafs.R ~/mybin/plotafs
cp r/plotafstip.R ~/mybin/plotafstip
cp r/plotafstipfinal.R ~/mybin/plotafstipfinal
cp r/plotgen0sfs.R ~/mybin/plotgen0sfs
cp r/plotafstip_ugly.R ~/mybin/plotafstip_ugly
cp cmd/multi_max_slope_categories ~/mybin
cp cmd/multi_max_together ~/mybin
cp cmd/make_multiplot_cfg ~/mybin
cp cmd/make_slopebox_cfg ~/mybin
cp cmd/slope_diffs ~/mybin
cp cmd/slope_box ~/mybin
