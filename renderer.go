package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Color palette
var (
	ColorBG        = rl.NewColor(0, 0, 0, 255)
	ColorPrimary   = rl.NewColor(255, 255, 255, 255)
	ColorSecondary = rl.NewColor(170, 170, 170, 255)
	ColorSubtle    = rl.NewColor(51, 51, 51, 255)
	ColorAccent    = rl.NewColor(80, 200, 120, 255)
	ColorBreak     = rl.NewColor(100, 160, 255, 255)
	ColorComplete  = rl.NewColor(255, 200, 80, 255)
)

// Renderer handles all drawing operations with correct layering.
type Renderer struct {
	fonts     Fonts
	fontScale float32
	screenW   float32
	screenH   float32

	// Precomputed stable digit width for "00:00" style rendering
	digitWidth float32
	colonWidth float32
	refSize    float32 // font size used for precomputed widths

	// Cached hint string to avoid per-frame formatting
	cachedHint      string
	cachedHintState TimerState
	cachedHintMode  TimerMode
}

func NewRenderer(fonts Fonts, fontScale float32) *Renderer {
	return &Renderer{
		fonts:     fonts,
		fontScale: fontScale,
	}
}

// updateMetrics refreshes screen dimensions and stable digit measurements.
func (r *Renderer) updateMetrics() {
	r.screenW = float32(rl.GetScreenWidth())
	r.screenH = float32(rl.GetScreenHeight())

	// Precompute widest digit width for stable centering
	fontSize := r.screenH * 0.14 * r.fontScale
	if fontSize != r.refSize {
		r.refSize = fontSize
		spacing := fontSize * 0.03
		maxW := float32(0)
		for d := 0; d <= 9; d++ {
			s := fmt.Sprintf("%d", d)
			w := rl.MeasureTextEx(r.fonts.Bold, s, fontSize, spacing)
			if w.X > maxW {
				maxW = w.X
			}
		}
		r.digitWidth = maxW
		cw := rl.MeasureTextEx(r.fonts.Bold, ":", fontSize, spacing)
		r.colonWidth = cw.X
	}
}

// DrawFrame renders the complete UI frame with correct layer order.
func (r *Renderer) DrawFrame(app *App) {
	r.updateMetrics()

	rl.BeginDrawing()

	// Frame persistence mode: very subtle overlay (opacity ~0.05)
	if app.Config.FramePersistence && !app.ForceClear {
		rl.DrawRectangle(0, 0, int32(r.screenW), int32(r.screenH), rl.NewColor(0, 0, 0, 13))
	} else {
		rl.ClearBackground(ColorBG)
		app.ForceClear = false
	}

	alpha := app.UI.EffectiveAlpha()

	// Layer 1: Progress ring (behind text)
	if app.Config.ShowProgressRing {
		r.drawProgressRing(app, alpha)
	}

	// Layer 2: Mode indicator (top)
	r.drawModeIndicator(app, alpha)

	// Layer 3: Label text (above timer)
	r.drawLabel(app, alpha)

	// Layer 4: Timer text (center, on top of ring)
	r.drawTimer(app, alpha)

	// Layer 5: Session and today info
	r.drawSessionInfo(app, alpha)
	r.drawTodayFocus(app, alpha)

	// Layer 6: Sound name popup
	r.drawSoundPopup(app, alpha)

	// Layer 7: Keyboard hints (bottom)
	r.drawStatusHint(app, alpha)

	// Layer 8: Overlays (settings, sound selector) — always on top
	r.drawSoundSelector(app, alpha)
	r.drawSettingsPanel(app)

	rl.EndDrawing()
}

