#!/bin/bash

# for updating middlewares after each change.

go get github.com/72sevenzy2/Xrouter-middlewares@latest

if [$? -ne 0]; then 
	   echo "failed to update."
	   exit 1
fi

go mod tidy

if [$? -ne 0]; then 
	   echo "go mod tidy failed."
	   exit 1
fi

echo "successfully updated."
