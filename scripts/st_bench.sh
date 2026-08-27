#!/bin/bash

cd .. || exit 1
go test -bench="^(BenchmarkStaticRoutes|BenchmarkStdlibStatic)" -count=10 > results.txt
benchstat results.txt
