#!/bin/bash
set -e

echo "Installing SCC with fixed dependencies..."

# Create a temporary directory
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR"

# Create a simple module to fix dependency issues
cat > go.mod << EOF
module github.com/cloudartisan/install-scc

go 1.20

// Force correct versions
replace golang.org/x/text => golang.org/x/text v0.3.8

require github.com/boyter/scc v1.12.1
EOF

# Install a specific working version
go install github.com/boyter/scc@v1.12.1

echo "SCC has been installed successfully. You can now run 'scc' from your terminal."