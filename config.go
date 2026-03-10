package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// KeyBinding stores a configurable keyboard shortcut.
type KeyBinding struct {
	Key  int32  `json:"key"`
	Ctrl bool   `json:"ctrl,omitempty"`
	Name string `json:"name"`
}

// AppConfig holds all user-configurable settings.
type AppConfig struct {
	FocusDuration      int     `json:"focus_duration_minutes"`
	ShortBreak         int     `json:"short_break_minutes"`
	LongBreak          int     `json:"long_break_minutes"`
	SessionsBeforeLong int     `json:"sessions_before_long_break"`
	DefaultTimer       int     `json:"default_timer_minutes"`
	SoundFile          string  `json:"sound"`
	FontScale          float32 `json:"font_scale"`
	TaskName           string  `json:"task_name"`
	Volume             float32 `json:"volume"`
	ShowProgressRing   bool    `json:"show_progress_ring"`
	EnableAnimations   bool    `json:"enable_animations"`
	FramePersistence   bool    `json:"frame_persistence"`
	AlarmSoundPath     string  `json:"alarm_sound_path,omitempty"`

	// Configurable key bindings (stored as raylib key codes)
	Keys KeyBindings `json:"keys"`
}

// KeyBindings holds all configurable shortcuts.
type KeyBindings struct {
	StartPause       KeyBinding `json:"start_pause"`
	Reset            KeyBinding `json:"reset"`
	FullscreenToggle KeyBinding `json:"fullscreen_toggle"`
	SettingsToggle   KeyBinding `json:"settings_toggle"`
	FadeToggle       KeyBinding `json:"fade_toggle"`
	PomodoroMode     KeyBinding `json:"pomodoro_mode"`
	SoundMenu        KeyBinding `json:"sound_menu"`
	Preset25         KeyBinding `json:"preset_25"`
	Preset50         KeyBinding `json:"preset_50"`
	PresetBreak      KeyBinding `json:"preset_break"`
	Quit             KeyBinding `json:"quit"`
}

// Raylib key constants used in config (duplicated to avoid importing rl in config)
const (
	keySpace  int32 = 32
	keyR      int32 = 82
	keyF      int32 = 70
	keyTab    int32 = 258
	keyP      int32 = 80
	keyS      int32 = 83
	keyOne    int32 = 49
	keyTwo    int32 = 50
	keyThree  int32 = 51
	keyEscape int32 = 256
)

func defaultKeyBindings() KeyBindings {
	return KeyBindings{
		StartPause:       KeyBinding{Key: keySpace, Name: "SPACE"},
		Reset:            KeyBinding{Key: keyR, Name: "R"},
		FullscreenToggle: KeyBinding{Key: keyF, Name: "F"},
		SettingsToggle:   KeyBinding{Key: keyTab, Name: "TAB"},
		FadeToggle:       KeyBinding{Key: keySpace, Ctrl: true, Name: "CTRL+SPACE"},
		PomodoroMode:     KeyBinding{Key: keyP, Name: "P"},
		SoundMenu:        KeyBinding{Key: keyS, Name: "S"},
		Preset25:         KeyBinding{Key: keyOne, Name: "1"},
		Preset50:         KeyBinding{Key: keyTwo, Name: "2"},
		PresetBreak:      KeyBinding{Key: keyThree, Name: "3"},
		Quit:             KeyBinding{Key: keyEscape, Name: "ESC"},
	}
}

var defaultConfig = AppConfig{
	FocusDuration:      25,
	ShortBreak:         5,
	LongBreak:          20,
	SessionsBeforeLong: 4,
	DefaultTimer:       25,
	SoundFile:          "bell",
	FontScale:          1.0,
	TaskName:           "Deep Work",
	Volume:             0.4,
	ShowProgressRing:   true,
	EnableAnimations:   true,
	FramePersistence:   false,
	Keys:               defaultKeyBindings(),
}

func LoadConfig() AppConfig {
	cfg := defaultConfig

	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		SaveConfig(cfg)
		return cfg
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaultConfig
	}

	clampConfig(&cfg)
	return cfg
}

func clampConfig(cfg *AppConfig) {
	if cfg.FontScale < 0.5 {
		cfg.FontScale = 0.5
	}
	if cfg.FontScale > 3.0 {
		cfg.FontScale = 3.0
	}
	if cfg.Volume < 0 {
		cfg.Volume = 0
	}
	if cfg.Volume > 1.0 {
		cfg.Volume = 1.0
	}
	if cfg.FocusDuration < 1 {
		cfg.FocusDuration = 25
	}
	if cfg.ShortBreak < 1 {
		cfg.ShortBreak = 5
	}
	if cfg.LongBreak < 1 {
		cfg.LongBreak = 20
	}
	if cfg.SessionsBeforeLong < 1 {
		cfg.SessionsBeforeLong = 4
	}
	// Ensure key bindings have defaults if zero
	if cfg.Keys.StartPause.Key == 0 {
		cfg.Keys = defaultKeyBindings()
	}
}

func SaveConfig(cfg AppConfig) {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(configPath(), data, 0644)
}

func configPath() string {
	exe, err := os.Executable()
	if err != nil {
		return "config.json"
	}
	return filepath.Join(filepath.Dir(exe), "config.json")
}
