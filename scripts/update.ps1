# for updating middlewares after each change.

go get github.com/72sevenzy2/Xrouter-middlewares@latest

if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to update Xrouter-middlewares"
    exit 1
}

go mod tidy

if ($LASTEXITCODE -ne 0) {
    Write-Host "go mod tidy failed"
    exit 1
}

Write-Host "Xrouter-middlewares updated successfully!"
