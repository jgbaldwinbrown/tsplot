#!/bin/bash
set -e

(cd cmd && (
	go build ts.go
	go build ts_separate.go
	go build resize_ranges.go
	go build multi_max_slope_categories.go
	go build multi_max_together.go
	go build make_multiplot_cfg.go
))

cp cmd/ts ~/mybin/ts
cp cmd/ts_separate ~/mybin/ts_separate
cp cmd/resize_ranges ~/mybin/resize_ranges
cp r/plotafs.R ~/mybin/plotafs
cp r/plotafstip.R ~/mybin/plotafstip
cp r/plotafstip_ugly.R ~/mybin/plotafstip_ugly
cp cmd/multi_max_slope_categories ~/mybin
cp cmd/multi_max_together ~/mybin
cp cmd/make_multiplot_cfg ~/mybin
