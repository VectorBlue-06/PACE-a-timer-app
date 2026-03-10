package main

import (
	"embed"
	"os"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/sounds/bell.wav
//go:embed assets/sounds/chime.wav
//go:embed assets/sounds/start.wav
//go:embed assets/sounds/pause.wav
//go:embed assets/sounds/resume.wav
//go:embed assets/sounds/reset.wav
var soundFS embed.FS

// SoundSystem manages alarm and UI feedback sounds.
type SoundSystem struct {
	// Alarm sounds (user-configurable)
	bellSound        rl.Sound
	chimeSound       rl.Sound
	customAlarmSound rl.Sound
	customAlarmLoaded bool

	// UI feedback sounds (not configurable, always available)
	startSound  rl.Sound
	pauseSound  rl.Sound
	resumeSound rl.Sound
	resetSound  rl.Sound

	loaded bool
}

func NewSoundSystem() *SoundSystem {
	return &SoundSystem{}
}

func (s *SoundSystem) Init(cfg AppConfig) {
	rl.InitAudioDevice()

	s.bellSound = s.loadEmbeddedSound("bell.wav")
	s.chimeSound = s.loadEmbeddedSound("chime.wav")
	s.startSound = s.loadEmbeddedSound("start.wav")
	s.pauseSound = s.loadEmbeddedSound("pause.wav")
	s.resumeSound = s.loadEmbeddedSound("resume.wav")
	s.resetSound = s.loadEmbeddedSound("reset.wav")

	if cfg.AlarmSoundPath != "" {
		s.LoadCustomAlarm(cfg.AlarmSoundPath)
	}

	s.loaded = true
}

func (s *SoundSystem) loadEmbeddedSound(name string) rl.Sound {
	data, err := soundFS.ReadFile("assets/sounds/" + name)
	if err != nil {
		return rl.Sound{}
	}
	tmpPath := writeTempWav("pace-"+name, data)
	if tmpPath == "" {
		return rl.Sound{}
	}
	snd := rl.LoadSound(tmpPath)
	os.Remove(tmpPath)
	return snd
}

// LoadCustomAlarm loads an external WAV/MP3 file as the alarm sound.
func (s *SoundSystem) LoadCustomAlarm(path string) {
	if s.customAlarmLoaded {
		rl.UnloadSound(s.customAlarmSound)
		s.customAlarmLoaded = false
	}
	if path == "" {
		return
	}
	if _, err := os.Stat(path); err != nil {
		return
	}
	s.customAlarmSound = rl.LoadSound(path)
	s.customAlarmLoaded = true
}

// PlayAlarm plays the configured alarm sound (timer completion).
func (s *SoundSystem) PlayAlarm(cfg AppConfig) {
	if !s.loaded {
		return
	}
	vol := cfg.Volume
	if cfg.SoundFile == "custom" && s.customAlarmLoaded {
		rl.SetSoundVolume(s.customAlarmSound, vol)
		rl.PlaySound(s.customAlarmSound)
		return
	}
	switch cfg.SoundFile {
	case "bell":
		rl.SetSoundVolume(s.bellSound, vol)
		rl.PlaySound(s.bellSound)
	case "chime":
		rl.SetSoundVolume(s.chimeSound, vol)
		rl.PlaySound(s.chimeSound)
	}
}

func (s *SoundSystem) PlayPreview(cfg AppConfig) {
	s.PlayAlarm(cfg)
}

const uiSoundVolume = 0.15

func (s *SoundSystem) PlayStart() {
	if !s.loaded {
		return
	}
	rl.SetSoundVolume(s.startSound, uiSoundVolume)
	rl.PlaySound(s.startSound)
}

func (s *SoundSystem) PlayPause() {
	if !s.loaded {
		return
	}
	rl.SetSoundVolume(s.pauseSound, uiSoundVolume)
	rl.PlaySound(s.pauseSound)
}

func (s *SoundSystem) PlayResume() {
	if !s.loaded {
		return
	}
	rl.SetSoundVolume(s.resumeSound, uiSoundVolume)
	rl.PlaySound(s.resumeSound)
}

func (s *SoundSystem) PlayReset() {
	if !s.loaded {
		return
	}
	rl.SetSoundVolume(s.resetSound, uiSoundVolume)
	rl.PlaySound(s.resetSound)
}

func (s *SoundSystem) Close() {
	if s.loaded {
		rl.UnloadSound(s.bellSound)
		rl.UnloadSound(s.chimeSound)
		rl.UnloadSound(s.startSound)
		rl.UnloadSound(s.pauseSound)
		rl.UnloadSound(s.resumeSound)
		rl.UnloadSound(s.resetSound)
		if s.customAlarmLoaded {
			rl.UnloadSound(s.customAlarmSound)
		}
	}
	rl.CloseAudioDevice()
}

// AlarmDisplayName returns the display name for the current alarm sound.
func AlarmDisplayName(cfg AppConfig) string {
	switch cfg.SoundFile {
	case "bell":
		return "Bell"
	case "chime":
		return "Chime"
	case "custom":
		if cfg.AlarmSoundPath != "" {
			base := filepath.Base(cfg.AlarmSoundPath)
			ext := filepath.Ext(base)
			return strings.TrimSuffix(base, ext)
		}
		return ""
	}
	return ""
}

func writeTempWav(name string, data []byte) string {
	dir := os.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return ""
	}
	return path
}
