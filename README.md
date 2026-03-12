<p align="center">
  <img src="./assets/PACE-banner.png" width="100%" alt="PACE banner">
</p>

<h1 align="center">PACE</h1>

<p align="center">
  A minimal productivity timer built for deep focus.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/platform-Windows-blue" alt="Platform">
  <img src="https://img.shields.io/badge/language-Go-00ADD8" alt="Go">
  <img src="https://img.shields.io/badge/renderer-Raylib%205.5-red" alt="Raylib">
  <img src="https://img.shields.io/badge/binary-~5MB-green" alt="Size">
  <img src="https://img.shields.io/badge/dependencies-zero-brightgreen" alt="Dependencies">
</p>



## Features

- **Single executable** — ~5 MB, no dependencies
- **Pomodoro-first** — Starts in Pomodoro mode by default
- **Two modes** — Pomodoro and Countdown timer
- **Customizable** — Remap keys and load your own alarm sound
- **Smooth UI** — 60 FPS animations and settings panel

---

## Screenshots

<p align="center">
  <img src="./assets/screenshots/home1.png" width="45%" alt="Home screen - Timer running">
  <img src="./assets/screenshots/home2.png" width="45%" alt="Home screen - Timer stopped">
</p>

<p align="center">
  <img src="./assets/screenshots/pomodoro-stop.png" width="45%" alt="Pomodoro mode - Timer stopped">
  <img src="./assets/screenshots/setting.png" width="45%" alt="Settings panel">
</p>



---

**Download for Windows:** [Releases](https://github.com/VectorBlue-06/PACE-a-timer-app/releases)

> Currently only Windows is supported. macOS and Linux coming soon!

---


---




## Installation

### Prerequisites

- [Go](https://go.dev/dl/) 1.21 or later
- [GCC](https://www.mingw-w64.org/) (MinGW-w64 on Windows)

### Build

```bash
# PowerShell
.\build.ps1

# Command Prompt
build.bat

# Manual
set CGO_ENABLED=1
set GOARCH=amd64
go build -ldflags "-s -w -H windowsgui -extldflags '-static'" -o pace.exe .
```

**Output:** `pace.exe` (~5 MB, fully self-contained)

---

## How to Use

1. Download and run `pace.exe`
2. Timer starts in **Pomodoro mode** by default
3. Press **Space** to start the timer
4. Press **1**, **2**, or **3** to switch between countdown presets (25 min, 50 min, 5 min break)
5. Press **P** to switch to Pomodoro mode (focus → break → focus cycle)
6. Press **TAB** to open settings and customize durations, sounds, and keybinds
7. Hold **ESC** for 1 second to quit

> NOTE - Controls may change in future for the ease of use.

---

## Usage

Run `PACE.exe`. The app starts in Pomodoro mode.

### Keyboard Shortcuts

| Key           | Action                          |
|---------------|---------------------------------|
| Space         | Start / Pause / Resume          |
| R             | Reset timer                     |
| M             | Minimize window                 |
| TAB           | Open / close settings panel     |
| Ctrl+Space    | Toggle frame persistence mode   |
| P             | Pomodoro mode                   |
| S             | Sound selector                  |
| 1             | 25 minute timer                 |
| 2             | 50 minute timer                 |
| 3             | 5 minute break                  |
| Hold ESC (1s) | Quit app                        |

### Settings Panel (TAB)

| Key   | Action            |
|-------|-------------------|
| ↑ ↓   | Navigate options  |
| ← →   | Adjust values     |
| ENTER | Browse alarm file |
| TAB   | Close settings    |

### Quit Behavior

- Press and hold **ESC** for 1 second to quit.
- While ESC is held, a smooth top-left status text (`quiting...`) appears.

### Config Location

- Config is stored in system app data: `%AppData%/PACE/config.json`.
- On first launch after update, legacy config next to the executable is read and migrated.

---

## Project Structure

```
do-it/
├── assets/
│   ├── fonts/              # Inter TTF files (embedded at compile time)
│   ├── sounds/             # WAV sound files (embedded at compile time)
│   ├── PACE-banner.png
│   └── PACE-logo.png
├── main.go                 # Window init and render loop
├── app.go                  # Central state machine and lifecycle
├── timer.go                # System-clock timer with digit transitions
├── pomodoro.go             # Pomodoro cycle state machine
├── renderer.go             # Layered rendering via Raylib
├── input.go                # Configurable keyboard event routing
├── sound.go                # Embedded sound loading and playback
├── fonts.go                # Embedded TTF loading via go:embed
├── ui.go                   # Animation engine (blink, scale, fade)
├── config.go               # JSON configuration with key bindings
├── dialog_windows.go       # Windows file picker for custom alarm
├── build.ps1               # PowerShell build script
├── build.bat               # Command Prompt build script
├── (AppData)/PACE/config.json  # Auto-generated user config
└── docs/
    └── DOCUMENTATION.md    # Full technical documentation
```

---

## Documentation

See [docs/DOCUMENTATION.md](DOCUMENTATION.md) for the full technical reference — architecture, rendering pipeline, animation system, configuration schema, and more.

---

## License

This project is provided as-is for personal use.
