package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/zeozeozeo/fine"
)

var (
	gopher                      *fine.Sprite
	gopherEntity, gopherEntity2 *fine.Entity
)

func loadSprites(app *fine.App) {
	var err error
	reader, err := os.Open("gopher.png")
	if err != nil {
		panic(err)
	}
	gopher, err = app.NewSpriteFromReader(reader)
	if err != nil {
		panic(err)
	}
}

func createEntities(app *fine.App) {
	gopherEntity = app.Entity(fine.NewVec2(1280/4, 720/4)).
		SetTexture(gopher).
		SetScale(fine.NewVec2(10, 10)).
		SetPivotCentered(true)

	gopherEntity2 = app.Entity(fine.NewVec2(1280/4, 720/4)).
		SetTexture(gopher).
		SetScale(fine.NewVec2(10, 10)).
		SetPivotCentered(true).SetAngle(90)

	app.ClearColor = color.RGBA{54, 54, 54, 255}
}

func main() {
	app := fine.NewApp("Demo", 1280, 720)
	app.SetUpdateFunc(update)
	loadSprites(app)
	createEntities(app)

	app.After(2, func(app *fine.App) {
		fmt.Println("Hello, World!")
	})

	// Start app
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func update(dt float64, app *fine.App) {
	gopherEntity.Position.Y += dt * 1000

	collision := gopherEntity.Collide()
	gopherEntity.Position.Y += collision.TopPushY
	fmt.Println(collision.TopPushY)

	// gopherEntity2.Position.X += dt * 100
}

/*
// This function is called every frame.
// dt is the time elapsed to draw the previous frame,
// app is the application that called the function.
func update(dt float64, app *fine.App) {
	// sine := (math.Sin((4*app.Time+3)*math.Pi/2) + 1) / 2
	// app.Camera.SetPosition(fine.NewVec2(math.Sin(app.Time)*200, math.Sin(app.Time*5)*400))
	// app.Camera.SetZoom(app.Time / 2)

	// gopherEntity.SetAngle(math.Sin(app.Time) * 360)
	// SetOpacity((math.Sin((4*app.Time+3)*math.Pi/2) + 1) / 2)

	if app.IsKeyJustDown(fine.KEY_w) {
		fmt.Println("W is just down!")
	}

	if app.IsMouseButtonJustUp(fine.MBUTTON_LEFT) {
		fmt.Println("LMB is just up!")
	}

	gopherEntity.Scale.X = (math.Sin(app.Time*5) + 1) * 10
	gopherEntity.Scale.Y = (math.Sin(app.Time*3.5) + 1) * 10
	// gopherEntity.FlipToggle(fine.FLIP_HORIZONTAL)
}
*/
