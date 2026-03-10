package main

import (
	"math"
	"time"
)

// UI holds all transient animation and visual feedback state.
type UI struct {
	// Timer scale animation (1.0 normal, 1.05 when running)
	timerScale       float32
	timerScaleTarget float32

	// Blink state (time-based, not frame-based)
	blinkVisible bool
	blinkTimer   float64

	// Startup fade
	fadeAlpha float32

	// State change transition
	stateAlpha float32

	// Completion pulse
	pulsePhase float64

	// Settings panel
	settingsAlpha       float32
	settingsAlphaTarget float32

	// Frame persistence fade
	framePersistence bool

	// Breathing animation
	breathPhase float64

	// Wall clock for time-based blink
	lastUpdate time.Time

	prevState TimerState
}

func NewUI() *UI {
	return &UI{
		timerScale:       1.0,
		timerScaleTarget: 1.0,
		blinkVisible:     true,
		stateAlpha:       1.0,
		fadeAlpha:        0.0,
		lastUpdate:       time.Now(),
	}
}

// Update advances all animation state using real time deltas.
func (u *UI) Update(dt float32, state TimerState, enableAnim bool) {
	now := time.Now()
	realDt := float32(now.Sub(u.lastUpdate).Seconds())
	u.lastUpdate = now
	if realDt > 0.1 {
		realDt = 0.1
	}

	// Startup fade-in (~333ms)
	if u.fadeAlpha < 1.0 {
		u.fadeAlpha += realDt * 3.0
		if u.fadeAlpha > 1.0 {
			u.fadeAlpha = 1.0
		}
	}

	// State change transition (~200ms)
	if state != u.prevState {
		u.stateAlpha = 0.7
		u.prevState = state
	}
	if u.stateAlpha < 1.0 {
		u.stateAlpha += realDt * 5.0
		if u.stateAlpha > 1.0 {
			u.stateAlpha = 1.0
		}
	}

	// Timer scale: 1.05 when running, 1.0 otherwise (~200ms transition)
	if enableAnim {
		if state == StateRunning {
			u.timerScaleTarget = 1.05
		} else {
			u.timerScaleTarget = 1.0
		}
		// Smooth lerp toward target
		diff := u.timerScaleTarget - u.timerScale
		u.timerScale += diff * realDt * 5.0 // ~200ms
		if math.Abs(float64(diff)) < 0.001 {
			u.timerScale = u.timerScaleTarget
		}
	} else {
		u.timerScale = 1.0
	}

	// Blink logic (time-based)
	u.updateBlink(state, realDt)

	// Completion pulse
	if state == StateCompleted {
		u.pulsePhase += float64(realDt) * 3.0
	} else {
		u.pulsePhase = 0
	}

	// Breathing during focus
	if state == StateRunning {
		u.breathPhase += float64(realDt) * 0.5
	}

	// Settings panel fade (~180ms)
	diff := u.settingsAlphaTarget - u.settingsAlpha
	u.settingsAlpha += diff * realDt * 5.56 // 1/0.18
	if math.Abs(float64(diff)) < 0.01 {
		u.settingsAlpha = u.settingsAlphaTarget
	}
}

func (u *UI) updateBlink(state TimerState, dt float32) {
	switch state {
	case StatePaused:
		// 1s visible, 1s invisible (2s cycle)
		u.blinkTimer += float64(dt)
		cycle := math.Mod(u.blinkTimer, 2.0)
		u.blinkVisible = cycle < 1.0

	case StateCompleted:
		// 0.5s visible, 0.5s invisible (1s cycle)
		u.blinkTimer += float64(dt)
		cycle := math.Mod(u.blinkTimer, 1.0)
		u.blinkVisible = cycle < 0.5

	default:
		u.blinkTimer = 0
		u.blinkVisible = true
	}
}

// EffectiveAlpha returns the combined alpha for the main content.
func (u *UI) EffectiveAlpha() float32 {
	a := u.fadeAlpha * u.stateAlpha
	if a < 0 {
		return 0
	}
	return a
}

// TimerVisible returns whether the timer text should be visible (for blink).
func (u *UI) TimerVisible() bool {
	return u.blinkVisible
}

// TimerScale returns the current scale multiplier for the timer text.
func (u *UI) TimerScale() float32 {
	return u.timerScale
}

// PulseAlpha returns a 0..1 pulsing value for completion state.
func (u *UI) PulseAlpha() float32 {
	return float32(0.6 + 0.4*math.Sin(u.pulsePhase))
}

// ShowSettings toggles the settings overlay visibility.
func (u *UI) ShowSettings(show bool) {
	if show {
		u.settingsAlphaTarget = 1.0
	} else {
		u.settingsAlphaTarget = 0.0
	}
}

// SettingsVisible returns true if the settings overlay has any visibility.
func (u *UI) SettingsVisible() bool {
	return u.settingsAlpha > 0.01
}

// SettingsAlpha returns the current settings panel opacity.
func (u *UI) SettingsAlpha() float32 {
	return u.settingsAlpha
}

// SettingsFullyOpen returns true if settings is fully visible.
func (u *UI) SettingsFullyOpen() bool {
	return u.settingsAlphaTarget > 0.5 && u.settingsAlpha > 0.95
}