// drawLabel draws the task/phase label above the timer.
func (r *Renderer) drawLabel(app *App, alpha float32) {
	label := app.Config.TaskName
	if app.Timer.Mode == ModePomodoro {
		label = app.Pomodoro.PhaseLabel()
	}

	fontSize := r.screenH * 0.025 * r.fontScale
	spacing := fontSize * 0.05
	size := rl.MeasureTextEx(r.fonts.Regular, label, fontSize, spacing)

	posX := (r.screenW - size.X) / 2
	posY := r.screenH*0.45 - r.screenH*0.10*r.fontScale - size.Y

	col := ColorSecondary
	col.A = uint8(float32(col.A) * alpha)
	rl.DrawTextEx(r.fonts.Regular, label, rl.NewVector2(posX, posY), fontSize, spacing, col)
}

// drawTimer draws the large MM:SS display with stable character-cell rendering.
// Each digit is drawn in a fixed-width cell to prevent horizontal jitter.
func (r *Renderer) drawTimer(app *App, alpha float32) {
	text := app.Timer.DisplayString() // "MM:SS"
	fontSize := r.screenH * 0.14 * r.fontScale
	spacing := fontSize * 0.03
	scale := app.UI.TimerScale()

	// Total width: 4 digits + 1 colon, with small gap
	gap := r.digitWidth * 0.08
	totalW := r.digitWidth*4 + r.colonWidth + gap*4
	scaledW := totalW * scale
	scaledFontSize := fontSize * scale

	startX := (r.screenW - scaledW) / 2
	textH := rl.MeasureTextEx(r.fonts.Bold, "0", scaledFontSize, spacing)
	posY := r.screenH*0.45 - textH.Y/2

	// Determine alpha: blink + digit transition + pulse
	visible := app.UI.TimerVisible()
	digitA := app.Timer.DigitAlpha()

	var baseAlpha float32
	if !visible {
		baseAlpha = 0
	} else {
		baseAlpha = alpha * digitA
		if app.Timer.State == StateCompleted {
			baseAlpha *= app.UI.PulseAlpha()
		}
	}

	col := ColorPrimary
	col.A = uint8(float32(col.A) * baseAlpha)

	// Draw each character in a fixed cell
	x := startX
	for i := 0; i < 5; i++ {
		ch := string(text[i])
		var cellW float32
		if text[i] == ':' {
			cellW = r.colonWidth * scale
		} else {
			cellW = r.digitWidth * scale
		}

		// Center character within its cell
		charSize := rl.MeasureTextEx(r.fonts.Bold, ch, scaledFontSize, spacing)
		charX := x + (cellW-charSize.X)/2

		rl.DrawTextEx(r.fonts.Bold, ch, rl.NewVector2(charX, posY), scaledFontSize, spacing, col)

		x += cellW + gap*scale
	}
}

// drawProgressRing draws a thin circular progress indicator behind the timer.
func (r *Renderer) drawProgressRing(app *App, alpha float32) {
	progress := float32(app.Timer.Progress())
	cx := r.screenW / 2
	cy := r.screenH * 0.45
	radius := r.screenH * 0.18 * r.fontScale

	// Background ring
	bgCol := ColorSubtle
	bgCol.A = uint8(float32(bgCol.A) * alpha * 0.4)
	drawArc(cx, cy, radius, 0, 360, 2.5, bgCol)

	// Progress ring color
	var ringCol rl.Color
	switch {
	case app.Timer.State == StateCompleted:
		ringCol = ColorComplete
	case app.Timer.Mode == ModePomodoro && app.Pomodoro.Phase != PhaseFocus:
		ringCol = ColorBreak
	default:
		ringCol = ColorAccent
	}
	ringCol.A = uint8(float32(ringCol.A) * alpha * 0.9)

	endAngle := progress * 360.0
	if endAngle > 0.5 {
		drawArc(cx, cy, radius, -90, -90+endAngle, 2.5, ringCol)
	}
}

