package main

import (
	"math"
	"time"
)

// TimerMode represents the current operating mode.
type TimerMode int

const (
	ModeCountdown TimerMode = iota
	ModePomodoro
	ModeStopwatch
)

// TimerState represents the current state of the timer.
type TimerState int

const (
	StateIdle TimerState = iota
	StateRunning
	StatePaused
	StateCompleted
)

// Timer handles time tracking independent of frame rate.
type Timer struct {
	Mode     TimerMode
	State    TimerState
	Duration time.Duration // Total duration for countdown/pomodoro
	Elapsed  time.Duration // Elapsed time

	startTime time.Time // When the timer was last started/resumed
	pausedAt  time.Duration // Elapsed time when paused

	// Completion callback
	OnComplete func()
}

func NewTimer() *Timer {
	return &Timer{
		Mode:  ModeCountdown,
		State: StateIdle,
	}
}

// SetDuration sets the countdown duration in minutes.
func (t *Timer) SetDuration(minutes int) {
	t.Duration = time.Duration(minutes) * time.Minute
	t.Reset()
}

// Start begins or resumes the timer.
func (t *Timer) Start() {
	if t.State == StateCompleted {
		t.Reset()
	}

	if t.State == StateIdle || t.State == StatePaused {
		t.startTime = time.Now()
		if t.State == StatePaused {
			// Resume from where we paused
			t.startTime = t.startTime.Add(-t.pausedAt)
		}
		t.State = StateRunning
	}
}

// Pause pauses the timer.
func (t *Timer) Pause() {
	if t.State == StateRunning {
		t.pausedAt = time.Since(t.startTime)
		t.State = StatePaused
	}
}

// Toggle starts or pauses depending on state.
func (t *Timer) Toggle() {
	switch t.State {
	case StateRunning:
		t.Pause()
	case StateIdle, StatePaused, StateCompleted:
		t.Start()
	}
}

// Reset resets the timer to initial state.
func (t *Timer) Reset() {
	t.State = StateIdle
	t.Elapsed = 0
	t.pausedAt = 0
}

// Update should be called every frame to update elapsed time.
func (t *Timer) Update() {
	if t.State != StateRunning {
		return
	}

	t.Elapsed = time.Since(t.startTime)

	// Check for completion in countdown/pomodoro modes
	if t.Mode != ModeStopwatch && t.Elapsed >= t.Duration {
		t.Elapsed = t.Duration
		t.State = StateCompleted
		if t.OnComplete != nil {
			t.OnComplete()
		}
	}
}

// Remaining returns the time left for countdown/pomodoro modes.
func (t *Timer) Remaining() time.Duration {
	if t.Mode == ModeStopwatch {
		return t.Elapsed
	}
	rem := t.Duration - t.Elapsed
	if rem < 0 {
		return 0
	}
	return rem
}

// Progress returns a value from 0.0 to 1.0 representing completion.
func (t *Timer) Progress() float64 {
	if t.Duration == 0 {
		return 0
	}
	if t.Mode == ModeStopwatch {
		return 0
	}
	p := float64(t.Elapsed) / float64(t.Duration)
	return math.Min(p, 1.0)
}

// DisplayMinutes returns the minutes component for display.
func (t *Timer) DisplayMinutes() int {
	d := t.Remaining()
	if t.Mode == ModeStopwatch {
		d = t.Elapsed
	}
	return int(d.Minutes()) % 60
}

// DisplaySeconds returns the seconds component for display.
func (t *Timer) DisplaySeconds() int {
	d := t.Remaining()
	if t.Mode == ModeStopwatch {
		d = t.Elapsed
	}
	return int(d.Seconds()) % 60
}

// DisplayString returns the formatted MM:SS string.
func (t *Timer) DisplayString() string {
	m := t.DisplayMinutes()
	s := t.DisplaySeconds()
	buf := [5]byte{
		byte('0' + m/10),
		byte('0' + m%10),
		':',
		byte('0' + s/10),
		byte('0' + s%10),
	}
	return string(buf[:])
}
