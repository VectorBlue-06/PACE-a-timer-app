package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Color palette
var (
	ColorBG        = rl.NewColor(0, 0, 0, 255)       // #000000
	ColorPrimary   = rl.NewColor(255, 255, 255, 255)  // #FFFFFF
	ColorSecondary = rl.NewColor(138, 138, 138, 255)  // #8A8A8A
	ColorSubtle    = rl.NewColor(51, 51, 51, 255)     // #333333
	ColorAccent    = rl.NewColor(80, 200, 120, 255)   // Soft green for progress
	ColorBreak     = rl.NewColor(100, 160, 255, 255)  // Soft blue for break
	ColorComplete  = rl.NewColor(255, 200, 80, 255)   // Warm gold for completion
)

// Renderer handles all drawing operations.
type Renderer struct {
	fonts      Fonts
	fontScale  float32
	screenW    int32
	screenH    int32

	// Animation state
	fadeAlpha    float32
	pulsePhase  float64
	stateAlpha  float32
	targetAlpha float32
	prevState   TimerState
}

func NewRenderer(fonts Fonts, fontScale float32) *Renderer {
	return &Renderer{
		fonts:       fonts,
		fontScale:   fontScale,
		stateAlpha:  1.0,
		targetAlpha: 1.0,
		fadeAlpha:   0.0,
	}
}

// Update updates animation state each frame.
func (r *Renderer) Update(dt float32, state TimerState) {
	r.screenW = int32(rl.GetScreenWidth())
	r.screenH = int32(rl.GetScreenHeight())

	// Fade in on startup
	if r.fadeAlpha < 1.0 {
		r.fadeAlpha += dt * 3.0
		if r.fadeAlpha > 1.0 {
			r.fadeAlpha = 1.0
		}
	}

	// State change animation
	if state != r.prevState {
		r.stateAlpha = 0.7
		r.prevState = state
	}
	if r.stateAlpha < 1.0 {
		r.stateAlpha += dt * 5.0
		if r.stateAlpha > 1.0 {
			r.stateAlpha = 1.0
		}
	}

	// Pulsing for completed state
	if state == StateCompleted {
		r.pulsePhase += float64(dt) * 3.0
	} else {
		r.pulsePhase = 0
	}
}

// DrawFrame renders the complete UI frame.
func (r *Renderer) DrawFrame(app *App) {
	rl.BeginDrawing()
	rl.ClearBackground(ColorBG)

	alpha := r.fadeAlpha * r.stateAlpha
	if alpha < 0 {
		alpha = 0
	}

	r.drawLabel(app, alpha)
	r.drawTimer(app, alpha)
	r.drawProgressRing(app, alpha)
	r.drawSessionInfo(app, alpha)
	r.drawStatusHint(app, alpha)
	r.drawModeIndicator(app, alpha)
	r.drawSoundSelector(app, alpha)

	rl.EndDrawing()
}

// drawLabel draws the task name / phase label above the timer.
func (r *Renderer) drawLabel(app *App, alpha float32) {
	label := app.Config.TaskName
	if app.Timer.Mode == ModePomodoro {
		label = app.Pomodoro.PhaseLabel()
	}
	if app.Timer.Mode == ModeStopwatch {
		label = "Stopwatch"
	}

	fontSize := float32(r.screenH) * 0.025 * r.fontScale
	spacing := fontSize * 0.05
	size := rl.MeasureTextEx(r.fonts.Regular, label, fontSize, spacing)

	posX := (float32(r.screenW) - size.X) / 2
	posY := float32(r.screenH)*0.45 - float32(r.screenH)*0.09*r.fontScale - size.Y

	col := ColorSecondary
	col.A = uint8(float32(col.A) * alpha)
	rl.DrawTextEx(r.fonts.Regular, label, rl.NewVector2(posX, posY), fontSize, spacing, col)
}

