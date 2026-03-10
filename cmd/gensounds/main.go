package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

func main() {
	os.MkdirAll("assets/sounds", 0755)

	files := map[string][]byte{
		"assets/sounds/bell.wav":   generateBell(),
		"assets/sounds/chime.wav":  generateChime(),
		"assets/sounds/start.wav":  generateStartTone(),
		"assets/sounds/pause.wav":  generatePauseTone(),
		"assets/sounds/resume.wav": generateResumeTone(),
		"assets/sounds/reset.wav":  generateResetTone(),
	}

	for path, data := range files {
		if err := os.WriteFile(path, data, 0644); err != nil {
			fmt.Printf("Error writing %s: %v\n", path, err)
			continue
		}
		fmt.Printf("Generated %s (%d bytes)\n", path, len(data))
	}
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

// Alarm: bell — multi-harmonic sine with exponential decay, 1.5s
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

// Alarm: chime — two-tone sequence with overlapping decay, 1.8s
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

// UI: start — quick ascending sweep, 100ms
func generateStartTone() []byte {
	sampleRate := 44100
	duration := 0.1
	samples := int(float64(sampleRate) * duration)
	data := make([]int16, samples)
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		envelope := math.Exp(-t * 25.0)
		freq := 600.0 + 300.0*(t/duration)
		sample := math.Sin(2*math.Pi*freq*t) * envelope * 0.4
		data[i] = int16(sample * 32000)
	}
	return encodeWav(data, sampleRate)
}

// UI: pause — soft low tone, 80ms
func generatePauseTone() []byte {
	sampleRate := 44100
	duration := 0.08
	samples := int(float64(sampleRate) * duration)
	data := make([]int16, samples)
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		envelope := math.Exp(-t * 30.0)
		sample := math.Sin(2*math.Pi*400*t) * envelope * 0.35
		data[i] = int16(sample * 32000)
	}
	return encodeWav(data, sampleRate)
}

// UI: resume — ascending sweep, 100ms (slightly different from start)
func generateResumeTone() []byte {
	sampleRate := 44100
	duration := 0.1
	samples := int(float64(sampleRate) * duration)
	data := make([]int16, samples)
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		envelope := math.Exp(-t * 25.0)
		freq := 550.0 + 250.0*(t/duration)
		sample := math.Sin(2*math.Pi*freq*t) * envelope * 0.35
		data[i] = int16(sample * 32000)
	}
	return encodeWav(data, sampleRate)
}

// UI: reset — very short tick, 60ms
func generateResetTone() []byte {
	sampleRate := 44100
	duration := 0.06
	samples := int(float64(sampleRate) * duration)
	data := make([]int16, samples)
	for i := 0; i < samples; i++ {
		t := float64(i) / float64(sampleRate)
		envelope := math.Exp(-t * 40.0)
		sample := math.Sin(2*math.Pi*300*t) * envelope * 0.3
		data[i] = int16(sample * 32000)
	}
	return encodeWav(data, sampleRate)
}
