#!/bin/bash
set -e

(cd scripts && (
	go build ts.go
	go build ts_separate.go
	go build resize_ranges.go
))

cp scripts/ts ~/mybin/ts
cp scripts/ts_separate ~/mybin/ts_separate
cp scripts/resize_ranges ~/mybin/resize_ranges
cp r/plotafs.R ~/mybin/plotafs