// drawTimer draws the large MM:SS display.
func (r *Renderer) drawTimer(app *App, alpha float32) {
	text := app.Timer.DisplayString()
	fontSize := float32(r.screenH) * 0.14 * r.fontScale
	spacing := fontSize * 0.03

	size := rl.MeasureTextEx(r.fonts.Bold, text, fontSize, spacing)
	posX := (float32(r.screenW) - size.X) / 2
	posY := float32(r.screenH)*0.45 - size.Y/2

	col := ColorPrimary
	// Pulse effect on completion
	if app.Timer.State == StateCompleted {
		pulse := float32(0.6 + 0.4*math.Sin(r.pulsePhase))
		col.A = uint8(float32(col.A) * alpha * pulse)
	} else {
		col.A = uint8(float32(col.A) * alpha)
	}

	rl.DrawTextEx(r.fonts.Bold, text, rl.NewVector2(posX, posY), fontSize, spacing, col)
}

// drawProgressRing draws a thin circular progress indicator around the timer.
func (r *Renderer) drawProgressRing(app *App, alpha float32) {
	if app.Timer.Mode == ModeStopwatch {
		return
	}

	progress := float32(app.Timer.Progress())
	cx := float32(r.screenW) / 2
	cy := float32(r.screenH) * 0.45
	radius := float32(r.screenH) * 0.16 * r.fontScale

	// Background ring
	bgCol := ColorSubtle
	bgCol.A = uint8(float32(bgCol.A) * alpha * 0.5)
	drawArc(cx, cy, radius, 0, 360, 3.0, bgCol)

	// Progress ring
	var ringCol rl.Color
	switch {
	case app.Timer.State == StateCompleted:
		ringCol = ColorComplete
	case app.Timer.Mode == ModePomodoro && app.Pomodoro.Phase != PhaseFocus:
		ringCol = ColorBreak
	default:
		ringCol = ColorAccent
	}
	ringCol.A = uint8(float32(ringCol.A) * alpha)

	// Draw from top (270°), clockwise
	endAngle := progress * 360.0
	if endAngle > 0.5 {
		drawArc(cx, cy, radius, -90, -90+endAngle, 3.0, ringCol)
	}
}

