package main

import (
	"embed"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/fonts/Inter-Regular.ttf
//go:embed assets/fonts/Inter-SemiBold.ttf
//go:embed assets/fonts/Inter-Bold.ttf
var fontFS embed.FS

// Fonts holds all loaded font resources.
type Fonts struct {
	Regular  rl.Font
	SemiBold rl.Font
	Bold     rl.Font
}

// Full Unicode range for crisp rendering. We include ASCII + Latin supplements.
var fontChars = func() []rune {
	chars := make([]rune, 0, 256)
	for i := 32; i < 256; i++ {
		chars = append(chars, rune(i))
	}
	return chars
}()

// LoadFonts loads all Inter font variants at high resolution for crisp rendering.
func LoadFonts(scale float32) Fonts {
	baseSize := int32(200 * scale)
	if baseSize < 100 {
		baseSize = 100
	}

	return Fonts{
		Regular:  loadEmbeddedFont("assets/fonts/Inter-Regular.ttf", baseSize),
		SemiBold: loadEmbeddedFont("assets/fonts/Inter-SemiBold.ttf", baseSize),
		Bold:     loadEmbeddedFont("assets/fonts/Inter-Bold.ttf", baseSize),
	}
}

func loadEmbeddedFont(path string, size int32) rl.Font {
	data, err := fontFS.ReadFile(path)
	if err != nil {
		// Fallback to default font if embedded font fails
		return rl.GetFontDefault()
	}

	font := rl.LoadFontFromMemory(".ttf", data, size, fontChars)
	rl.SetTextureFilter(font.Texture, rl.FilterBilinear)
	return font
}

// UnloadFonts frees all font resources.
func UnloadFonts(f Fonts) {
	rl.UnloadFont(f.Regular)
	rl.UnloadFont(f.SemiBold)
	rl.UnloadFont(f.Bold)
}
