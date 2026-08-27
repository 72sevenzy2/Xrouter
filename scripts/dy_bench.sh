#!/bin/bash

cd .. || exit 1
go test -bench="^(BenchmarkDynamicRoutes|BenchmarkStdlibDynamic)$" -run=^$ -count=10 -benchmem > results.txt
benchstat results.txt
