package fine

// Draws a new frame.
func (app *App) DrawFrame() error {
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
		// Follow parent (TODO: Rotate around parent)
		if entity.Parent != nil {
			entity.Position.X += entity.Parent.positionDelta.X
			entity.Position.Y += entity.Parent.positionDelta.Y
		}

		if entity.UpdateFunc != nil {
			entity.UpdateFunc(app.DeltaTime, app, entity)
		}
		if err := app.DrawEntity(entity); err != nil {
			return err
		}

		if entity.previousPosition.X != entity.Position.X || entity.previousPosition.Y != entity.Position.Y {
			entity.positionDelta = entity.Position.Sub(entity.previousPosition)
			entity.previousPosition = entity.Position
		}
	}

	// Check queued functions
	for idx, cronFunc := range app.QueuedFunctions {
		cronFunc.Left -= app.DeltaTime

		if cronFunc.Left <= 0 && !cronFunc.shouldRepeat {
			if cronFunc.Func != nil {
				cronFunc.Func(app)
			}
			app.QueuedFunctions = append(app.QueuedFunctions[:idx], app.QueuedFunctions[idx+1:]...)
		} else if cronFunc.Left <= 0 && cronFunc.shouldRepeat {
			cronFunc.Func(app)
			// Start counting down again
			cronFunc.Left = cronFunc.startTime
		}
	}

	if app.PostRender != nil {
		app.PostRender(app)
	}

	return nil
}

// Blits an entity to the screen.
func (app *App) DrawEntity(entity *Entity) error {
	if !entity.Visible || entity.Opacity == 0 {
		return nil
	}

	// Draw textures
	if entity.Texture != nil {
		if entity.Texture.Surface != nil && entity.Texture.Tex == nil {
			// The surface was loaded but the texture was not, load it immediately
			if err := entity.Texture.LoadTexture(app); err != nil {
				return err
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
	return nil
}
