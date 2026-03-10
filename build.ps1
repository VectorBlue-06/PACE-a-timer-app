# Build script for do-it timer
# Produces a single static executable with no external DLL dependencies

$env:GOARCH = "amd64"
$env:CGO_ENABLED = "1"

Write-Host "Building do-it..." -ForegroundColor Cyan

go build -ldflags "-s -w -H windowsgui -extldflags '-static'" -o do-it.exe .

if ($LASTEXITCODE -eq 0) {
    $size = [math]::Round((Get-Item do-it.exe).Length / 1MB, 2)
    Write-Host "Build successful: do-it.exe ($size MB)" -ForegroundColor Green
} else {
    Write-Host "Build failed." -ForegroundColor Red
    exit 1
}
