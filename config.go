package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AppConfig holds all user-configurable settings.
type AppConfig struct {
	FocusDuration     int     `json:"focus_duration_minutes"`
	ShortBreak        int     `json:"short_break_minutes"`
	LongBreak         int     `json:"long_break_minutes"`
	SessionsBeforeLong int    `json:"sessions_before_long_break"`
	DefaultTimer      int     `json:"default_timer_minutes"`
	SoundFile         string  `json:"sound"`
	FontScale         float32 `json:"font_scale"`
	TaskName          string  `json:"task_name"`
	Volume            float32 `json:"volume"`
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
}

func LoadConfig() AppConfig {
	cfg := defaultConfig

	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		// No config file — write default and return
		SaveConfig(cfg)
		return cfg
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return defaultConfig
	}

	// Clamp values
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

	return cfg
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
