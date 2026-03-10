@echo off
:: Build script for PACE timer
:: Produces a single static executable with no external DLL dependencies

set GOARCH=amd64
set CGO_ENABLED=1

echo Building PACE...
go build -ldflags "-s -w -H windowsgui -extldflags '-static'" -o pace.exe .

if %ERRORLEVEL% EQU 0 (
    echo Build successful: pace.exe
) else (
    echo Build failed.
    exit /b 1
)
