package main

import (
	_ "embed"
	"image/color"

	"github.com/zeozeozeo/fine"
)

//go:embed Tuffy_Bold.ttf
var fontData []byte
var text string = "Hello, world!"
var font *fine.Font

func main() {
	app := fine.NewApp("Drawing text", 1280, 720)

	// Load the font
	var err error
	font, err = app.LoadFontFromData(fontData, 96) // Font data, font size
	if err != nil {
		panic(err)
	}

	// Render the text to a texture
	renderedString := font.RenderStringOrPanic(text, color.RGBA{255, 255, 255, 255})

	// Create the entity holding the text
	app.Entity(fine.NewVec2(-float64(renderedString.Width), -float64(renderedString.Height)*2)).
		SetTexture(renderedString).
		SetUpdateFunc(entityUpdate)

	// Start app
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func entityUpdate(dt float64, app *fine.App, entity *fine.Entity) {
	// Warning: rendering text could take some time (~0.5ms),
	// so it's not the best idea to re-render it every frame.
	renderedString := font.RenderStringOrPanic(
		text[:int(app.Time*10)%(len(text)+1)],
		color.RGBA{255, 255, 255, 255},
	)
	app.FreeSprite(entity.Texture) // We don't need the previous texture anymore, free it
	entity.SetTexture(renderedString)
}