func drawArc(cx, cy, radius, startAngle, endAngle, thickness float32, color rl.Color) {
	segments := int(math.Abs(float64(endAngle-startAngle)) / 1.5)
	if segments < 24 {
		segments = 24
	}
	if segments > 240 {
		segments = 240
	}

	step := (endAngle - startAngle) / float32(segments)
	for i := 0; i < segments; i++ {
		a1 := (startAngle + step*float32(i)) * math.Pi / 180.0
		a2 := (startAngle + step*float32(i+1)) * math.Pi / 180.0

		x1 := cx + radius*float32(math.Cos(float64(a1)))
		y1 := cy + radius*float32(math.Sin(float64(a1)))
		x2 := cx + radius*float32(math.Cos(float64(a2)))
		y2 := cy + radius*float32(math.Sin(float64(a2)))

		rl.DrawLineEx(rl.NewVector2(x1, y1), rl.NewVector2(x2, y2), thickness, color)
	}
}

// drawSessionInfo draws session info in the upper portion below mode indicator.
func (r *Renderer) drawSessionInfo(app *App, alpha float32) {
	if app.Timer.Mode != ModePomodoro {
		return
	}

	fontSize := r.screenH * 0.016 * r.fontScale
	spacing := fontSize * 0.03

	sessionText := app.Pomodoro.SessionLabel()
	size := rl.MeasureTextEx(r.fonts.Regular, sessionText, fontSize, spacing)
	posX := (r.screenW - size.X) / 2
	posY := r.screenH * 0.06

	col := ColorSubtle
	col.A = uint8(float32(col.A) * alpha * 0.8)
	rl.DrawTextEx(r.fonts.Regular, sessionText, rl.NewVector2(posX, posY), fontSize, spacing, col)
}

// drawTodayFocus draws total focus time in the lower area.
func (r *Renderer) drawTodayFocus(app *App, alpha float32) {
	if app.Timer.Mode != ModePomodoro {
		return
	}

	fontSize := r.screenH * 0.016 * r.fontScale
	spacing := fontSize * 0.03

	todayText := app.Pomodoro.TodayLabel()
	size := rl.MeasureTextEx(r.fonts.Regular, todayText, fontSize, spacing)
	posX := (r.screenW - size.X) / 2
	posY := r.screenH * 0.85

	col := ColorSubtle
	col.A = uint8(float32(col.A) * alpha * 0.6)
	rl.DrawTextEx(r.fonts.Regular, todayText, rl.NewVector2(posX, posY), fontSize, spacing, col)
}

// drawSoundPopup shows the alarm sound name in the bottom-right corner.
func (r *Renderer) drawSoundPopup(app *App, alpha float32) {
	popupAlpha := app.UI.SoundPopupAlpha()
	if popupAlpha <= 0 {
		return
	}

	text := app.UI.SoundPopupText()
	if text == "" {
		return
	}

	fontSize := r.screenH * 0.014 * r.fontScale
	spacing := fontSize * 0.03
	size := rl.MeasureTextEx(r.fonts.Regular, text, fontSize, spacing)

	posX := r.screenW - size.X - r.screenW*0.03
	posY := r.screenH * 0.92

	col := ColorSecondary
	col.A = uint8(float32(col.A) * alpha * popupAlpha * 0.7)
	rl.DrawTextEx(r.fonts.Regular, text, rl.NewVector2(posX, posY), fontSize, spacing, col)
}

// drawStatusHint draws contextual keyboard hints at the bottom.
func (r *Renderer) drawStatusHint(app *App, alpha float32) {
	state := app.Timer.State
	mode := app.Timer.Mode

	// Cache hint string — only rebuild when state or mode changes
	if r.cachedHint == "" || state != r.cachedHintState || mode != r.cachedHintMode {
		r.cachedHintState = state
		r.cachedHintMode = mode
		keys := app.Config.Keys
		switch state {
		case StateIdle:
			r.cachedHint = keys.StartPause.Name + " start  \u00b7  " + keys.Reset.Name + " reset  \u00b7  " + keys.SettingsToggle.Name + " settings  \u00b7  " + keys.Quit.Name + " quit"
		case StateRunning:
			r.cachedHint = keys.StartPause.Name + " pause  \u00b7  " + keys.Reset.Name + " reset  \u00b7  " + keys.Quit.Name + " quit"
		case StatePaused:
			r.cachedHint = keys.StartPause.Name + " resume  \u00b7  " + keys.Reset.Name + " reset  \u00b7  " + keys.Quit.Name + " quit"
		case StateCompleted:
			r.cachedHint = keys.StartPause.Name + " continue  \u00b7  " + keys.Reset.Name + " reset  \u00b7  " + keys.Quit.Name + " quit"
		}
	}

	hint := r.cachedHint

	fontSize := r.screenH * 0.013 * r.fontScale
	spacing := fontSize * 0.02
	size := rl.MeasureTextEx(r.fonts.Regular, hint, fontSize, spacing)

	posX := (r.screenW - size.X) / 2
	posY := r.screenH * 0.94

	col := ColorSubtle
	col.A = uint8(float32(40) * alpha)
	rl.DrawTextEx(r.fonts.Regular, hint, rl.NewVector2(posX, posY), fontSize, spacing, col)
}

