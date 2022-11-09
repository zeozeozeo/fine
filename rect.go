package fine

import (
	"image/color"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
)

type Rectangle struct {
	Filled   bool    // Specifies whether the rect is filled or not.
	Rounding int32   // The rounding of the rectangle, 0 = disabled.
	app      *App    // The app this rectangle belongs to.
	entity   *Entity // The entity this rectangle belongs to.
}

// Draws a rectangle to the screen.
func (rect *Rectangle) Draw() {
	rectX, rectY := rect.app.Camera.WorldToScreen(
		NewVec2(
			rect.entity.Position.X*rect.app.Camera.Zoom,
			rect.entity.Position.Y*rect.app.Camera.Zoom,
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

	rect.entity.Width, rect.entity.Height = rect.entity.Width*rect.entity.Scale.X, rect.entity.Height*rect.entity.Scale.X
	sdlRect := &sdl.Rect{
		X: int32(float64(rectX)) + rect.app.Width/2,
		Y: int32(float64(rectY)) + rect.app.Height/2,
		W: int32(rect.entity.Width * rect.app.Camera.Zoom),
		H: int32(rect.entity.Height * rect.app.Camera.Zoom),
	}

	if !rect.app.isRectOnScreen(sdlRect.X, sdlRect.Y, sdlRect.W, sdlRect.H) {
		return
	}

	switch {
	case rect.Filled:
		if rect.Rounding != 0 {
			rect.app.Renderer.FillRect(sdlRect)
		} else {
			gfx.RoundedRectangleRGBA(
				rect.app.Renderer,
				sdlRect.X,
				sdlRect.Y,
				sdlRect.X+sdlRect.W,
				sdlRect.Y+sdlRect.H,
				rect.Rounding,
				rect.entity.Color.R,
				rect.entity.Color.G,
				rect.entity.Color.B,
				rect.entity.Color.A,
			)
		}
	default:
		rect.app.Renderer.DrawRect(sdlRect)
	}

	rect.app.Renderer.SetDrawColor(prevR, prevG, prevB, prevA)
}

// Creates a new rectangle on the scene.
func (app *App) Rect(position Vec2, w, h float64, color color.RGBA, isFilled bool) *Entity {
	entity := &Entity{
		Position:         position,
		Scene:            app.Scene,
		Scale:            NewVec2(1, 1),
		Visible:          true,
		Opacity:          1,
		Color:            color,
		Width:            w,
		Height:           h,
		DoCollide:        true,
		app:              app,
		previousPosition: position,
	}

	rectShape := &Rectangle{Filled: isFilled, app: app, entity: entity}
	entity.Shape = rectShape

	app.Scene.Entities = append(app.Scene.Entities, entity)
	return entity
}
