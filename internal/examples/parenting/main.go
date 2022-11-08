// Parents entity2 to entity1. Child entities will follow their parent.
package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/zeozeozeo/fine"
)

var (
	entity1 *fine.Entity
	entity2 *fine.Entity
)

func main() {
	app := fine.NewApp("Entity parenting", 1280, 720)

	entity1 = app.Rect(fine.NewVec2(10, -400), 100, 100, color.RGBA{255, 255, 128, 255}, true)
	entity2 = app.Rect(fine.NewVec2(200, -450), 100, 100, color.RGBA{0, 255, 255, 255}, true).
		ParentTo(entity1)

	entity1.SetUpdateFunc(func(dt float64, app *fine.App, entity *fine.Entity) {
		entity.Position.Y += 30 * dt
		entity2.Position.X = 200 * math.Sin(app.Time*5)
	})

	app.After(3, func(app *fine.App) {
		entity2.RemoveParent()
		fmt.Println("unparented entity2")
	})
	app.After(6, func(app *fine.App) {
		entity2.ParentTo(entity1)
		fmt.Println("parented entity2 to entity1")
	})

	// Start the application
	if err := app.Run(); err != nil {
		panic(err)
	}
}
