package main

import rl "github.com/gen2brain/raylib-go/raylib"

// HandleInput processes all keyboard input and returns true if the app should exit.
func HandleInput(app *App) bool {
	keys := app.Config.Keys

	// Settings panel input takes priority when open
	if app.UI.SettingsFullyOpen() {
		return handleSettingsInput(app)
	}

	// Sound menu takes priority
	if app.ShowSoundMenu {
		return handleSoundMenu(app)
	}

	// Quit
	if isKeyPressed(keys.Quit) {
		return true
	}

	// Settings toggle (TAB)
	if isKeyPressed(keys.SettingsToggle) {
		app.ShowSettings = !app.ShowSettings
		app.UI.ShowSettings(app.ShowSettings)
		return false
	}

	// Frame persistence toggle (CTRL+SPACE)
	if isKeyComboPressed(keys.FadeToggle) {
		app.Config.FramePersistence = !app.Config.FramePersistence
		SaveConfig(app.Config)
		return false
	}

	// Start/Pause (Space — but only without ctrl)
	if isKeyPressed(keys.StartPause) && !rl.IsKeyDown(rl.KeyLeftControl) && !rl.IsKeyDown(rl.KeyRightControl) {
		if app.Timer.State == StateCompleted && app.Timer.Mode == ModePomodoro {
			app.Pomodoro.Advance()
			app.Timer.Start()
		} else {
			app.Timer.Toggle()
		}
	}

	// Reset
	if isKeyPressed(keys.Reset) {
		app.Timer.Reset()
		if app.Timer.Mode == ModePomodoro {
			app.Pomodoro.Phase = PhaseFocus
			app.Pomodoro.CurrentSession = 1
			app.Pomodoro.Setup()
		}
	}

	// Fullscreen
	if isKeyPressed(keys.FullscreenToggle) {
		rl.ToggleBorderlessWindowed()
	}

	// Pomodoro mode
	if isKeyPressed(keys.PomodoroMode) {
		app.Timer.Mode = ModePomodoro
		app.Timer.Reset()
		app.Pomodoro = NewPomodoro(app.Config, app.Timer)
		app.Pomodoro.Setup()
	}

	// Stopwatch mode
	if isKeyPressed(keys.StopwatchMode) {
		app.Timer.Mode = ModeStopwatch
		app.Timer.Reset()
	}

	// Sound selector
	if isKeyPressed(keys.SoundMenu) {
		app.ShowSoundMenu = true
	}

	// Number presets (only when idle/paused and not in pomodoro)
	if app.Timer.State == StateIdle || app.Timer.State == StatePaused {
		if app.Timer.Mode != ModePomodoro {
			if isKeyPressed(keys.Preset25) {
				app.Timer.Mode = ModeCountdown
				app.Timer.SetDuration(25)
			}
			if isKeyPressed(keys.Preset50) {
				app.Timer.Mode = ModeCountdown
				app.Timer.SetDuration(50)
			}
			if isKeyPressed(keys.PresetBreak) {
				app.Timer.Mode = ModeCountdown
				app.Timer.SetDuration(5)
			}
		}
	}

	return false
}

func handleSoundMenu(app *App) bool {
	if rl.IsKeyPressed(rl.KeyEscape) || isKeyPressed(app.Config.Keys.SoundMenu) {
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

// Settings panel navigation: 10 items
const settingsItemCount = 10

func handleSettingsInput(app *App) bool {
	// Close settings
	if rl.IsKeyPressed(rl.KeyTab) || rl.IsKeyPressed(rl.KeyEscape) {
		app.ShowSettings = false
		app.UI.ShowSettings(false)
		return false
	}

	// Navigate
	if rl.IsKeyPressed(rl.KeyUp) {
		app.SettingsIndex--
		if app.SettingsIndex < 0 {
			app.SettingsIndex = settingsItemCount - 1
		}
	}
	if rl.IsKeyPressed(rl.KeyDown) {
		app.SettingsIndex++
		if app.SettingsIndex >= settingsItemCount {
			app.SettingsIndex = 0
		}
	}

	// Adjust values
	left := rl.IsKeyPressed(rl.KeyLeft)
	right := rl.IsKeyPressed(rl.KeyRight)
	if !left && !right {
		return false
	}

	cfg := &app.Config
	delta := 1
	if left {
		delta = -1
	}

	switch app.SettingsIndex {
	case 0: // Focus Duration
		cfg.FocusDuration += delta * 5
		if cfg.FocusDuration < 1 {
			cfg.FocusDuration = 1
		}
		if cfg.FocusDuration > 120 {
			cfg.FocusDuration = 120
		}
	case 1: // Short Break
		cfg.ShortBreak += delta
		if cfg.ShortBreak < 1 {
			cfg.ShortBreak = 1
		}
		if cfg.ShortBreak > 30 {
			cfg.ShortBreak = 30
		}
	case 2: // Long Break
		cfg.LongBreak += delta * 5
		if cfg.LongBreak < 1 {
			cfg.LongBreak = 1
		}
		if cfg.LongBreak > 60 {
			cfg.LongBreak = 60
		}
	case 3: // Sessions
		cfg.SessionsBeforeLong += delta
		if cfg.SessionsBeforeLong < 1 {
			cfg.SessionsBeforeLong = 1
		}
		if cfg.SessionsBeforeLong > 10 {
			cfg.SessionsBeforeLong = 10
		}
	case 4: // Sound
		sounds := []string{"bell", "chime", "none"}
		idx := 0
		for i, s := range sounds {
			if s == cfg.SoundFile {
				idx = i
				break
			}
		}
		idx += delta
		if idx < 0 {
			idx = len(sounds) - 1
		}
		if idx >= len(sounds) {
			idx = 0
		}
		cfg.SoundFile = sounds[idx]
		app.Sound.PlayPreview(*cfg)
	case 5: // Volume
		cfg.Volume += float32(delta) * 0.1
		if cfg.Volume < 0 {
			cfg.Volume = 0
		}
		if cfg.Volume > 1.0 {
			cfg.Volume = 1.0
		}
	case 6: // Font Scale
		cfg.FontScale += float32(delta) * 0.1
		if cfg.FontScale < 0.5 {
			cfg.FontScale = 0.5
		}
		if cfg.FontScale > 3.0 {
			cfg.FontScale = 3.0
		}
	case 7: // Progress Ring
		cfg.ShowProgressRing = !cfg.ShowProgressRing
	case 8: // Animations
		cfg.EnableAnimations = !cfg.EnableAnimations
	case 9: // Frame Persistence
		cfg.FramePersistence = !cfg.FramePersistence
	}

	SaveConfig(*cfg)
	return false
}

// isKeyPressed checks if a key binding was pressed (non-ctrl shortcuts).
func isKeyPressed(kb KeyBinding) bool {
	if kb.Ctrl {
		return false // Use isKeyComboPressed for ctrl combos
	}
	return rl.IsKeyPressed(kb.Key)
}

// isKeyComboPressed checks if a ctrl+key combo was pressed.
func isKeyComboPressed(kb KeyBinding) bool {
	if !kb.Ctrl {
		return rl.IsKeyPressed(kb.Key)
	}
	return rl.IsKeyPressed(kb.Key) && (rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl))
}
