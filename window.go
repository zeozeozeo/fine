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

	defer app.FreeSprites()
	defer app.UnloadFontLibrary()

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
		finalFlags |= uint32(flag)
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

	// Start draw loop
	startTime := time.Now()
	if err := sdl.GLSetSwapInterval(app.SwapInterval); err != nil {
		return err
	}
	app.Running = true

	for app.Running {
		// app.ScrollDeltaX, app.ScrollDeltaY = 0, 0
		app.IsScrolling = false

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
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
			case *sdl.MouseWheelEvent:
				app.ScrollDeltaX, app.ScrollDeltaY = t.PreciseX, t.PreciseY
				app.IsScrolling = true
			}
		}

		// Update
		app.Time = time.Since(startTime).Seconds()
		app.DeltaTime = app.Time - app.PreviousFrameTime
		app.PreviousFrameTime = app.Time
		app.GetWindowSize() // This will update the window size for us :)

		// Draw
		if err := app.DrawFrame(); err != nil {
			return err
		}
		app.Renderer.Present()
	}

	return nil
}

type WindowFlag int

const (
	WINDOW_FULLSCREEN         WindowFlag = sdl.WINDOW_FULLSCREEN         // Fullscreen window.
	WINDOW_OPENGL             WindowFlag = sdl.WINDOW_OPENGL             // Window usable with OpenGL context.
	WINDOW_SHOWN              WindowFlag = sdl.WINDOW_SHOWN              // Window is visible.
	WINDOW_HIDDEN             WindowFlag = sdl.WINDOW_HIDDEN             // Window is not visible.
	WINDOW_BORDERLESS         WindowFlag = sdl.WINDOW_BORDERLESS         // No window decoration.
	WINDOW_RESIZABLE          WindowFlag = sdl.WINDOW_RESIZABLE          // Window can be resized.
	WINDOW_MINIMIZED          WindowFlag = sdl.WINDOW_MINIMIZED          // Window is minimized.
	WINDOW_MAXIMIZED          WindowFlag = sdl.WINDOW_MAXIMIZED          // Window is maximized.
	WINDOW_INPUT_GRABBED      WindowFlag = sdl.WINDOW_INPUT_GRABBED      // Window has grabbed input focus.
	WINDOW_INPUT_FOCUS        WindowFlag = sdl.WINDOW_INPUT_FOCUS        // Window has input focus.
	WINDOW_MOUSE_FOCUS        WindowFlag = sdl.WINDOW_MOUSE_FOCUS        // Window has mouse focus.
	WINDOW_FULLSCREEN_DESKTOP WindowFlag = sdl.WINDOW_FULLSCREEN_DESKTOP // Fullscreen window at the current desktop resolution.
	WINDOW_FOREIGN            WindowFlag = sdl.WINDOW_FOREIGN            // Window not created by SDL.
	WINDOW_ALLOW_HIGHDPI      WindowFlag = sdl.WINDOW_ALLOW_HIGHDPI      // Window should be created in high-DPI mode if supported (>= SDL 2.0.1).
	WINDOW_MOUSE_CAPTURE      WindowFlag = sdl.WINDOW_MOUSE_CAPTURE      // Window has mouse captured (unrelated to INPUT_GRABBED, >= SDL 2.0.4).
	WINDOW_ALWAYS_ON_TOP      WindowFlag = sdl.WINDOW_ALWAYS_ON_TOP      // Window should always be above others (X11 only, >= SDL 2.0.5).
	WINDOW_SKIP_TASKBAR       WindowFlag = sdl.WINDOW_SKIP_TASKBAR       // Window should not be added to the taskbar (X11 only, >= SDL 2.0.5).
	WINDOW_UTILITY            WindowFlag = sdl.WINDOW_UTILITY            // Window should be treated as a utility window (X11 only, >= SDL 2.0.5).
	WINDOW_TOOLTIP            WindowFlag = sdl.WINDOW_TOOLTIP            // Window should be treated as a tooltip (X11 only, >= SDL 2.0.5).
	WINDOW_POPUP_MENU         WindowFlag = sdl.WINDOW_POPUP_MENU         // Window should be treated as a popup menu (X11 only, >= SDL 2.0.5).
	WINDOW_VULKAN             WindowFlag = sdl.WINDOW_VULKAN             // Window usable for Vulkan surface (>= SDL 2.0.6).
)
