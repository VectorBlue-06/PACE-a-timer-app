package main

import (
	"encoding/binary"
	"math"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// SoundSystem manages notification sounds.
type SoundSystem struct {
	bellSound  rl.Sound
	chimeSound rl.Sound
	loaded     bool
}

func NewSoundSystem() *SoundSystem {
	return &SoundSystem{}
}

// Init generates and loads notification sounds.
func (s *SoundSystem) Init() {
	rl.InitAudioDevice()

	// Generate sounds to temp files, load, then clean up
	bellPath := writeTempWav("do-it-bell.wav", generateBell())
	chimePath := writeTempWav("do-it-chime.wav", generateChime())

	if bellPath != "" {
		s.bellSound = rl.LoadSound(bellPath)
		os.Remove(bellPath)
	}
	if chimePath != "" {
		s.chimeSound = rl.LoadSound(chimePath)
		os.Remove(chimePath)
	}

	s.loaded = true
}

// Play plays the configured notification sound.
func (s *SoundSystem) Play(cfg AppConfig) {
	if !s.loaded {
		return
	}

	vol := cfg.Volume
	switch cfg.SoundFile {
	case "bell":
		rl.SetSoundVolume(s.bellSound, vol)
		rl.PlaySound(s.bellSound)
	case "chime":
		rl.SetSoundVolume(s.chimeSound, vol)
		rl.PlaySound(s.chimeSound)
	case "none":
		// Silent
	}
}

// PlayPreview plays a preview of the selected sound.
func (s *SoundSystem) PlayPreview(cfg AppConfig) {
	s.Play(cfg)
}

// Close releases audio resources.
func (s *SoundSystem) Close() {
	if s.loaded {
		rl.UnloadSound(s.bellSound)
		rl.UnloadSound(s.chimeSound)
	}
	rl.CloseAudioDevice()
}

// generateBell creates a soft bell-like tone (sine wave with exponential decay).
func generateBell() []byte {
	sampleRate := 44100
	duration := 1.5
	samples := int(float64(sampleRate) * duration)
	data := make([]int16, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		// Multi-harmonic bell sound
		envelope := math.Exp(-t * 3.0)
		sample := 0.0
		sample += 0.6 * math.Sin(2*math.Pi*880*t)  // fundamental
		sample += 0.3 * math.Sin(2*math.Pi*1320*t)  // 3rd harmonic
		sample += 0.1 * math.Sin(2*math.Pi*1760*t)  // 2nd overtone
		sample *= envelope * 0.5
		data[i] = int16(sample * 32000)
	}

	return encodeWav(data, sampleRate)
}

// generateChime creates a gentle two-tone chime.
func generateChime() []byte {
	sampleRate := 44100
	duration := 1.8
	samples := int(float64(sampleRate) * duration)
	data := make([]int16, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		sample := 0.0

		// First tone: C5 (523 Hz)
		if t < 1.0 {
			env := math.Exp(-t * 2.5)
			sample += 0.5 * math.Sin(2*math.Pi*523*t) * env
			sample += 0.2 * math.Sin(2*math.Pi*785*t) * env
		}

		// Second tone: E5 (659 Hz), starts at 0.3s
		if t > 0.3 {
			t2 := t - 0.3
			env := math.Exp(-t2 * 2.5)
			sample += 0.5 * math.Sin(2*math.Pi*659*t2) * env
			sample += 0.2 * math.Sin(2*math.Pi*988*t2) * env
		}

		sample *= 0.4
		data[i] = int16(sample * 32000)
	}

	return encodeWav(data, sampleRate)
}

// encodeWav creates a minimal 16-bit mono WAV file.
func encodeWav(samples []int16, sampleRate int) []byte {
	dataSize := len(samples) * 2
	fileSize := 44 + dataSize

	buf := make([]byte, fileSize)
	copy(buf[0:4], "RIFF")
	binary.LittleEndian.PutUint32(buf[4:8], uint32(fileSize-8))
	copy(buf[8:12], "WAVE")

	// fmt chunk
	copy(buf[12:16], "fmt ")
	binary.LittleEndian.PutUint32(buf[16:20], 16)           // chunk size
	binary.LittleEndian.PutUint16(buf[20:22], 1)            // PCM
	binary.LittleEndian.PutUint16(buf[22:24], 1)            // mono
	binary.LittleEndian.PutUint32(buf[24:28], uint32(sampleRate))
	binary.LittleEndian.PutUint32(buf[28:32], uint32(sampleRate*2)) // byte rate
	binary.LittleEndian.PutUint16(buf[32:34], 2)            // block align
	binary.LittleEndian.PutUint16(buf[34:36], 16)           // bits per sample

	// data chunk
	copy(buf[36:40], "data")
	binary.LittleEndian.PutUint32(buf[40:44], uint32(dataSize))

	for i, s := range samples {
		binary.LittleEndian.PutUint16(buf[44+i*2:46+i*2], uint16(s))
	}

	return buf
}

// writeTempWav writes WAV data to a temp file and returns the path.
func writeTempWav(name string, data []byte) string {
	dir := os.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return ""
	}
	return path
}
