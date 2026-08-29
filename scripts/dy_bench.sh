#!/bin/bash

cd .. || exit 1
go test -bench="^(BenchmarkXrouter|BenchmarkGolangStdlibRouter)$" -run=^$ -count=10 > results.txt
benchstat results.txt
