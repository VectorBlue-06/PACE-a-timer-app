@echo off
:: Build script for do-it timer
:: Produces a single static executable with no external DLL dependencies

set GOARCH=amd64
set CGO_ENABLED=1

echo Building do-it...
go build -ldflags "-s -w -H windowsgui -extldflags '-static'" -o do-it.exe .

if %ERRORLEVEL% EQU 0 (
    echo Build successful: do-it.exe
) else (
    echo Build failed.
    exit /b 1
)
