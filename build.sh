#!/bin/bash

# Build script for cross-platform releases
# Creates binaries for Windows, macOS, and Linux

VERSION=${1:-"v1.0.0"}
BUILD_DIR="build"

echo "Building go-loadtest version $VERSION..."

# Create build directory
mkdir -p $BUILD_DIR

# Build for Linux (64-bit)
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=$VERSION" -o $BUILD_DIR/go-loadtest-linux-amd64 main.go

# Build for Linux (ARM64)
echo "Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=$VERSION" -o $BUILD_DIR/go-loadtest-linux-arm64 main.go

# Build for macOS (Intel)
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=$VERSION" -o $BUILD_DIR/go-loadtest-darwin-amd64 main.go

# Build for macOS (Apple Silicon)
echo "Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.Version=$VERSION" -o $BUILD_DIR/go-loadtest-darwin-arm64 main.go

# Build for Windows (64-bit)
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=$VERSION" -o $BUILD_DIR/go-loadtest-windows-amd64.exe main.go

# Build for Windows (32-bit)
echo "Building for Windows (386)..."
GOOS=windows GOARCH=386 go build -ldflags "-X main.Version=$VERSION" -o $BUILD_DIR/go-loadtest-windows-386.exe main.go

echo ""
echo "Build complete! Binaries are in the $BUILD_DIR directory:"
ls -lh $BUILD_DIR/
echo ""
echo "To create a release:"
echo "  1. Create a git tag: git tag $VERSION"
echo "  2. Push the tag: git push origin $VERSION"
echo "  3. Create a GitHub release and upload the binaries from $BUILD_DIR/"
