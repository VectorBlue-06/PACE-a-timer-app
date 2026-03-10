# do-it — Technical Documentation

## Overview

**do-it** is a minimal fullscreen productivity timer built with Go and Raylib. It compiles to a single static Windows executable (~5 MB) with no runtime dependencies. All assets — fonts and sounds — are embedded directly in the binary.

The design philosophy prioritizes calm, focused aesthetics: pure black background, crisp white Inter typography, and smooth GPU-accelerated rendering at 60 FPS.

---

## Architecture

### Module Layout

| File | Responsibility |
|----------------|-----------------------------------------------|
| `main.go` | Window initialization, render loop |
| `app.go` | Central state machine, lifecycle management |
| `timer.go` | System-clock-based time tracking |
| `pomodoro.go` | Pomodoro cycle logic (focus/break/long break) |
| `renderer.go` | All OpenGL drawing via Raylib |
| `input.go` | Keyboard event routing |
| `sound.go` | Procedural WAV generation, audio playback |
| `fonts.go` | Embedded TTF loading via `//go:embed` |
| `ui.go` | Transient animation state |
| `config.go` | JSON configuration load/save |

### Data Flow

```
main.go
  └─ creates App
       ├─ Timer      (system time based)
       ├─ Pomodoro   (cycle state machine)
       ├─ Renderer   (draws everything)
       ├─ SoundSystem (generates and plays sounds)
       └─ UI         (animation state)

Per frame:
  1. Timer.Update()     — read system clock
  2. UI.Update()        — advance animations
  3. Renderer.Update()  — update fade/pulse state
  4. HandleInput()      — process keyboard
  5. Renderer.DrawFrame() — render everything
```

---

## Timer System

### Clock Independence

The timer does **not** count frames. It uses `time.Now()` and `time.Since()` from Go's standard library, which provides nanosecond-resolution wall-clock time on Windows.

```
Start:  startTime = time.Now()
Update: elapsed = time.Since(startTime)
Pause:  pausedAt = elapsed; save startTime offset
Resume: startTime = time.Now() - pausedAt
```

This ensures the timer remains accurate regardless of frame drops, system load, or vsync behavior.

### Display Format

The timer displays `MM:SS` with zero-padded digits. The display string is computed without `fmt.Sprintf` for zero allocations:

```go
buf := [5]byte{
    byte('0' + m/10), byte('0' + m%10),
    ':',
    byte('0' + s/10), byte('0' + s%10),
}
```

### Timer Modes

| Mode | Behavior |
|-----------|----------------------------------------------|
| Countdown | Counts down from a set duration |
| Pomodoro | Countdown with automatic phase transitions |
| Stopwatch | Counts up indefinitely |

---

## Pomodoro System

### Default Cycle

```
[Focus 25m] → [Short Break 5m] → [Focus 25m] → [Short Break 5m] →
[Focus 25m] → [Short Break 5m] → [Focus 25m] → [Long Break 20m] → repeat
```

The default is 4 focus sessions before a long break. All durations are configurable via `config.json`.

### State Machine

```
PhaseFocus ──complete──► PhaseShortBreak (if session < total)
PhaseFocus ──complete──► PhaseLongBreak  (if session >= total)
PhaseShortBreak ──complete──► PhaseFocus (session++)
PhaseLongBreak  ──complete──► PhaseFocus (session = 1)
```

Transitions are triggered by the user pressing Space after timer completion (not automatic), giving the user control over when to start the next phase.

### Session Tracking

The pomodoro system tracks:
- Current session number (1-based)
- Total completed focus sessions today
- Total focus minutes today

This information is displayed subtly at the bottom of the screen in Pomodoro mode.

---

## Rendering System

### Font Rendering

Fonts are loaded at high resolution (200px base size × font scale) and scaled down for display. This produces crisp text at any screen resolution.

The Inter typeface is used in three weights:
- **Inter-Bold** — Timer digits (large, centered)
- **Inter-SemiBold** — Mode indicators (subtle, top of screen)
- **Inter-Regular** — Labels, hints, session info

Font files are embedded at compile time using Go's `//go:embed` directive and loaded from memory via `rl.LoadFontFromMemory`. No font files needed at runtime.

### Typography Scale

All font sizes are proportional to screen height:

| Element | Size | Weight |
|----------------|----------------------|----------|
| Timer digits | 14% of screen height | Bold |
| Task label | 2.5% of screen height| Regular |
| Session info | 1.8% of screen height| Regular |
| Keyboard hints | 1.4% of screen height| Regular |
| Mode indicator | 1.3% of screen height| SemiBold |

### Color System

| Name | Hex | Usage |
|-----------|---------|-------------------------------|
| Background| #000000 | Screen background |
| Primary | #FFFFFF | Timer digits |
| Secondary | #8A8A8A | Labels, session info |
| Subtle | #333333 | Hints, ring background |
| Accent | #50C878 | Progress ring (focus) |
| Break | #64A0FF | Progress ring (break) |
| Complete | #FFC850 | Progress ring (completed) |

### Progress Ring

