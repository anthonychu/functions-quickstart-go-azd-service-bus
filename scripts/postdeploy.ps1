$ErrorActionPreference = "Stop"

Write-Host "Building Go binary for linux/amd64..."
Push-Location src
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o bin/linux-amd64/handler .
if ($LASTEXITCODE -ne 0) {
    Pop-Location
    throw "Go build failed"
}

Write-Host "Getting function app name..."
$FUNCTION_APP_NAME = azd env get-value SERVICE_PROCESSOR_NAME
if ($LASTEXITCODE -ne 0) {
    Pop-Location
    throw "Failed to get function app name"
}

Write-Host "Publishing to Azure Functions app: $FUNCTION_APP_NAME"
func azure functionapp publish $FUNCTION_APP_NAME --custom
if ($LASTEXITCODE -ne 0) {
    Pop-Location
    throw "Deployment failed"
}

Pop-Location
Write-Host "Deployment complete!"
