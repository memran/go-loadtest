# PowerShell build script for Windows
# Creates binaries for Windows, macOS, and Linux

param(
    [string]$Version = "v1.0.0"
)

$BUILD_DIR = "build"

Write-Host "Building go-loadtest version $Version..." -ForegroundColor Green

# Create build directory
if (-not (Test-Path $BUILD_DIR)) {
    New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null
}

# Build for Linux (64-bit)
Write-Host "Building for Linux (amd64)..." -ForegroundColor Yellow
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags "-X main.Version=$Version" -o "$BUILD_DIR/go-loadtest-linux-amd64" main.go

# Build for Linux (ARM64)
Write-Host "Building for Linux (arm64)..." -ForegroundColor Yellow
$env:GOOS = "linux"
$env:GOARCH = "arm64"
go build -ldflags "-X main.Version=$Version" -o "$BUILD_DIR/go-loadtest-linux-arm64" main.go

# Build for macOS (Intel)
Write-Host "Building for macOS (amd64)..." -ForegroundColor Yellow
$env:GOOS = "darwin"
$env:GOARCH = "amd64"
go build -ldflags "-X main.Version=$Version" -o "$BUILD_DIR/go-loadtest-darwin-amd64" main.go

# Build for macOS (Apple Silicon)
Write-Host "Building for macOS (arm64)..." -ForegroundColor Yellow
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -ldflags "-X main.Version=$Version" -o "$BUILD_DIR/go-loadtest-darwin-arm64" main.go

# Build for Windows (64-bit)
Write-Host "Building for Windows (amd64)..." -ForegroundColor Yellow
$env:GOOS = "windows"
$env:GOARCH = "amd64"
go build -ldflags "-X main.Version=$Version" -o "$BUILD_DIR/go-loadtest-windows-amd64.exe" main.go

# Build for Windows (32-bit)
Write-Host "Building for Windows (386)..." -ForegroundColor Yellow
$env:GOOS = "windows"
$env:GOARCH = "386"
go build -ldflags "-X main.Version=$Version" -o "$BUILD_DIR/go-loadtest-windows-386.exe" main.go

# Reset environment variables
$env:GOOS = ""
$env:GOARCH = ""

Write-Host ""
Write-Host "Build complete! Binaries are in the $BUILD_DIR directory:" -ForegroundColor Green
Get-ChildItem $BUILD_DIR | Format-Table Name, Length

Write-Host ""
Write-Host "To create a release:" -ForegroundColor Cyan
Write-Host "  1. Create a git tag: git tag $Version"
Write-Host "  2. Push the tag: git push origin $Version"
Write-Host "  3. Create a GitHub release and upload the binaries from $BUILD_DIR/"
