#!/bin/bash
# for updating middlewares after each change.

run() {
    "$@"
    if [ $? -ne 0 ]; then
        echo "failed: $*"
        exit 1
    fi
}

run go get -u github.com/72sevenzy2/Xrouter-middlewares@latest
run go mod tidy
echo "successfully updated middlewares."

run go get -u github.com/72sevenzy2/json-parser@latest
run go mod tidy
echo "successfully updated json-parser."
