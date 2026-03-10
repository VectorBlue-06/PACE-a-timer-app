# PACE — Technical Documentation

## Overview

**PACE** is a minimal fullscreen productivity timer built with Go and Raylib. It compiles to a single static Windows executable (~5 MB) with no runtime dependencies. All assets — fonts and sounds — are embedded directly in the binary.

The design philosophy prioritizes calm, focused aesthetics: pure black background, crisp white Inter typography, and smooth GPU-accelerated rendering at 60 FPS.

---

## Architecture

### Module Layout

| File | Responsibility |
|----------------|--------------------------------------------------|
| `main.go` | Window initialization, render loop |
| `app.go` | Central state machine, lifecycle management |
| `timer.go` | System-clock-based time tracking, digit transitions |
| `pomodoro.go` | Pomodoro cycle logic (focus/break/long break) |
| `renderer.go` | Layered OpenGL drawing via Raylib |
| `input.go` | Configurable keyboard event routing |
| `sound.go` | Procedural WAV generation, audio playback |
| `fonts.go` | Embedded TTF loading via `//go:embed` |
| `ui.go` | Animation state machine (blink, scale, fade) |
| `config.go` | JSON configuration with key bindings |

### Data Flow

```
main.go
  └─ creates App
       ├─ Timer      (system time based, tracks digit transitions)
       ├─ Pomodoro   (cycle state machine)
       ├─ Renderer   (layered drawing)
       ├─ SoundSystem (generates and plays sounds)
       └─ UI         (animation state: blink, scale, fade, settings)

Per frame:
  1. Timer.Update(dt)     — read system clock, advance digit fade
  2. UI.Update(dt, ...)   — advance blink, scale, fade, settings animations
  3. HandleInput()        — process keyboard with configurable bindings
  4. Renderer.DrawFrame() — render all layers in correct order
```

---

## Rendering System

### Layer Order

The renderer draws in strict order to prevent the progress ring from overlapping text:

1. **Background** — Full clear or frame persistence fade
2. **Progress ring** — Thin circular arc behind everything
3. **Mode indicator** — Subtle label at top
4. **Label text** — Task name / phase above timer
5. **Timer text** — Large MM:SS digits, centered
6. **Session info** — Session count and today's focus
7. **Keyboard hints** — Subtle hints at bottom
8. **Overlays** — Settings panel, sound selector (always on top)

### Stable Text Rendering

Timer digits are rendered in **fixed-width character cells** to prevent horizontal jitter when values change (e.g., "25:00" → "24:59").

Implementation:
1. At startup, measure the widest digit (0–9) in the Bold font
2. Each digit is drawn centered within a cell of that maximum width
3. The colon has its own fixed cell width
4. Total width = 4 × digitWidth + 1 × colonWidth + gaps
5. The entire block is centered on screen

This ensures the timer appears visually locked to the center regardless of which digits are displayed.

### Frame Persistence Mode

An optional rendering mode where the screen is not fully cleared each frame. Instead, a semi-transparent black rectangle is drawn over the previous frame:

```go
rl.DrawRectangle(0, 0, screenW, screenH, Color(0, 0, 0, 25))
```

This creates soft motion trails for a dreamlike visual effect. Toggled via Ctrl+Space or the settings panel.

---

## Animation System

All animations use real time deltas from `time.Now()`, not frame counting.

### Timer Scale

When the timer starts running, it smoothly scales from 1.0× to 1.05× over ~200ms, creating subtle emphasis. When paused, it returns to 1.0×.

### Blink Patterns

| State | Visible | Invisible | Cycle |
|-----------|---------|-----------|-------|
| Paused | 1.0s | 1.0s | 2.0s |
| Completed | 0.5s | 0.5s | 1.0s |

Blink is computed via `math.Mod(elapsed, cycle)` using wall-clock time, making it frame-rate independent.

### Digit Transitions

