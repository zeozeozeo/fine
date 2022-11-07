package fine

import (
	"image/color"
	"io"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Font struct {
	Data       *ttf.Font // SDL font data.
	Size       int       // Font size.
	WrapLength int       // Wrap length, negative values = no wrap (default: -1)
	app        *App
}

func (app *App) LoadFont(path string, size int) (*Font, error) {
	if err := app.loadFontLibraryIfNotLoaded(); err != nil {
		return nil, err
	}

	data, err := ttf.OpenFont(path, size)
	if err != nil {
		return nil, err
	}
	return &Font{
		Data: data,
		Size: size,
		app:  app,
	}, nil
}

func (app *App) LoadFontFromReader(reader io.Reader, size int) (*Font, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return app.LoadFontFromData(data, size)
}

// Loads a font from bytes.
func (app *App) LoadFontFromData(data []byte, size int) (*Font, error) {
	if err := app.loadFontLibraryIfNotLoaded(); err != nil {
		return nil, err
	}

	rwops, err := sdl.RWFromMem(data)
	if err != nil {
		return nil, err
	}

	fontData, err := ttf.OpenFontRW(rwops, 1, size)
	if err != nil {
		return nil, err
	}
	return &Font{
		Data: fontData,
		Size: size,
		app:  app,
	}, nil
}

func (app *App) loadFontLibraryIfNotLoaded() error {
	if !ttf.WasInit() {
		if err := ttf.Init(); err != nil {
			return err
		}
	}
	return nil
}

// Renders a string to a sprite. You can use this to set the entity's texture.
func (font *Font) RenderString(text string, color color.RGBA) (*Sprite, error) {
	var surface *sdl.Surface
	var err error
	if len(text) == 0 {
		text = " "
	}

	if font.WrapLength < 0 {
		surface, err = font.Data.RenderUTF8Blended(text, sdl.Color(color))
	} else {
		surface, err = font.Data.RenderUTF8BlendedWrapped(text, sdl.Color(color), font.WrapLength)
	}
	if err != nil {
		return nil, err
	}

	return font.app.NewSpriteFromSurface(surface)
}

// Renders a string to a sprite, panics if failed.
func (font *Font) RenderStringOrPanic(text string, color color.RGBA) *Sprite {
	sprite, err := font.RenderString(text, color)
	if err != nil {
		panic(err)
	}
	return sprite
}

func (app *App) UnloadFontLibrary() {
	if ttf.WasInit() {
		ttf.Quit()
	}
}
