package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/zeozeozeo/fine"
)

var (
	gopher       *fine.Sprite
	gopherEntity *fine.Entity
)

func loadSprites(app *fine.App) {
	var err error
	gopher, err = app.NewSpriteFromPath("gopher.png")
	if err != nil {
		panic(err)
	}
}

func createEntities(app *fine.App) {
	gopherEntity = app.Entity(fine.NewVec2(1280/4, 720/4)).
		SetTexture(gopher).
		SetScale(fine.NewVec2(10, 10)).
		SetPivotCentered(true)

	for i := float64(0); i < 50; i += 10 {
		app.Rect(fine.NewVec2(100+i, 100+i), 200, 250, color.RGBA{255, 0, 0, 255})
	}

	for i := float64(0); i < 1000; i += 10 {
		app.Line(fine.NewVec2(200, 200), fine.NewVec2(200+i, 400+i/8), color.RGBA{0, 255, 0, 255})
	}
}

func main() {
	app := fine.NewApp("Demo", 1280, 720)
	app.SetUpdateFunc(update)
	loadSprites(app)
	createEntities(app)

	app.After(2, func(a *fine.App) {
		fmt.Println("Hello, World!")
	})

	// Start app
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// This function is called every frame.
// dt is the time elapsed to draw the previous frame,
// app is the application that called the function.
func update(dt float64, app *fine.App) {
	// sine := (math.Sin((4*app.Time+3)*math.Pi/2) + 1) / 2
	// app.Camera.SetPosition(fine.NewVec2(math.Sin(app.Time)*200, math.Sin(app.Time*5)*400)).
	// app.Camera.SetZoom(app.Time / 2)
	app.Camera.SetPosition(gopherEntity.Position)

	gopherEntity.SetAngle(math.Sin(app.Time) * 360)
	// SetOpacity((math.Sin((4*app.Time+3)*math.Pi/2) + 1) / 2)

	if app.IsKeyJustDown(fine.KEY_w) {
		fmt.Println("W is just down!")
	}

	if app.IsMouseButtonJustUp(fine.MBUTTON_LEFT) {
		fmt.Println("LMB is just up!")
	}

	gopherEntity.Scale.X = (math.Sin(app.Time*5) + 1) * 10
	gopherEntity.Scale.Y = (math.Sin(app.Time*3.5) + 1) * 10
}
