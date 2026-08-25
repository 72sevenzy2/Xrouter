#!/bin/bash

cd .. || exit 1

go test router_test.go -v
