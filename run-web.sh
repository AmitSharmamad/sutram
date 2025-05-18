#!/bin/bash

set -e # Exit on error
set -x # Echo commands

echo "Building the web application..."
mkdir -p bin

# Generate the wire_gen.go file if wire tool is available
if command -v wire &> /dev/null; then
  echo "Generating wire_gen.go..."
  wire ./cmd/web
fi

# Build the application
go build -v -o bin/web-server cmd/web/main.go

echo "Running the web server..."
./bin/web-server
