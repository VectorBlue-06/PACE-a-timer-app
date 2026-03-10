package main

import rl "github.com/gen2brain/raylib-go/raylib"

// rl_GetFrameTime wraps raylib's GetFrameTime.
func rl_GetFrameTime() float32 {
	return rl.GetFrameTime()
}

func main() {
	// Configure window before creation
	rl.SetConfigFlags(rl.FlagBorderlessWindowedMode | rl.FlagMsaa4xHint | rl.FlagVsyncHint)

	// Get monitor size for borderless fullscreen
	monitor := rl.GetCurrentMonitor()
	monW := rl.GetMonitorWidth(monitor)
	monH := rl.GetMonitorHeight(monitor)
	if monW <= 0 || monH <= 0 {
		monW = 1920
		monH = 1080
	}

	rl.InitWindow(int32(monW), int32(monH), "do-it")
	rl.SetTargetFPS(60)
	rl.SetExitKey(0) // We handle ESC ourselves

	app := NewApp()
	app.InitGraphics()

	for !rl.WindowShouldClose() && !app.ShouldExit {
		app.Update()
		app.Draw()
	}

	app.Close()
	rl.CloseWindow()
}
