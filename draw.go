package fine

import (
	"log"
)

// Draws a new frame.
func (app *App) DrawFrame() {
	if app.Update != nil {
		app.Update(app.DeltaTime, app)
	}

	// Clear the slices of just up/down keys
	app.JustDownKeys = nil
	app.JustDownMouseButtons = nil
	app.JustUpKeys = nil
	app.JustUpMouseButtons = nil

	if app.DoClear {
		prevR, prevG, prevB, prevA, err := app.Renderer.GetDrawColor()
		if err != nil {
			prevR, prevG, prevB, prevA = 0, 0, 0, 0
		}
		app.Renderer.SetDrawColor(app.ClearColor.R, app.ClearColor.G, app.ClearColor.B, app.ClearColor.A)
		app.Renderer.Clear()
		app.Renderer.SetDrawColor(prevR, prevG, prevB, prevA)
	}

	// TODO: Proper layer system
	// Draw entities
	for _, entity := range app.Scene.Entities {
		app.DrawEntity(entity)
	}

	// Check queued functions
	for idx, cronFunc := range app.QueuedFunctions {
		cronFunc.Left -= app.DeltaTime
		if cronFunc.Left <= 0 {
			if cronFunc.Func != nil {
				cronFunc.Func(app)
			}
			app.QueuedFunctions = append(app.QueuedFunctions[:idx], app.QueuedFunctions[idx+1:]...)
		}
	}

	if app.PostRender != nil {
		app.PostRender(app)
	}
}

// Blits an entity to the screen.
func (app *App) DrawEntity(entity *Entity) {
	if !entity.Visible || entity.Opacity == 0 {
		return
	}

	// Draw textures
	if entity.Texture != nil {
		if entity.Texture.Surface != nil && entity.Texture.Tex == nil {
			// The surface was loaded but the texture was not, load it immediately
			if err := entity.Texture.LoadTexture(app); err != nil {
				log.Printf("[fatal] couldn't load entity texture: %s", err)
				app.Close()
				return
			}
		}

		screenX, screenY := app.Camera.WorldToScreen(
			NewVec2(
				entity.Position.X*app.Camera.Zoom,
				entity.Position.Y*app.Camera.Zoom,
			),
		)
		entity.Texture.Render(app, screenX, screenY, entity)
	} else if entity.Shape != nil {
		entity.Shape.Draw()
	}
}
