# ====================================
# Go Cross Compile Build Script
# Windows -> macOS, Linux, Windows
# ====================================

$ErrorActionPreference = "Stop"

$targets = @(
    @{GOOS="darwin"; GOARCH="amd64"; OUT="dist/capella-erp6-be-darwin-amd64"}   # macOS Intel
    @{GOOS="darwin"; GOARCH="arm64"; OUT="dist/capella-erp6-be-darwin-arm64"}   # macOS ARM (M1/M2/M3)
    @{GOOS="linux";  GOARCH="amd64"; OUT="dist/capella-erp6-be-linux-amd64"}    # Linux Intel
    @{GOOS="linux";  GOARCH="arm64"; OUT="dist/capella-erp6-be-linux-arm64"}    # Linux ARM
    @{GOOS="windows"; GOARCH="amd64"; OUT="dist/capella-erp6-be-windows-amd64.exe"} # Windows
)

foreach ($t in $targets) {
    Write-Host "==========================================" -ForegroundColor Cyan
    Write-Host "Building $($t.OUT) for $($t.GOOS)/$($t.GOARCH)..."

    $env:GOOS        = $t.GOOS
    $env:GOARCH      = $t.GOARCH
    $env:CGO_ENABLED = 0

    go build -trimpath -ldflags="-s -w" -o $t.OUT .

    Write-Host "Success: $($t.OUT)" -ForegroundColor Green
}

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "ALL BUILD DONE!"
