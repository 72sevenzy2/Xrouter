#!/bin/bash

cd .. || exit 1
go test -bench="^(BenchmarkDynamicRoutes|BenchmarkStdlibDynamic)$" -run=^$ -count=10 > results.txt
benchstat results.txt
