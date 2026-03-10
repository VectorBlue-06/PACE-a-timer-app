package main

// App holds the complete application state.
type App struct {
	Config   AppConfig
	Timer    *Timer
	Pomodoro *Pomodoro
	Renderer *Renderer
	Sound    *SoundSystem
	UI       *UI

	ShowSoundMenu bool
	ShouldExit    bool
}

// NewApp creates and initializes the full application.
func NewApp() *App {
	cfg := LoadConfig()
	timer := NewTimer()
	pomo := NewPomodoro(cfg, timer)

	app := &App{
		Config:   cfg,
		Timer:    timer,
		Pomodoro: pomo,
		Sound:    NewSoundSystem(),
		UI:       NewUI(),
	}

	// Set up timer completion callback
	timer.OnComplete = func() {
		app.Sound.Play(app.Config)
		if app.Timer.Mode == ModePomodoro {
			// Auto-advance handled by user pressing Space
		}
	}

	// Default timer setup
	switch cfg.DefaultTimer {
	case 0:
		timer.Mode = ModeStopwatch
	default:
		timer.Mode = ModeCountdown
		timer.SetDuration(cfg.DefaultTimer)
	}

	return app
}

// InitGraphics sets up fonts and renderer (must be called after window creation).
func (app *App) InitGraphics() {
	fonts := LoadFonts(app.Config.FontScale)
	app.Renderer = NewRenderer(fonts, app.Config.FontScale)
	app.Sound.Init()
}

// Update runs one frame of application logic.
func (app *App) Update() {
	dt := float32(1.0 / 60.0)
	if fps := float32(60); fps > 0 {
		dt = 1.0 / fps
	}
	// Use actual frame time for smoother animation
	frameTime := rl_GetFrameTime()
	if frameTime > 0 && frameTime < 0.1 {
		dt = frameTime
	}

	app.Timer.Update()
	app.UI.Update(dt, app.Timer.Mode, app.Timer.State)
	app.Renderer.Update(dt, app.Timer.State)

	if HandleInput(app) {
		app.ShouldExit = true
	}
}

// Draw renders the current frame.
func (app *App) Draw() {
	app.Renderer.DrawFrame(app)
}

// Close releases all resources.
func (app *App) Close() {
	UnloadFonts(app.Renderer.fonts)
	app.Sound.Close()
}
