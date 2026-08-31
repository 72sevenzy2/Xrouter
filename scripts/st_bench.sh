#!/bin/bash

cd .. || exit 1
cd benchmarks || exit 1
go test -bench="^(BenchmarkXrouterStaticRoutes|BenchmarkGolangStdlibStaticRoutes)$" -run=^$ -cpuprofile=2cpu.prof -memprofile=2mem.prof -count=10 > staticBenchmarks.txt
benchstat staticBenchmarks.txt