// drawModeIndicator draws a subtle mode label at the very top.
func (r *Renderer) drawModeIndicator(app *App, alpha float32) {
	var mode string
	switch app.Timer.Mode {
	case ModeCountdown:
		mode = "COUNTDOWN"
	case ModePomodoro:
		mode = "POMODORO"
	}

	fontSize := r.screenH * 0.012 * r.fontScale
	spacing := fontSize * 0.15
	size := rl.MeasureTextEx(r.fonts.SemiBold, mode, fontSize, spacing)

	posX := (r.screenW - size.X) / 2
	posY := r.screenH * 0.03

	col := ColorSubtle
	col.A = uint8(float32(35) * alpha)
	rl.DrawTextEx(r.fonts.SemiBold, mode, rl.NewVector2(posX, posY), fontSize, spacing, col)
}

// drawSoundSelector draws the sound selection overlay when active.
func (r *Renderer) drawSoundSelector(app *App, alpha float32) {
	if !app.ShowSoundMenu {
		return
	}

	rl.DrawRectangle(0, 0, int32(r.screenW), int32(r.screenH), rl.NewColor(0, 0, 0, 200))

	titleSize := r.screenH * 0.03 * r.fontScale
	itemSize := r.screenH * 0.022 * r.fontScale
	spacing := titleSize * 0.03

	title := "Select Sound"
	tSize := rl.MeasureTextEx(r.fonts.SemiBold, title, titleSize, spacing)
	posX := (r.screenW - tSize.X) / 2
	posY := r.screenH * 0.35

	col := ColorPrimary
	col.A = uint8(float32(col.A) * alpha)
	rl.DrawTextEx(r.fonts.SemiBold, title, rl.NewVector2(posX, posY), titleSize, spacing, col)

	sounds := []struct {
		key, name, id string
	}{
		{"1", "Bell", "bell"},
		{"2", "Chime", "chime"},
		{"3", "None", "none"},
	}

	for i, s := range sounds {
		text := s.key + "  " + s.name
		iSize := rl.MeasureTextEx(r.fonts.Regular, text, itemSize, spacing)
		ix := (r.screenW - iSize.X) / 2
		iy := posY + tSize.Y + float32(i+1)*itemSize*2.0

		itemCol := ColorSecondary
		if app.Config.SoundFile == s.id {
			itemCol = ColorPrimary
		}
		itemCol.A = uint8(float32(itemCol.A) * alpha)
		rl.DrawTextEx(r.fonts.Regular, text, rl.NewVector2(ix, iy), itemSize, spacing, itemCol)
	}

	hint := "Press 1-3 to select  ·  S or ESC to close"
	hintSize := r.screenH * 0.015 * r.fontScale
	hSize := rl.MeasureTextEx(r.fonts.Regular, hint, hintSize, spacing)
	hx := (r.screenW - hSize.X) / 2
	hy := r.screenH * 0.65

	hCol := ColorSubtle
	hCol.A = uint8(float32(55) * alpha)
	rl.DrawTextEx(r.fonts.Regular, hint, rl.NewVector2(hx, hy), hintSize, spacing, hCol)
}

