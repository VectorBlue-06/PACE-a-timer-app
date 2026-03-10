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

func (s *SoundSystem) Init() {
	rl.InitAudioDevice()

	bellPath := writeTempWav("pace-bell.wav", generateBell())
	chimePath := writeTempWav("pace-chime.wav", generateChime())

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
	}
}

func (s *SoundSystem) PlayPreview(cfg AppConfig) {
	s.Play(cfg)
}

func (s *SoundSystem) Close() {
	if s.loaded {
		rl.UnloadSound(s.bellSound)
		rl.UnloadSound(s.chimeSound)
	}
	rl.CloseAudioDevice()
}

func generateBell() []byte {
	sampleRate := 44100
	duration := 1.5
	samples := int(float64(sampleRate) * duration)
	data := make([]int16, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		envelope := math.Exp(-t * 3.0)
		sample := 0.0
		sample += 0.6 * math.Sin(2*math.Pi*880*t)
		sample += 0.3 * math.Sin(2*math.Pi*1320*t)
		sample += 0.1 * math.Sin(2*math.Pi*1760*t)
		sample *= envelope * 0.5
		data[i] = int16(sample * 32000)
	}

	return encodeWav(data, sampleRate)
}

func generateChime() []byte {
	sampleRate := 44100
	duration := 1.8
	samples := int(float64(sampleRate) * duration)
	data := make([]int16, samples)

	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		sample := 0.0

		if t < 1.0 {
			env := math.Exp(-t * 2.5)
			sample += 0.5 * math.Sin(2*math.Pi*523*t) * env
			sample += 0.2 * math.Sin(2*math.Pi*785*t) * env
		}

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

func encodeWav(samples []int16, sampleRate int) []byte {
	dataSize := len(samples) * 2
	fileSize := 44 + dataSize

	buf := make([]byte, fileSize)
	copy(buf[0:4], "RIFF")
	binary.LittleEndian.PutUint32(buf[4:8], uint32(fileSize-8))
	copy(buf[8:12], "WAVE")

	copy(buf[12:16], "fmt ")
	binary.LittleEndian.PutUint32(buf[16:20], 16)
	binary.LittleEndian.PutUint16(buf[20:22], 1)
	binary.LittleEndian.PutUint16(buf[22:24], 1)
	binary.LittleEndian.PutUint32(buf[24:28], uint32(sampleRate))
	binary.LittleEndian.PutUint32(buf[28:32], uint32(sampleRate*2))
	binary.LittleEndian.PutUint16(buf[32:34], 2)
	binary.LittleEndian.PutUint16(buf[34:36], 16)

	copy(buf[36:40], "data")
	binary.LittleEndian.PutUint32(buf[40:44], uint32(dataSize))

	for i, s := range samples {
		binary.LittleEndian.PutUint16(buf[44+i*2:46+i*2], uint16(s))
	}

	return buf
}

func writeTempWav(name string, data []byte) string {
	dir := os.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return ""
	}
	return path
}
