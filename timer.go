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
	Duration time.Duration
	Elapsed  time.Duration

	startTime time.Time
	pausedAt  time.Duration

	// Digit transition: always 1.0 for instant swap (no blink)
	digitAlpha float32

	OnComplete func()
}

func NewTimer() *Timer {
	return &Timer{
		Mode:       ModeCountdown,
		State:      StateIdle,
		digitAlpha: 1.0,
	}
}

func (t *Timer) SetDuration(minutes int) {
	t.Duration = time.Duration(minutes) * time.Minute
	t.Reset()
}

func (t *Timer) Start() {
	if t.State == StateCompleted {
		t.Reset()
	}
	if t.State == StateIdle || t.State == StatePaused {
		t.startTime = time.Now()
		if t.State == StatePaused {
			t.startTime = t.startTime.Add(-t.pausedAt)
		}
		t.State = StateRunning
	}
}

func (t *Timer) Pause() {
	if t.State == StateRunning {
		t.pausedAt = time.Since(t.startTime)
		t.State = StatePaused
	}
}

func (t *Timer) Toggle() {
	switch t.State {
	case StateRunning:
		t.Pause()
	case StateIdle, StatePaused, StateCompleted:
		t.Start()
	}
}

func (t *Timer) Reset() {
	t.State = StateIdle
	t.Elapsed = 0
	t.pausedAt = 0
	t.digitAlpha = 1.0
}

// Update should be called every frame.
func (t *Timer) Update(dt float32) {
	if t.State == StateRunning {
		t.Elapsed = time.Since(t.startTime)
		if t.Elapsed >= t.Duration {
			t.Elapsed = t.Duration
			t.State = StateCompleted
			if t.OnComplete != nil {
				t.OnComplete()
			}
		}
	}

	// Keep digit alpha at 1.0 — instant swap, no blink
	t.digitAlpha = 1.0
}

func (t *Timer) Remaining() time.Duration {
	rem := t.Duration - t.Elapsed
	if rem < 0 {
		return 0
	}
	return rem
}

func (t *Timer) Progress() float64 {
	if t.Duration == 0 {
		return 0
	}
	return math.Min(float64(t.Elapsed)/float64(t.Duration), 1.0)
}

func (t *Timer) DisplayMinutes() int {
	d := t.Remaining()
	return int(d.Minutes()) % 60
}

func (t *Timer) DisplaySeconds() int {
	d := t.Remaining()
	return int(d.Seconds()) % 60
}

// DisplayString returns MM:SS with zero-allocation formatting.
func (t *Timer) DisplayString() string {
	m := t.DisplayMinutes()
	s := t.DisplaySeconds()
	buf := [5]byte{
		byte('0' + m/10), byte('0' + m%10),
		':',
		byte('0' + s/10), byte('0' + s%10),
	}
	return string(buf[:])
}

// DigitAlpha returns the current digit transition alpha (0..1).
func (t *Timer) DigitAlpha() float32 {
	return t.digitAlpha
}