// drawSettingsPanel draws the settings overlay with fade animation.
func (r *Renderer) drawSettingsPanel(app *App) {
	if !app.UI.SettingsVisible() {
		return
	}

	sa := app.UI.SettingsAlpha()

	// Dimmed background
	rl.DrawRectangle(0, 0, int32(r.screenW), int32(r.screenH),
		rl.NewColor(0, 0, 0, uint8(220*sa)))

	titleSize := r.screenH * 0.035 * r.fontScale
	labelSize := r.screenH * 0.018 * r.fontScale
	valueSize := r.screenH * 0.018 * r.fontScale
	spacing := labelSize * 0.03

	// Title
	title := "Settings"
	tSize := rl.MeasureTextEx(r.fonts.SemiBold, title, titleSize, spacing)
	tx := (r.screenW - tSize.X) / 2
	ty := r.screenH * 0.15

	titleCol := ColorPrimary
	titleCol.A = uint8(float32(titleCol.A) * sa)
	rl.DrawTextEx(r.fonts.SemiBold, title, rl.NewVector2(tx, ty), titleSize, spacing, titleCol)

	// Settings items
	type settingsItem struct {
		label string
		value string
	}

	alarmName := "Not set"
	if app.Config.AlarmSoundPath != "" {
		alarmName = AlarmDisplayName(app.Config)
		if alarmName == "" {
			alarmName = "Not set"
		}
	}

	items := []settingsItem{
		{"Focus Duration", fmt.Sprintf("%d min", app.Config.FocusDuration)},
		{"Short Break", fmt.Sprintf("%d min", app.Config.ShortBreak)},
		{"Long Break", fmt.Sprintf("%d min", app.Config.LongBreak)},
		{"Sessions", fmt.Sprintf("%d", app.Config.SessionsBeforeLong)},
		{"Alarm Sound", app.Config.SoundFile},
		{"Volume", fmt.Sprintf("%.0f%%", app.Config.Volume*100)},
		{"Font Scale", fmt.Sprintf("%.1fx", app.Config.FontScale)},
		{"Progress Ring", boolLabel(app.Config.ShowProgressRing)},
		{"Animations", boolLabel(app.Config.EnableAnimations)},
		{"Frame Persistence", boolLabel(app.Config.FramePersistence)},
		{"Custom Alarm File", alarmName},
	}

	startY := ty + tSize.Y + r.screenH*0.04
	rowH := labelSize * 2.4
	colLabel := r.screenW*0.35
	colValue := r.screenW*0.55
	selIdx := app.SettingsIndex

	for i, item := range items {
		y := startY + float32(i)*rowH

		// Highlight selected row
		labelCol := ColorSecondary
		valCol := ColorPrimary
		if i == selIdx {
			labelCol = ColorPrimary
			valCol = ColorAccent
		}
		labelCol.A = uint8(float32(labelCol.A) * sa)
		valCol.A = uint8(float32(valCol.A) * sa)

		rl.DrawTextEx(r.fonts.Regular, item.label, rl.NewVector2(colLabel, y), labelSize, spacing, labelCol)
		rl.DrawTextEx(r.fonts.SemiBold, item.value, rl.NewVector2(colValue, y), valueSize, spacing, valCol)
	}

	// Navigation hint
	hint := "↑↓ navigate  ·  ←→ adjust  ·  ENTER browse  ·  TAB close"
	hintSize := r.screenH * 0.013 * r.fontScale
	hSize := rl.MeasureTextEx(r.fonts.Regular, hint, hintSize, spacing)
	hx := (r.screenW - hSize.X) / 2
	hy := r.screenH * 0.88

	hCol := ColorSubtle
	hCol.A = uint8(float32(50) * sa)
	rl.DrawTextEx(r.fonts.Regular, hint, rl.NewVector2(hx, hy), hintSize, spacing, hCol)
}

func boolLabel(b bool) string {
	if b {
		return "On"
	}
	return "Off"
}
