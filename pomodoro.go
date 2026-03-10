package main

import "fmt"

// PomodoroPhase represents the current phase of a pomodoro cycle.
type PomodoroPhase int

const (
	PhaseFocus PomodoroPhase = iota
	PhaseShortBreak
	PhaseLongBreak
)

// Pomodoro manages the pomodoro work/break cycle.
type Pomodoro struct {
	Phase           PomodoroPhase
	CurrentSession  int // 1-based session number
	TotalSessions   int
	FocusMinutes    int
	ShortBreakMin   int
	LongBreakMin    int
	CompletedToday  int
	TotalFocusToday int // total focus minutes completed today

	timer *Timer
}

func NewPomodoro(cfg AppConfig, timer *Timer) *Pomodoro {
	return &Pomodoro{
		Phase:          PhaseFocus,
		CurrentSession: 1,
		TotalSessions:  cfg.SessionsBeforeLong,
		FocusMinutes:   cfg.FocusDuration,
		ShortBreakMin:  cfg.ShortBreak,
		LongBreakMin:   cfg.LongBreak,
		timer:          timer,
	}
}

// Setup configures the timer for the current pomodoro phase.
func (p *Pomodoro) Setup() {
	p.timer.Mode = ModePomodoro
	switch p.Phase {
	case PhaseFocus:
		p.timer.SetDuration(p.FocusMinutes)
	case PhaseShortBreak:
		p.timer.SetDuration(p.ShortBreakMin)
	case PhaseLongBreak:
		p.timer.SetDuration(p.LongBreakMin)
	}
}

// Advance moves to the next phase in the pomodoro cycle.
func (p *Pomodoro) Advance() {
	switch p.Phase {
	case PhaseFocus:
		p.CompletedToday++
		p.TotalFocusToday += p.FocusMinutes
		if p.CurrentSession >= p.TotalSessions {
			p.Phase = PhaseLongBreak
		} else {
			p.Phase = PhaseShortBreak
		}
	case PhaseShortBreak:
		p.CurrentSession++
		p.Phase = PhaseFocus
	case PhaseLongBreak:
		p.CurrentSession = 1
		p.Phase = PhaseFocus
	}
	p.Setup()
}

// PhaseLabel returns a human-readable label for the current phase.
func (p *Pomodoro) PhaseLabel() string {
	switch p.Phase {
	case PhaseFocus:
		return "Focus"
	case PhaseShortBreak:
		return "Short Break"
	case PhaseLongBreak:
		return "Long Break"
	}
	return ""
}

// SessionLabel returns "Session X of Y".
func (p *Pomodoro) SessionLabel() string {
	return fmt.Sprintf("Session %d of %d", p.CurrentSession, p.TotalSessions)
}

// TodayLabel returns total focus time today.
func (p *Pomodoro) TodayLabel() string {
	h := p.TotalFocusToday / 60
	m := p.TotalFocusToday % 60
	if h > 0 {
		return fmt.Sprintf("Today: %dh %dm", h, m)
	}
	return fmt.Sprintf("Today: %dm", m)
}
