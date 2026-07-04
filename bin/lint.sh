#!/bin/bash
set -e
cd "$(dirname "$0")/.."
gofmt -w .
golangci-lint run ./...
