package fine

type CollisionInfo struct {
	ForwardPushX  float64 // The amount to add to the X position to push the entity out of the collision to the forward.
	BackwardPushX float64 // The amount to add to the X position to push the entity out of the collision to the backward.
	TopPushY      float64 // The amount to add to the Y position to push the entity out of the collision to the top.
	BottomPushY   float64 // The amount to add to the Y position to push the entity out of the collision to the bottom.
	IsFalling     bool    // This is true if the entity isn't being pushed to the top.
	Entity        *Entity // The entity this collision happened with.
}

// Checks collision with other entites.
// This feature is still experimental. Suitable for platformer games.
func (entity *Entity) Collide() CollisionInfo {
	// Check collisions with other entities
	pos1 := entity.Position
	collision := CollisionInfo{}

	for _, ent := range entity.Scene.Entities {
		if ent == entity || !ent.DoCollide {
			continue
		}
		pos2 := ent.Position

		// Find collisions
		if !(pos1.X+entity.Width < pos2.X ||
			pos1.Y+entity.Height < pos2.Y ||
			pos1.X > pos2.X+ent.Width ||
			pos1.Y > pos2.Y+ent.Height) {
			//////////////////////////////////////////////////////////
			collision.ForwardPushX = (pos2.X + ent.Width) - pos1.X
			collision.BackwardPushX = pos2.X - (pos1.X + entity.Width)
			collision.TopPushY = pos2.Y - (pos1.Y + entity.Height)
			collision.BottomPushY = (pos2.Y + ent.Height) - pos1.Y
			collision.Entity = ent
			break
		}
	}

	// If the entity is not being pushed to the top, it is falling
	// NOTE:  This won't work if the entity isn't constantly falling.
	// FIXME: Instead, check if this entity has an another entity under itself.
	//        That will only work with platformer games, though...
	if collision.TopPushY == 0 {
		collision.IsFalling = true
	}

	return collision
}

// Checks if a rect is on the screen (should be drawn) or not.
func (app *App) isRectOnScreen(x, y, w, h int32) bool {
	return !(x+w < 0 || y+h < 0 || x > app.Width+w || y > app.Height+h)
}
