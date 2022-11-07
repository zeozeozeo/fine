package fine

import (
	"io"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Font struct {
	Data       *ttf.Font // SDL font data.
	Size       int       // Font size.
	WrapLength int       // Wrap length, negative values = no wrap (default: -1)
}

type TextShape struct {
	Text string
	Font *Font
}

func (app *App) LoadFont(path string, size int) (*Font, error) {
	data, err := ttf.OpenFont(path, size)
	if err != nil {
		return nil, err
	}
	return &Font{
		Data: data,
		Size: size,
	}, nil
}

func (app *App) LoadFontFromReader(reader io.Reader, size int) (*Font, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
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
	}, nil
}

func (app *App) Text(text string, position Vec2, font *Font) *Entity {
	// TODO
	return &Entity{}
}
