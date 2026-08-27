#!/bin/bash

cd .. || exit 1
go test -bench="^(BenchmarkStaticRoutes|BenchmarkStdlibStatic)" -benchmem -count=10 > results.txt
benchstat results.txt
