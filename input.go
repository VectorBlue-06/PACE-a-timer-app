package main

import rl "github.com/gen2brain/raylib-go/raylib"

// HandleInput processes all keyboard input and returns true if the app should exit.
func HandleInput(app *App) bool {
	// Sound menu takes priority
	if app.ShowSoundMenu {
		return handleSoundMenu(app)
	}

	// ESC — exit
	if rl.IsKeyPressed(rl.KeyEscape) {
		return true
	}

	// Space — toggle timer
	if rl.IsKeyPressed(rl.KeySpace) {
		if app.Timer.State == StateCompleted && app.Timer.Mode == ModePomodoro {
			app.Pomodoro.Advance()
			app.Timer.Start()
		} else {
			app.Timer.Toggle()
		}
	}

	// R — reset
	if rl.IsKeyPressed(rl.KeyR) {
		app.Timer.Reset()
		if app.Timer.Mode == ModePomodoro {
			app.Pomodoro.Phase = PhaseFocus
			app.Pomodoro.CurrentSession = 1
			app.Pomodoro.Setup()
		}
	}

	// F — toggle fullscreen
	if rl.IsKeyPressed(rl.KeyF) {
		rl.ToggleBorderlessWindowed()
	}

	// P — pomodoro mode
	if rl.IsKeyPressed(rl.KeyP) {
		app.Timer.Mode = ModePomodoro
		app.Timer.Reset()
		app.Pomodoro = NewPomodoro(app.Config, app.Timer)
		app.Pomodoro.Setup()
	}

	// W — stopwatch mode
	if rl.IsKeyPressed(rl.KeyW) {
		app.Timer.Mode = ModeStopwatch
		app.Timer.Reset()
	}

	// S — sound selector
	if rl.IsKeyPressed(rl.KeyS) {
		app.ShowSoundMenu = true
	}

	// Number presets (only in countdown mode or idle)
	if app.Timer.State == StateIdle || app.Timer.State == StatePaused {
		if app.Timer.Mode != ModePomodoro {
			if rl.IsKeyPressed(rl.KeyOne) {
				app.Timer.Mode = ModeCountdown
				app.Timer.SetDuration(25)
			}
			if rl.IsKeyPressed(rl.KeyTwo) {
				app.Timer.Mode = ModeCountdown
				app.Timer.SetDuration(50)
			}
			if rl.IsKeyPressed(rl.KeyThree) {
				app.Timer.Mode = ModeCountdown
				app.Timer.SetDuration(5) // Break timer
			}
		}
	}

	return false
}

func handleSoundMenu(app *App) bool {
	if rl.IsKeyPressed(rl.KeyEscape) || rl.IsKeyPressed(rl.KeyS) {
		app.ShowSoundMenu = false
		return false
	}

	if rl.IsKeyPressed(rl.KeyOne) {
		app.Config.SoundFile = "bell"
		app.Sound.PlayPreview(app.Config)
		SaveConfig(app.Config)
	}
	if rl.IsKeyPressed(rl.KeyTwo) {
		app.Config.SoundFile = "chime"
		app.Sound.PlayPreview(app.Config)
		SaveConfig(app.Config)
	}
	if rl.IsKeyPressed(rl.KeyThree) {
		app.Config.SoundFile = "none"
		SaveConfig(app.Config)
	}

	return false
}
