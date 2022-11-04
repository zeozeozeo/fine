package fine

import (
	"fmt"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// Creates a new window and starts the draw loop.
func (app *App) Run() error {
	runtime.LockOSThread()

	// Initialize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}
	defer sdl.Quit()

	// Create window
	if len(app.WindowFlags) < 1 {
		return fmt.Errorf("cannot initialize window without any window flags")
	}

	// Get window flags
	var finalFlags uint32
	for _, flag := range app.WindowFlags {
		finalFlags |= flag
	}

	if app.ScaleQuality != 0 {
		sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, fmt.Sprint(app.ScaleQuality))
	}

	// Create window
	window, err := sdl.CreateWindow(
		app.Title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		app.Width,
		app.Height,
		finalFlags,
	)
	if err != nil {
		return err
	}
	defer window.Destroy()
	app.Window = window

	// Create renderer
	if app.Renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED); err != nil {
		return err
	}
	app.Renderer.Clear()
	defer app.Renderer.Destroy()
	defer app.FreeSprites()

	// Start draw loop
	startTime := time.Now()
	if err := sdl.GLSetSwapInterval(app.SwapInterval); err != nil {
		return err
	}
	app.Running = true

	for app.Running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			// TODO: More events
			switch t := event.(type) {
			case *sdl.QuitEvent:
				// Ask close function if we need to close
				doCloseResponse := true
				if app.OnClose != nil {
					doCloseResponse = app.OnClose()
				}

				if !app.IgnoreClose && doCloseResponse {
					app.Running = false
					break
				}
			case *sdl.KeyboardEvent:
				app.handleKeyboardEvent(t)
			case *sdl.MouseButtonEvent:
				app.handleMouseButtonEvent(t)
			}
		}

		// Update
		app.Time = time.Since(startTime).Seconds()
		app.DeltaTime = app.Time - app.PreviousFrameTime
		app.PreviousFrameTime = app.Time

		// Draw
		app.DrawFrame()
		app.Renderer.Present()
	}

	return nil
}
