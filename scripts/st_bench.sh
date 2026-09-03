#!/bin/bash

cd .. || exit 1
cd benchmarks || exit 1
go test -bench="^(BenchmarkXrouterStaticRoutes|BenchmarkGolangStdlibStaticRoutes)$" -run=^$ -benchmem -cpuprofile=2cpu.prof -memprofile=2mem.prof -count=10 > staticBenchmarks.txt
benchstat staticBenchmarks.txt

go tool pprof 2mem.prof
