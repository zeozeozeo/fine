package fine

import "image/color"

type Line struct {
	Start  Vec2    // Start of the line in the world.
	End    Vec2    // End of the line in the world.
	app    *App    // The app this line belongs to.
	entity *Entity // The entity this rectangle belongs to.
}

// Draws a line to the screen.
func (line *Line) Draw() {
	halfWidth, halfHeight := float64(line.app.Width)/2, float64(line.app.Height)/2
	x1, y1 := line.app.Camera.WorldToScreen(
		NewVec2(
			(line.Start.X*line.app.Camera.Zoom)+halfWidth,
			line.Start.Y*line.app.Camera.Zoom+halfHeight,
		),
	)

	x2, y2 := line.app.Camera.WorldToScreen(
		NewVec2(
			(line.End.X*line.app.Camera.Zoom)+halfWidth,
			(line.End.Y*line.app.Camera.Zoom)+halfHeight,
		),
	)

	prevR, prevG, prevB, prevA, err := line.app.Renderer.GetDrawColor()
	if err != nil {
		prevR, prevG, prevB, prevA = 0, 0, 0, 0
	}

	line.app.Renderer.SetDrawColor(
		line.entity.Color.R,
		line.entity.Color.G,
		line.entity.Color.B,
		line.entity.Color.A,
	)

	// TODO: Don't draw line if it's over the screen
	line.app.Renderer.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2))

	line.app.Renderer.SetDrawColor(prevR, prevG, prevB, prevA)
}

// Creates a new line on the scene.
func (app *App) Line(start, end Vec2, color color.RGBA) *Entity {
	entity := &Entity{
		Scene:   app.Scene,
		Scale:   NewVec2(1, 1),
		Visible: true,
		Opacity: 1,
		Color:   color,
		app: app,
	}

	lineShape := &Line{Start: start, End: end, app: app, entity: entity}
	entity.Shape = lineShape

	app.Scene.Entities = append(app.Scene.Entities, entity)
	return entity
}

// Set the start of the line in the world.
func (line *Line) SetStart(start Vec2) *Line {
	line.Start = start
	return line
}

// Set the end of the line in the world.
func (line *Line) SetEnd(end Vec2) *Line {
	line.End = end
	return line
}
