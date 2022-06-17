#!/bin/bash
set -e

(cd scripts && (
	go build ts.go
	go build ts_separate.go
))

cp scripts/ts ~/mybin/ts
cp scripts/ts_separate ~/mybin/ts_separate
cp r/plotafs.R ~/mybin/plotafs
