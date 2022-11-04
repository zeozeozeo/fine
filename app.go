package fine

import (
	"image/color"

	"github.com/veandco/go-sdl2/sdl"
)

type OnCloseFunc func() bool                      // Function that is called after the window closes. Should return whether the window should close or not.
type UpdateFunc func(float64, *App)               // Function that is called every frame.
type PostRenderFunc func(*App)                    // Function that is called at the end of drawing each frame.
type OnKeyEventFunc func(Key, KeyDirection, *App) // Function that is called when a key is down or up. Recieves the key, the direction and the app.

type App struct {
	Title       string         // Window title (TODO: Change title).
	Width       int32          // Window width.
	Height      int32          // Window height.
	Update      UpdateFunc     // Will be called every frame, recieves the delta time in seconds.
	PostRender  PostRenderFunc // Will be called after the rendering of all entities on a frame.
	OnClose     OnCloseFunc    // Will be called when the window should be closed. Should return whether the window should close or not.
	IgnoreClose bool           // Specifies if the window should ignore quit events.
	Running     bool           // Specifies if the app is running or not.
	WindowFlags []uint32       // SDL2 window flags. See https://github.com/veandco/go-sdl2/blob/1cc5f3cf35451b753c5838cb3760d2126b3d315e/sdl/video.go#L273
	DoClear     bool           // Do clear the previous frame before drawing the new frame?
	ClearColor  color.RGBA     // The color to clear the background of the new frame.

	// Time.

	Time              float64       // The amount of seconds since the app has started.
	PreviousFrameTime float64       // The time of the previous frame.
	DeltaTime         float64       // The duration between two frames (delta time).
	SwapInterval      int           // OpenGL swap interval. Default: 1 (vsync).
	QueuedFunctions   []*QueuedFunc // All functions that are queued to be called.

	Scene         *Scene        // The main scene of the app.
	Renderer      *sdl.Renderer // SDL renderer.
	LoadedSprites []*Sprite     // All loaded sprites.
	Camera        *Camera       // The main camera.
	ScaleQuality  int           // SDL scale quality.

	// Input.

	HeldKeys             []Key          // All keys that are currently held. Call app.IsKeyDown if you want to see if a key is down.
	JustDownKeys         []Key          // All keys that were pressed on this frame.
	JustUpKeys           []Key          // All keys that were released on this frame.
	OnKeyEvent           OnKeyEventFunc // Function that is called when a key is down or up.
	DownMouseButtons     []ButtonInfo   // Pressed mouse buttons information.
	JustDownMouseButtons []ButtonInfo   // Mouse buttons pressed on this frame information.
	JustUpMouseButtons   []ButtonInfo   // Mouse buttons released on this frame information.
}

// Creates a new app with a window title and size.
func NewApp(title string, width, height int32) *App {
	return &App{
		Title:  title,
		Width:  width,
		Height: height,
		WindowFlags: []uint32{
			sdl.WINDOW_SHOWN,
			sdl.WINDOW_OPENGL,
		},
		SwapInterval: 1,
		Scene:        &Scene{},
		DoClear:      true,
		Camera:       &Camera{Position: NewVec2(0, 0), Zoom: 1},
	}
}

// Sets the update function (called every frame).
func (app *App) SetUpdateFunc(newFunc UpdateFunc) *App {
	app.Update = newFunc
	return app
}

// Sets the close function (called before the window closes).
func (app *App) SetCloseFunc(newFunc OnCloseFunc) *App {
	app.OnClose = newFunc
	return app
}

// Sets the post render function (called when the frame is finished rendering).
func (app *App) SetPostRenderFunc(newFunc PostRenderFunc) *App {
	app.PostRender = newFunc
	return app
}

// Sets the post render function (called when the frame is finished rendering).
func (app *App) SetKeyEventFunc(newFunc OnKeyEventFunc) *App {
	app.OnKeyEvent = newFunc
	return app
}

// Sets the window width and height.
func (app *App) SetWindowSize(width, height int32) *App {
	app.Width, app.Height = width, height
	return app
}

// Sets the title of the window.
func (app *App) SetWindowTitle(title string) *App {
	// TODO: Make the title update after app.Run is called
	app.Title = title
	return app
}

// Closes the window. This will work even if app.IgnoreClose is true.
func (app *App) Close() {
	app.Running = false
}

// Returns an estimate amount of frames per second.
func (app *App) GetFPS() float64 {
	return 1 / app.DeltaTime
}

// Sets the OpenGL swap interval.
func (app *App) SetSwapInterval(interval int) error {
	var err error
	if app.Running {
		err = sdl.GLSetSwapInterval(interval)
	}
	app.SwapInterval = interval
	return err
}

// Fills the entire screen with an RGBA color.
func (app *App) FillScreen(r, g, b, a uint8) {
	if app.Running {
		prevR, prevG, prevB, prevA, err := app.Renderer.GetDrawColor()
		if err != nil {
			prevR, prevG, prevB, prevA = 0, 0, 0, 0
		}
		app.Renderer.SetDrawColor(r, g, b, a)
		app.Renderer.Clear()
		app.Renderer.SetDrawColor(prevR, prevG, prevB, prevA)
	}
}

// Frees all loaded sprites.
// Returns the amount of free'd sprites.
func (app *App) FreeSprites() int {
	amount := 0
	for idx, sprite := range app.LoadedSprites {
		sprite.Free()

		// FIXME: append() is slow for removing elements from slices, see
		// https://stackoverflow.com/a/37335777
		app.LoadedSprites = append(app.LoadedSprites[:idx], app.LoadedSprites[idx+1:]...)
		amount++
	}
	return amount
}

// Returns the text in the OS clipboard.
func (app *App) GetClipboard() string {
	clipboardText, err := sdl.GetClipboardText()
	if err != nil {
		return ""
	}
	return clipboardText
}

// Sets text in the OS clipboard.
func (app *App) SetClipboard(text string) error {
	return sdl.SetClipboardText(text)
}
