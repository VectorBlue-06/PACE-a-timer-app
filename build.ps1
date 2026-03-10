# Build script for PACE timer
# Produces a single static executable with no external DLL dependencies

$env:GOARCH = "amd64"
$env:CGO_ENABLED = "1"

Write-Host "Building PACE..." -ForegroundColor Cyan

go build -ldflags "-s -w -H windowsgui -extldflags '-static'" -o pace.exe .

if ($LASTEXITCODE -eq 0) {
    $size = [math]::Round((Get-Item pace.exe).Length / 1MB, 2)
    Write-Host "Build successful: pace.exe ($size MB)" -ForegroundColor Green
} else {
    Write-Host "Build failed." -ForegroundColor Red
    exit 1
}
