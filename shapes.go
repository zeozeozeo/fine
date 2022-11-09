package fine

import (
	"image/color"
	"math"

	"github.com/veandco/go-sdl2/gfx"
)

// TODO: Implement all shapes from gfx.

// Returns x, y, width, height, do render.
func (entity *Entity) getDrawData() (int32, int32, int32, int32, bool) {
	if !entity.Visible {
		return 0, 0, 0, 0, false
	}

	posX, posY := entity.app.Camera.WorldToScreen(
		NewVec2(
			entity.Position.X*entity.app.Camera.Zoom,
			entity.Position.Y*entity.app.Camera.Zoom,
		),
	)

	ax := int32(math.Round(float64(posX) + (float64(entity.app.Width) / 2)))
	ay := int32(math.Round(float64(posY) + (float64(entity.app.Height) / 2)))
	aw := int32(entity.Width * entity.app.Camera.Zoom)
	ah := int32(entity.Height * entity.app.Camera.Zoom)
	if !entity.app.isRectOnScreen(ax, ay, aw, ah) {
		return 0, 0, 0, 0, false
	}

	return ax, ay, aw, ah, true
}

func (entity *Entity) getColors() (uint8, uint8, uint8, uint8) {
	return entity.Color.R,
		entity.Color.G,
		entity.Color.B,
		entity.Color.A
}

type CircleShape struct {
	HasAA  bool // Is the circle antialiased?
	Filled bool // Is the circle filled?
	Radius float64
	entity *Entity
	app    *App
}

func (circle *CircleShape) Draw() {
	x, y, _, _, draw := circle.entity.getDrawData()
	if !draw {
		return
	}
	r, g, b, a := circle.entity.getColors()

	switch {
	case circle.Filled:
		gfx.FilledCircleRGBA(circle.app.Renderer, x, y, int32(math.Round(circle.Radius)), r, g, b, a)
	case circle.HasAA:
		gfx.AACircleRGBA(circle.app.Renderer, x, y, int32(math.Round(circle.Radius)), r, g, b, a)
	default:
		gfx.CircleRGBA(circle.app.Renderer, x, y, int32(math.Round(circle.Radius)), r, g, b, a)
	}
}

// Creates a new circle.
// radius: the radius of the circle.
// filled: specifies whether the circle is filled.
// hasAA: specifies whether the circle has antialiasing.
func (app *App) Circle(position Vec2, radius float64, color color.RGBA, filled, hasAA bool) *Entity {
	entity := baseEntity(app, position, radius*2, radius*2, color)
	entity.Shape = &CircleShape{HasAA: hasAA, Filled: filled, entity: entity, app: app, Radius: radius}
	return entity
}

func baseEntity(app *App, position Vec2, w, h float64, color color.RGBA) *Entity {
	return &Entity{
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
}
