#!/bin/bash

cd .. || exit 1
cd benchmarks || exit 1
go test -bench="^(BenchmarkXrouter|BenchmarkGolangStdlibRouter)$" -run=^$ -cpuprofile=1cpu.prof -memprofile=1mem.prof -count=10 > dynamicBenchmarks.txt
benchstat dynamicBenchmarks.txt
