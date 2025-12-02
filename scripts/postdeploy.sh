#!/bin/bash
set -e

echo "Building Go binary for linux/amd64..."
cd src
GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/handler .

echo "Getting function app name..."
FUNCTION_APP_NAME=$(azd env get-value SERVICE_PROCESSOR_NAME)

echo "Publishing to Azure Functions app: $FUNCTION_APP_NAME"
func azure functionapp publish "$FUNCTION_APP_NAME" --custom

echo "Deployment complete!"