A thin circular progress ring (3px) surrounds the timer display. It animates smoothly as time progresses, drawn from the top (12 o'clock position) clockwise.

The ring changes color based on state:
- **Green** — Focus session active
- **Blue** — Break period
- **Gold** — Timer completed

The ring is rendered as line segments along an arc, producing smooth curves at any resolution.

### Animations

| Animation | Duration | Trigger |
|--------------------|----------|-------------------------|
| Startup fade-in | ~333ms | Application launch |
| State change fade | ~200ms | Start/pause/complete |
| Completion pulse | Continuous| Timer reaches zero |

All animations use actual frame time (`GetFrameTime()`) rather than fixed increments, ensuring consistent behavior regardless of frame rate.

### Vertical Positioning

The timer sits at **45%** of screen height (slightly above center), following macOS-style layout conventions where primary content sits above the optical center. This creates a calm, balanced visual composition.

---

## Sound System

### Procedural Generation

Sounds are generated mathematically at runtime — no WAV files are embedded. This keeps the binary small.

#### Bell Sound
Multi-harmonic sine wave with exponential decay:
```
f₁ = 880 Hz  (amplitude 0.6)
f₂ = 1320 Hz (amplitude 0.3)
f₃ = 1760 Hz (amplitude 0.1)
envelope = e^(-3t)
duration = 1.5 seconds
```

#### Chime Sound
Two-tone sequence with overlapping decay:
```
Tone 1: 523 Hz (C5) + 785 Hz harmonic, starts at t=0
Tone 2: 659 Hz (E5) + 988 Hz harmonic, starts at t=0.3s
envelope = e^(-2.5t) per tone
duration = 1.8 seconds
```

### Audio Pipeline

1. WAV data is generated in memory as `[]byte`
2. Written to a temporary file (OS temp directory)
3. Loaded by Raylib's audio system
4. Temp file immediately deleted

This is necessary because Raylib's `LoadSound` requires a file path. The temp files exist only momentarily during initialization.

### Sound Options

| Sound | Description |
|-------|--------------------------|
| bell | Warm multi-harmonic bell |
| chime | Two-tone ascending chime |
| none | Silent |

Sound selection is accessible via the `S` key overlay menu, and persisted to `config.json`.

---

## Input System

### Keyboard-First Design

All interaction is keyboard-driven. There is no mouse cursor, no buttons, no click targets. This reduces visual noise and keeps the interface focused.

### Input Routing

The sound selection overlay captures input when visible, preventing timer controls from firing while choosing sounds. ESC is handled globally for application exit.

### Key Bindings

```
Space     Toggle timer (start/pause/resume/advance)
R         Reset timer to initial state
F         Toggle borderless fullscreen
P         Switch to Pomodoro mode
W         Switch to Stopwatch mode
S         Open/close sound selector
1         Set 25-minute countdown (when idle)
2         Set 50-minute countdown (when idle)
3         Set 5-minute break countdown (when idle)
ESC       Exit application
```

Number keys (1/2/3) only work when the timer is idle or paused, and not in Pomodoro mode (where the cycle manages durations automatically).

---

## Configuration

### File Location

`config.json` is created next to the executable on first run. If the file is missing or malformed, defaults are used.

### Options

| Field | Type | Default | Description |
|----------------------------|---------|---------|-------------------------------|
| `focus_duration_minutes` | int | 25 | Pomodoro focus session length |
| `short_break_minutes` | int | 5 | Short break length |
| `long_break_minutes` | int | 20 | Long break length |
| `sessions_before_long_break`| int | 4 | Sessions before long break |
| `default_timer_minutes` | int | 25 | Default countdown on startup |
| `sound` | string | "bell" | "bell", "chime", or "none" |
| `font_scale` | float | 1.0 | Global font size multiplier |
| `task_name` | string | "Deep Work"| Label shown above timer |
| `volume` | float | 0.4 | Sound volume (0.0 – 1.0) |

### Validation

Values are clamped to safe ranges:
- `font_scale`: 0.5 – 3.0
- `volume`: 0.0 – 1.0
- All durations: minimum 1 minute

---

## Build System

### Static Compilation

The build uses CGO to compile Raylib's C source code directly into the Go binary:

```
CGO_ENABLED=1 → raylib C source compiled via GCC
-ldflags "-s -w" → strip debug symbols
-H windowsgui → no console window
-extldflags '-static' → static linking (no DLL dependencies)
```

### Dependencies at Build Time

| Tool | Version | Purpose |
|--------|---------|------------------------------|
| Go | 1.21+ | Compiler and build system |
| GCC | MinGW-w64 | C compiler for CGO (raylib) |

### Dependencies at Runtime

**None.** The executable is fully self-contained.

### Build Output

| Metric | Value |
|--------|-------|
| Binary size | ~5 MB |
| External DLLs | 0 |
| Config files | 1 (auto-created) |

---

## Performance

### Design Targets

| Metric | Target |
|------------|----------------|
| Frame rate | 60 FPS (vsync) |
| CPU usage | < 2% |
| Memory | < 30 MB |
| Startup | < 200ms |

### Rendering Efficiency

- VSync enabled (no busy-loop rendering)
- MSAA 4x for smooth lines
- Bilinear texture filtering for crisp font rendering
- No dynamic allocations in the render loop
- Timer string computed without `fmt.Sprintf`

### Font Loading

Fonts are loaded once at startup at high resolution. The GPU handles downscaling via bilinear filtering, ensuring crisp text at any display size without runtime resampling.

---

## Design Decisions

### Why Raylib?

- Direct OpenGL access for smooth rendering
- Minimal dependency footprint
- Compiles from source via CGO (no shared libraries)
- Built-in font rendering with TTF support
- Built-in audio with WAV support
- Well-maintained, stable API

### Why Embedded Assets?

Embedding fonts via `//go:embed` and generating sounds procedurally means the entire application is a single file. No installer, no asset folder, no PATH configuration. Copy the exe anywhere and run it.

### Why System Time?

Frame-counting timers drift. A 60 FPS timer loses ~1 second every 10 minutes due to frame timing variance. Using `time.Now()` provides nanosecond accuracy from the OS clock, independent of rendering performance.

### Why Borderless Fullscreen?

Borderless fullscreen provides the immersive, distraction-free experience without the display mode switching delay of true fullscreen. The window covers the entire screen but can be alt-tabbed smoothly.
