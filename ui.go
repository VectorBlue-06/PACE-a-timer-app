package main

// UI holds transient UI state for animations and visual feedback.
type UI struct {
	// Breathing animation for ambient dimming during focus
	breathPhase float64

	// Transition tracking
	lastMode  TimerMode
	lastState TimerState
}

func NewUI() *UI {
	return &UI{
		lastMode:  ModeCountdown,
		lastState: StateIdle,
	}
}

// Update advances UI animation state.
func (u *UI) Update(dt float32, mode TimerMode, state TimerState) {
	// Breathing animation during running state
	if state == StateRunning {
		u.breathPhase += float64(dt) * 0.5
	}

	u.lastMode = mode
	u.lastState = state
}
