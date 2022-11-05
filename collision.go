package fine

import (
	"fmt"
	"math/rand"
)

// Checks collision with other entites.
// This feature is still expiremental. Suitable for platformer games.
func (entity *Entity) Collide() {
	// Check collisions with other entities
	pos1 := entity.Position
	for _, ent := range entity.Scene.Entities {
		if ent == entity {
			continue
		}
		pos2 := ent.Position

		if !(pos1.X+entity.Width < pos2.X ||
			pos1.Y+entity.Height < pos2.Y ||
			pos1.X > pos2.X+ent.Width ||
			pos1.Y > pos2.Y+ent.Height) {
			fmt.Printf("colliding %d\n", rand.Intn(500))
		}
	}
}

func (app *App) isRectOnScreen(x, y, w, h int32) bool {
	return !(x+w < 0 ||
		y+h < 0 ||
		x > app.Width+w ||
		y > app.Height+h)
}