When the displayed MM:SS string changes, the new digits fade in over ~120ms (alpha 0 → 1). This is tracked per-frame in `timer.go` using a `digitAlpha` field.

### State Change Fade

When the timer changes state (idle → running, running → paused, etc.), a brief opacity dip to 70% occurs, recovering to 100% over ~200ms.

### Startup Fade

The entire UI fades in from black over ~333ms on application launch.

### Completion Pulse

When the timer reaches zero, a continuous sinusoidal pulse modulates the timer text alpha between 60% and 100%.

### Settings Panel

The settings overlay fades in/out over ~180ms with a smooth alpha transition.

### Animation Timing Summary

| Animation | Duration |
|---------------------|----------|
| Timer scale | ~200ms |
| Digit fade | ~120ms |
| State change fade | ~200ms |
| Startup fade-in | ~333ms |
| Settings panel fade | ~180ms |
| Pause blink cycle | 2.0s |
| Complete blink cycle | 1.0s |

---

## Timer System

### Clock Independence

The timer uses `time.Now()` and `time.Since()` for nanosecond-resolution wall-clock timing. It never counts frames.

```
Start:  startTime = time.Now()
Update: elapsed = time.Since(startTime)
Pause:  pausedAt = elapsed
Resume: startTime = time.Now() - pausedAt
```

### Timer Modes

| Mode | Behavior |
|-----------|----------------------------------------------|
| Countdown | Counts down from a set duration |
| Pomodoro | Countdown with phase-based cycle management |
| Stopwatch | Counts up indefinitely |

---

## Pomodoro System

### Default Cycle

```
[Focus 25m] → [Short Break 5m] → [Focus 25m] → [Short Break 5m] →
[Focus 25m] → [Short Break 5m] → [Focus 25m] → [Long Break 20m] → repeat
```

### State Machine

```
PhaseFocus ──complete──► PhaseShortBreak (if session < total)
PhaseFocus ──complete──► PhaseLongBreak  (if session >= total)
PhaseShortBreak ──complete──► PhaseFocus (session++)
PhaseLongBreak  ──complete──► PhaseFocus (session = 1)
```

Transitions are user-initiated (press Space after completion).

---

## Input System

### Configurable Key Bindings

All keyboard shortcuts are stored in `config.json` under the `keys` object. Each binding specifies:

- `key` — Raylib key code (int32)
- `ctrl` — Whether Ctrl must be held (boolean)
- `name` — Display name shown in UI hints

### Default Bindings

| Action | Key | Code |
|---------------------|-------------|------|
| Start / Pause | Space | 32 |
| Reset | R | 82 |
| Fullscreen | F | 70 |
| Settings | Tab | 258 |
| Frame Persistence | Ctrl+Space | 32 |
| Pomodoro Mode | P | 80 |
| Stopwatch Mode | W | 87 |
| Sound Selector | S | 83 |
| Preset 25m | 1 | 49 |
| Preset 50m | 2 | 50 |
| Preset Break | 3 | 51 |
| Quit | Escape | 256 |

### Input Priority

1. Settings panel (when open) captures all input
2. Sound selector (when open) captures all input
3. Normal input routing

### Settings Panel Controls

When the settings panel is open:
- **↑↓** navigate between options
- **←→** adjust the selected value
- **TAB** or **ESC** closes the panel

---

## Settings Panel

Pressing TAB opens a full-screen settings overlay with these options:

| Setting | Adjustable | Step |
|---------------------|------------|------|
| Focus Duration | 1–120 min | ±5 |
| Short Break | 1–30 min | ±1 |
| Long Break | 1–60 min | ±5 |
| Sessions | 1–10 | ±1 |
| Sound | bell/chime/none | cycle |
| Volume | 0–100% | ±10% |
| Font Scale | 0.5–3.0x | ±0.1 |
| Progress Ring | On/Off | toggle |
| Animations | On/Off | toggle |
| Frame Persistence | On/Off | toggle |

