#!/bin/bash

echo "Tidying modules..."
go mod tidy

echo "Building Go binary..."
go build -o app

if [ $? -eq 0 ]; then
    echo "✅ Build completed successfully"
else
    echo "❌ Build failed"
    exit 1
fi
