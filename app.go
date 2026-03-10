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
	ShowSettings  bool
	SettingsIndex int
	ForceClear    bool
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

	timer.OnComplete = func() {
		app.Sound.PlayAlarm(app.Config)
		name := AlarmDisplayName(app.Config)
		if name != "" {
			app.UI.TriggerSoundPopup(name)
		}
	}

	// Default timer setup
	if cfg.DefaultTimer <= 0 {
		cfg.DefaultTimer = 25
	}
	timer.Mode = ModeCountdown
	timer.SetDuration(cfg.DefaultTimer)

	return app
}

// InitGraphics sets up fonts and renderer (must be called after window creation).
func (app *App) InitGraphics() {
	fonts := LoadFonts(app.Config.FontScale)
	app.Renderer = NewRenderer(fonts, app.Config.FontScale)
	app.Sound.Init(app.Config)
}

// Update runs one frame of application logic.
func (app *App) Update() {
	// Use actual frame time for smooth animation
	dt := rl_GetFrameTime()
	if dt <= 0 || dt > 0.1 {
		dt = 1.0 / 60.0
	}

	app.Timer.Update(dt)
	app.UI.Update(dt, app.Timer.State, app.Config.EnableAnimations)

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