Changes are saved to `config.json` immediately.

---

## Color System

| Name | Hex | Usage |
|-----------|---------|-------------------------------|
| Background| #000000 | Screen background |
| Primary | #FFFFFF | Timer digits, active labels |
| Secondary | #8A8A8A | Labels, session info |
| Subtle | #333333 | Hints, ring background |
| Accent | #50C878 | Progress ring (focus) |
| Break | #64A0FF | Progress ring (break) |
| Complete | #FFC850 | Progress ring (completed) |

---

## Typography

### Font Weights

| Weight | Usage |
|------------|-------------------------------------|
| Inter-Bold | Timer digits |
| Inter-SemiBold | Mode indicator, settings titles |
| Inter-Regular | Labels, hints, session info |

### Sizing (relative to screen height)

| Element | Size |
|----------------|------|
| Timer digits | 14% |
| Task label | 2.5% |
| Session info | 1.6% |
| Keyboard hints | 1.3% |
| Mode indicator | 1.2% |

### Positioning

The timer sits at **45%** of screen height (slightly above vertical center), following macOS-style layout conventions.

---

## Sound System

### Procedural Generation

Sounds are generated mathematically at runtime — no WAV files are shipped.

**Bell:** Multi-harmonic sine wave (880/1320/1760 Hz) with exponential decay, 1.5s duration.

**Chime:** Two-tone sequence (C5 at 523 Hz, E5 at 659 Hz) with overlapping decay, 1.8s duration.

### Audio Pipeline

1. WAV data generated in memory
2. Written to OS temp directory
3. Loaded by Raylib audio system
4. Temp file immediately deleted

---

## Progress Ring

A thin (2.5px) circular progress ring drawn behind the timer text at 18% of screen height radius.

| State | Color |
|-----------|---------|
| Focus | #50C878 |
| Break | #64A0FF |
| Completed | #FFC850 |

The ring animates smoothly from 12 o'clock clockwise. A background ring at 40% opacity provides visual context.

---

## Configuration

### File Location

`config.json` is created next to the executable on first run.

### Full Schema

| Field | Type | Default | Range |
|----------------------------|---------|---------|-------------|
| `focus_duration_minutes` | int | 25 | 1–120 |
| `short_break_minutes` | int | 5 | 1–30 |
| `long_break_minutes` | int | 20 | 1–60 |
| `sessions_before_long_break`| int | 4 | 1–10 |
| `default_timer_minutes` | int | 25 | 0+ |
| `sound` | string | "bell" | bell/chime/none |
| `font_scale` | float | 1.0 | 0.5–3.0 |
| `task_name` | string | "Deep Work" | — |
| `volume` | float | 0.4 | 0.0–1.0 |
| `show_progress_ring` | bool | true | — |
| `enable_animations` | bool | true | — |
| `frame_persistence` | bool | false | — |
| `keys` | object | (defaults) | — |

---

## Build System

### Static Compilation

```
CGO_ENABLED=1 → raylib C source compiled via GCC
-ldflags "-s -w" → strip debug symbols
-H windowsgui → no console window
-extldflags '-static' → static linking (no DLL dependencies)
```

### Dependencies

**Build time:** Go 1.21+, GCC (MinGW-w64)
**Runtime:** None

### Output

| Metric | Value |
|--------|-------|
| Binary size | ~5 MB |
| External DLLs | 0 |
| Config files | 1 (auto-created) |

---

## Performance

| Metric | Target |
|------------|----------------|
| Frame rate | 60 FPS (vsync) |
| CPU usage | < 2% |
| Memory | < 30 MB |
| Startup | < 200ms |

- VSync enabled (no busy-loop rendering)
- MSAA 4x for smooth lines
- Bilinear texture filtering for crisp fonts
- No dynamic allocations in the render loop
- Timer string computed without fmt.Sprintf
- All animations use real time deltas
