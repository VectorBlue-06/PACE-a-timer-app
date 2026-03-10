# Build script for PACE timer
# Produces a single static executable with no external DLL dependencies






























































































































MIT## License```assets/fonts/  Inter font files (embedded at compile time)config.go      JSON configuration with key bindingsui.go          Animation state machine (blink, scale, fade)fonts.go       Embedded font loading (Inter)sound.go       Procedural sound generationinput.go       Configurable keyboard handlingrenderer.go    Layered drawing, stable text, overlayspomodoro.go    Pomodoro cycle managementtimer.go       Time tracking with digit transitionsapp.go         Application state and lifecyclemain.go        Entry point and window setup```## Project StructureSettings can also be adjusted in real-time via the TAB settings panel.```}  "keys": { ... }  "frame_persistence": false,  "enable_animations": true,  "show_progress_ring": true,  "volume": 0.4,  "task_name": "Deep Work",  "font_scale": 1.0,  "sound": "bell",  "default_timer_minutes": 25,  "sessions_before_long_break": 4,  "long_break_minutes": 20,  "short_break_minutes": 5,  "focus_duration_minutes": 25,{```jsonA `config.json` file is auto-created next to the executable on first run:## Configuration**Output:** `pace.exe` (~5 MB, fully self-contained)```go build -ldflags "-s -w -H windowsgui -extldflags '-static'" -o pace.exe .set GOARCH=amd64set CGO_ENABLED=1# Manualbuild.bat# Command Prompt.\build.ps1# PowerShell```bash**Requirements:** Go 1.21+ and GCC (MinGW-w64 on Windows)## Build| TAB | Close settings || ← → | Adjust values || ↑ ↓ | Navigate options ||-------|-------------------|| Key | Action |### Settings Panel Navigation| ESC | Exit || 3 | 5 minute break || 2 | 50 minute timer || 1 | 25 minute timer || S | Sound selector || W | Stopwatch mode || P | Pomodoro mode || Ctrl+Space | Toggle frame persistence mode || TAB | Open / close settings panel || F | Toggle fullscreen || R | Reset timer || Space | Start / Pause / Resume ||-------------|-------------------------------|| Key | Action |All keys are configurable via `config.json`.## Keyboard ControlsLaunches fullscreen with a 25-minute countdown. Press Space to start.```pace.exe```## Quick Start- **Procedural Sound** — Bell and chime generated at runtime, no WAV files shipped- **Embedded Fonts** — Inter typeface baked into the binary- **Configurable Keys** — All keyboard shortcuts stored in config- **Settings Panel** — TAB opens a full settings overlay- **Frame Persistence** — Optional trailing-fade rendering mode (Ctrl+Space)- **Progress Ring** — Thin circular indicator behind the timer- **Digit Transitions** — 120ms fade when digits change- **Scale Animation** — Subtle 5% scale-up when timer starts- **Blink Feedback** — Paused (1s cycle) and completed (0.5s cycle) blink patterns- **Stable Typography** — Fixed-width digit cells prevent horizontal jitter- **GPU Rendered** — Smooth 60 FPS via OpenGL (Raylib)- **Single Executable** — No DLLs, no installers, zero runtime dependencies- **Stopwatch Mode** — Open-ended session tracking- **Pomodoro Mode** — Full cycle with focus, short break, long break- **Countdown Timer** — 25m, 50m, or custom focus durations## Features---![Platform](https://img.shields.io/badge/Platform-Windows-blue?style=flat)![Raylib](https://img.shields.io/badge/Raylib-5.5-red?style=flat)![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)Pure black background. Large white typography. Zero distractions.A minimal, fullscreen productivity timer for deep work.
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
