#!/bin/bash

cd .. || exit 1
go test middleware_test.go -v
