# do-it

A minimal, fullscreen productivity timer for deep work.

Pure black background. Large white typography. Zero distractions.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Raylib](https://img.shields.io/badge/Raylib-5.5-red?style=flat)
![Platform](https://img.shields.io/badge/Platform-Windows-blue?style=flat)

---

## Features

- **Countdown Timer** — Set custom focus durations (25m, 50m, or custom)
- **Pomodoro Mode** — Full pomodoro cycle with focus/break transitions
- **Stopwatch Mode** — Track open-ended work sessions
- **Single Executable** — No DLLs, no installers, no dependencies
- **GPU Rendered** — Smooth 60 FPS via OpenGL (Raylib)
- **Embedded Fonts** — Inter typeface baked into the binary
- **Sound Notifications** — Procedurally generated bell/chime sounds
- **Configurable** — JSON config for durations, sounds, font scale

## Quick Start

```
do-it.exe
```

That's it. The application launches fullscreen with a 25-minute countdown.

## Keyboard Controls

| Key | Action |
|-------|-------------------------|
| Space | Start / Pause / Resume |
| R | Reset timer |
| F | Toggle fullscreen |
| P | Pomodoro mode |
| W | Stopwatch mode |
| S | Sound selector |
| 1 | 25 minute timer |
| 2 | 50 minute timer |
| 3 | 5 minute break |
| ESC | Exit |

## Build

**Requirements:** Go 1.21+ and GCC (MinGW-w64 on Windows)

```bash
# PowerShell
.\build.ps1

# Or Command Prompt
build.bat

# Or manually
set CGO_ENABLED=1
set GOARCH=amd64
go build -ldflags "-s -w -H windowsgui -extldflags '-static'" -o do-it.exe .
```

**Output:** `do-it.exe` (~5 MB, fully self-contained)

## Configuration

A `config.json` file is auto-created next to the executable on first run:

```json
{
  "focus_duration_minutes": 25,
  "short_break_minutes": 5,
  "long_break_minutes": 20,
  "sessions_before_long_break": 4,
  "default_timer_minutes": 25,
  "sound": "bell",
  "font_scale": 1.0,
  "task_name": "Deep Work",
  "volume": 0.4
}
```

## Project Structure

```
main.go        Entry point and window setup
app.go         Application state and lifecycle
timer.go       Time tracking (system clock based)
pomodoro.go    Pomodoro cycle management
renderer.go    All drawing and animation
input.go       Keyboard handling
sound.go       Procedural sound generation
fonts.go       Embedded font loading (Inter)
ui.go          Transient UI state
config.go      JSON configuration
assets/fonts/  Inter font files (embedded at compile time)
```

## License

MIT
