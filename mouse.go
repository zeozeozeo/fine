package fine

import (
	"github.com/veandco/go-sdl2/sdl"
)

type ButtonInfo struct {
	MouseID uint32      // The ID of the mouse.
	Button  MouseButton // The mouse button that is pressed.
	Clicks  int         // The amount of times the button has been clicked in a row (max: 255).
}

func (app *App) handleMouseButtonEvent(event *sdl.MouseButtonEvent) {
	button := MouseButton(event.Button)
	isPressed := app.IsMouseButtonDown(button)
	buttonInfo := ButtonInfo{
		MouseID: event.Which,
		Button:  button,
		Clicks:  int(event.Clicks),
	}

	if event.State == sdl.RELEASED {
		app.JustUpMouseButtons = append(app.JustUpMouseButtons, buttonInfo)
	}

	if event.State == sdl.PRESSED && !isPressed {
		app.DownMouseButtons = append(app.DownMouseButtons, buttonInfo)
		app.JustDownMouseButtons = append(app.JustDownMouseButtons, buttonInfo)
	} else if isPressed {
		// Remove pressed button
		for idx, downButton := range app.DownMouseButtons {
			if button == downButton.Button {
				app.DownMouseButtons = append(app.DownMouseButtons[:idx], app.DownMouseButtons[idx+1:]...)
				break
			}
		}
	}
}

// Get the mouse X and Y coordinates on screen.
func (app *App) GetMousePos() (int, int) {
	x, y, _ := sdl.GetMouseState()
	return int(x), int(y)
}

// Checks if a mouse button is currently pressed.
func (app *App) IsMouseButtonDown(button MouseButton) bool {
	for _, downButton := range app.DownMouseButtons {
		if button == downButton.Button {
			return true
		}
	}
	return false
}

// Checks if a mouse button is pressed on this frame.
func (app *App) IsMouseButtonJustDown(button MouseButton) bool {
	for _, justDownButton := range app.JustDownMouseButtons {
		if button == justDownButton.Button {
			return true
		}
	}
	return false
}

// Checks if a mouse button is released on this frame.
func (app *App) IsMouseButtonJustUp(button MouseButton) bool {
	for _, justUpButton := range app.JustUpMouseButtons {
		if button == justUpButton.Button {
			return true
		}
	}
	return false
}

// Mouse buttons.
type MouseButton int

const (
	MBUTTON_LEFT   MouseButton = sdl.BUTTON_LEFT   // Left mouse button.
	MBUTTON_MIDDLE MouseButton = sdl.BUTTON_MIDDLE // Middle mouse button.
	MBUTTON_RIGHT  MouseButton = sdl.BUTTON_RIGHT  // Right mouse button.
	MBUTTON_X1     MouseButton = sdl.BUTTON_X1     // X1 mouse button.
	MBUTTON_X2     MouseButton = sdl.BUTTON_X2     // X2 mouse button.
)
