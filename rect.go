package fine

import (
	"image/color"

	"github.com/veandco/go-sdl2/sdl"
)

type Rectangle struct {
	Position Vec2    // The world position of the rect.
	W        float64 // The width of the rect.
	H        float64 // The height of the rect.
	Filled   bool    // Specifies whether the rect is filled or not.
	app      *App    // The app this rectangle belongs to.
	entity   *Entity // The entity this rectangle belongs to.
}

// Draws a rectangle to the screen.
func (rect *Rectangle) Draw() {
	// TODO: Don't draw rect if it's not on the screen
	rectX, rectY := rect.app.Camera.WorldToScreen(
		NewVec2(
			rect.Position.X*rect.app.Camera.Zoom,
			rect.Position.Y*rect.app.Camera.Zoom,
		),
	)

	prevR, prevG, prevB, prevA, err := rect.app.Renderer.GetDrawColor()
	if err != nil {
		prevR, prevG, prevB, prevA = 0, 0, 0, 0
	}

	rect.app.Renderer.SetDrawColor(
		rect.entity.Color.R,
		rect.entity.Color.G,
		rect.entity.Color.B,
		rect.entity.Color.A,
	)

	// TODO: check -1 scale values
	sdlRect := &sdl.Rect{
		X: int32(rectX),
		Y: int32(rectY),
		W: int32(rect.W * rect.app.Camera.Zoom * rect.entity.Scale.X),
		H: int32(rect.H * rect.app.Camera.Zoom * rect.entity.Scale.Y),
	}

	if rect.Filled {
		rect.app.Renderer.FillRect(sdlRect)
	} else {
		rect.app.Renderer.DrawRect(sdlRect)
	}

	rect.app.Renderer.SetDrawColor(prevR, prevG, prevB, prevA)
}

// Creates a new rectangle on the scene.
func (app *App) Rect(position Vec2, w, h float64, color color.RGBA, isFilled bool) *Entity {
	entity := &Entity{
		Scene:   app.Scene,
		Scale:   NewVec2(1, 1),
		Visible: true,
		Opacity: 1,
		Color:   color,
	}

	rectShape := &Rectangle{Position: position, W: w, H: h, Filled: isFilled, app: app, entity: entity}
	entity.Shape = rectShape

	app.Scene.Entities = append(app.Scene.Entities, entity)
	return entity
}

// Sets the rect's width and height.
func (rect *Rectangle) SetSize(w, h float64) *Rectangle {
	rect.W, rect.H = w, h
	return rect
}
