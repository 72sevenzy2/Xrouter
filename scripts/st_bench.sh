#!/bin/bash

cd .. || exit 1
cd benchmarks || exit 1
go test -bench="^(BenchmarkStaticRoutes|BenchmarkStdlibStatic)$" -run=^$ -count=10 > staticBenchmarks.txt
benchstat staticBenchmarks.txt
