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

func baseEntity(app *App, position Vec2, w, h float64, color color.RGBA) *Entity {
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
	app.Scene.Entities = append(app.Scene.Entities, entity)
	return entity
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

// renderer *sdl.Renderer, vx, vy []int16, r, g, b, a uint8

type Polygon struct {
	Point1 Vec2 // First point.
	Point2 Vec2 // Second point.
	Point3 Vec2 // Third point.
	app    *App
	entity *Entity
}

func (app *App) Polygon(point1, point2, point3 Vec2, color color.RGBA) *Entity {
	entity := baseEntity(app, Vec2{}, 0, 0, color)
	entity.Shape = &Polygon{Point1: point1, Point2: point2, Point3: point3, app: app, entity: entity}
	return entity
}

func (poly *Polygon) Draw() {
	posX, posY := poly.app.Camera.WorldToScreen(
		NewVec2(
			poly.Point1.X*poly.app.Camera.Zoom,
			poly.Point1.Y*poly.app.Camera.Zoom,
		),
	)

	ax1 := int16(math.Round(float64(posX) + (float64(poly.app.Width) / 2)))
	ay1 := int16(math.Round(float64(posY) + (float64(poly.app.Height) / 2)))

	posX, posY = poly.app.Camera.WorldToScreen(
		NewVec2(
			poly.Point2.X*poly.entity.app.Camera.Zoom,
			poly.Point2.Y*poly.entity.app.Camera.Zoom,
		),
	)

	ax2 := int16(math.Round(float64(posX) + (float64(poly.app.Width) / 2)))
	ay2 := int16(math.Round(float64(posY) + (float64(poly.app.Height) / 2)))

	posX, posY = poly.app.Camera.WorldToScreen(
		NewVec2(
			poly.Point3.X*poly.entity.app.Camera.Zoom,
			poly.Point3.Y*poly.entity.app.Camera.Zoom,
		),
	)

	ax3 := int16(math.Round(float64(posX) + (float64(poly.app.Width) / 2)))
	ay3 := int16(math.Round(float64(posY) + (float64(poly.app.Height) / 2)))

	if poly.entity.Texture != nil {
		gfx.TexturedPolygon(
			poly.app.Renderer,
			[]int16{ax1, ax2, ax3},
			[]int16{ay1, ay2, ay3},
			poly.entity.Texture.Surface,
			// TODO: Custom texture offset
			0,
			0,
		)
	} else {
		gfx.PolygonRGBA(
			poly.app.Renderer,
			[]int16{ax1, ax2, ax3},
			[]int16{ay1, ay2, ay3},
			poly.entity.Color.R,
			poly.entity.Color.G,
			poly.entity.Color.B,
			poly.entity.Color.A,
		)
	}
}