// drawArc draws a smooth arc using line segments.
func drawArc(cx, cy, radius, startAngle, endAngle, thickness float32, color rl.Color) {
	segments := int(math.Abs(float64(endAngle-startAngle)) / 2.0)
	if segments < 16 {
		segments = 16
	}
	if segments > 180 {
		segments = 180
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

// drawSessionInfo draws session and today's focus stats at the bottom.
func (r *Renderer) drawSessionInfo(app *App, alpha float32) {
	if app.Timer.Mode != ModePomodoro {
		return
	}

	fontSize := float32(r.screenH) * 0.018 * r.fontScale
	spacing := fontSize * 0.03

	// Session info
	sessionText := app.Pomodoro.SessionLabel()
	size := rl.MeasureTextEx(r.fonts.Regular, sessionText, fontSize, spacing)
	posX := (float32(r.screenW) - size.X) / 2
	posY := float32(r.screenH) * 0.72

	col := ColorSubtle
	col.A = uint8(float32(col.A) * alpha)
	rl.DrawTextEx(r.fonts.Regular, sessionText, rl.NewVector2(posX, posY), fontSize, spacing, col)

	// Today's focus
	todayText := app.Pomodoro.TodayLabel()
	size = rl.MeasureTextEx(r.fonts.Regular, todayText, fontSize, spacing)
	posX = (float32(r.screenW) - size.X) / 2
	posY += fontSize * 1.8

	rl.DrawTextEx(r.fonts.Regular, todayText, rl.NewVector2(posX, posY), fontSize, spacing, col)
}

// drawStatusHint draws a subtle keyboard hint at the bottom of the screen.
func (r *Renderer) drawStatusHint(app *App, alpha float32) {
	var hint string
	switch app.Timer.State {
	case StateIdle:
		hint = "SPACE to start  ·  R reset  ·  1/2/3 presets  ·  P pomodoro  ·  W stopwatch  ·  S sound  ·  ESC quit"
	case StateRunning:
		hint = "SPACE to pause  ·  R reset  ·  ESC quit"
	case StatePaused:
		hint = "SPACE to resume  ·  R reset  ·  ESC quit"
	case StateCompleted:
		if app.Timer.Mode == ModePomodoro {
			hint = "SPACE for next  ·  R reset  ·  ESC quit"
		} else {
			hint = "SPACE to restart  ·  R reset  ·  ESC quit"
		}
	}

	fontSize := float32(r.screenH) * 0.014 * r.fontScale
	spacing := fontSize * 0.02
	size := rl.MeasureTextEx(r.fonts.Regular, hint, fontSize, spacing)

	posX := (float32(r.screenW) - size.X) / 2
	posY := float32(r.screenH) * 0.94

	col := ColorSubtle
	col.A = uint8(float32(50) * alpha)
	rl.DrawTextEx(r.fonts.Regular, hint, rl.NewVector2(posX, posY), fontSize, spacing, col)
}

// drawModeIndicator draws a subtle mode label at the top.
func (r *Renderer) drawModeIndicator(app *App, alpha float32) {
	var mode string
	switch app.Timer.Mode {
	case ModeCountdown:
		mode = "COUNTDOWN"
	case ModePomodoro:
		mode = "POMODORO"
	case ModeStopwatch:
		mode = "STOPWATCH"
	}

	fontSize := float32(r.screenH) * 0.013 * r.fontScale
	spacing := fontSize * 0.15
	size := rl.MeasureTextEx(r.fonts.SemiBold, mode, fontSize, spacing)

	posX := (float32(r.screenW) - size.X) / 2
	posY := float32(r.screenH) * 0.04

	col := ColorSubtle
	col.A = uint8(float32(40) * alpha)
	rl.DrawTextEx(r.fonts.SemiBold, mode, rl.NewVector2(posX, posY), fontSize, spacing, col)
}

// drawSoundSelector draws the sound selection overlay when active.
func (r *Renderer) drawSoundSelector(app *App, alpha float32) {
	if !app.ShowSoundMenu {
		return
	}

	// Semi-transparent overlay
	rl.DrawRectangle(0, 0, r.screenW, r.screenH, rl.NewColor(0, 0, 0, 200))

	titleSize := float32(r.screenH) * 0.03 * r.fontScale
	itemSize := float32(r.screenH) * 0.022 * r.fontScale
	spacing := titleSize * 0.03

	// Title
	title := "Select Sound"
	tSize := rl.MeasureTextEx(r.fonts.SemiBold, title, titleSize, spacing)
	posX := (float32(r.screenW) - tSize.X) / 2
	posY := float32(r.screenH) * 0.35

	col := ColorPrimary
	col.A = uint8(float32(col.A) * alpha)
	rl.DrawTextEx(r.fonts.SemiBold, title, rl.NewVector2(posX, posY), titleSize, spacing, col)

	// Options
	sounds := []struct {
		key  string
		name string
	}{
		{"1", "Bell"},
		{"2", "Chime"},
		{"3", "None"},
	}

	for i, s := range sounds {
		text := s.key + "  " + s.name
		iSize := rl.MeasureTextEx(r.fonts.Regular, text, itemSize, spacing)
		ix := (float32(r.screenW) - iSize.X) / 2
		iy := posY + tSize.Y + float32(i+1)*itemSize*2.0

		itemCol := ColorSecondary
		// Highlight current selection
		current := ""
		switch i {
		case 0:
			current = "bell"
		case 1:
			current = "chime"
		case 2:
			current = "none"
		}
		if app.Config.SoundFile == current {
			itemCol = ColorPrimary
		}
		itemCol.A = uint8(float32(itemCol.A) * alpha)
		rl.DrawTextEx(r.fonts.Regular, text, rl.NewVector2(ix, iy), itemSize, spacing, itemCol)
	}

	// Hint
	hint := "Press 1-3 to select  ·  S or ESC to close"
	hSize := rl.MeasureTextEx(r.fonts.Regular, hint, float32(r.screenH)*0.015*r.fontScale, spacing)
	hx := (float32(r.screenW) - hSize.X) / 2
	hy := float32(r.screenH) * 0.65

	hCol := ColorSubtle
	hCol.A = uint8(float32(60) * alpha)
	rl.DrawTextEx(r.fonts.Regular, hint, rl.NewVector2(hx, hy), float32(r.screenH)*0.015*r.fontScale, spacing, hCol)
}
