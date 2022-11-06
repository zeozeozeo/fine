// TODO: Load image from reader
package fine

import (
	"fmt"
	"io"
	"math"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Sprite struct {
	Surface   *sdl.Surface  // SDL surface.
	Tex       *sdl.Texture  // SDL texture.
	Width     int32         // Width of the sprite.
	Height    int32         // Height of the sprite.
	BlendMode sdl.BlendMode // Texture blend mode.
}

// Blend modes.
type BlendMode int

const (
	BLENDMODE_NONE  BlendMode = sdl.BLENDMODE_NONE  // No blending.
	BLENDMODE_BLEND BlendMode = sdl.BLENDMODE_BLEND // Alpha blending.
	BLENDMODE_ADD   BlendMode = sdl.BLENDMODE_ADD   // Additive blending.
	BLENDMODE_MOD   BlendMode = sdl.BLENDMODE_MOD   // Color modulate.
)

// Loads a new sprite from path. This uses the file extension
// to identify the format.
// Remember to call Free() on the surface when you're done with it.
func (app *App) NewSpriteFromPath(path string) (*Sprite, error) {
	surface, err := img.Load(path)
	if err != nil {
		return nil, err
	}

	return app.NewSpriteFromSurface(surface)
}

// Creates a new sprite from a reader.
func (app *App) NewSpriteFromReader(reader io.Reader) (*Sprite, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return app.NewSpriteFromData(data)
}

// Creates a new sprite from bytes.
func (app *App) NewSpriteFromData(data []byte) (*Sprite, error) {
	rwops, err := sdl.RWFromMem(data)
	if err != nil {
		return nil, err
	}

	surface, err := img.LoadRW(rwops, true)
	if err != nil {
		return nil, err
	}

	return app.NewSpriteFromSurface(surface)
}

// Creates a new sprite from an SDL surface.
func (app *App) NewSpriteFromSurface(surface *sdl.Surface) (*Sprite, error) {
	var sprite *Sprite

	if app.Running && app.Renderer != nil {
		texture, err := app.Renderer.CreateTextureFromSurface(surface)
		if err != nil {
			surface.Free()
			return nil, err
		}
		sprite = &Sprite{Surface: surface, Tex: texture, BlendMode: sdl.BLENDMODE_BLEND}
	} else {
		// The texture will be loaded later (when the texture is first rendered)
		sprite = &Sprite{Surface: surface, Tex: nil, BlendMode: sdl.BLENDMODE_BLEND}
	}

	sprite.Width = sprite.Surface.W
	sprite.Height = sprite.Surface.H

	app.LoadedSprites = append(app.LoadedSprites, sprite)
	return sprite, nil
}

// Loads the texture from the surface, if it exists.
func (sprite *Sprite) LoadTexture(app *App) error {
	if sprite.Surface == nil || !app.Running || app.Renderer == nil {
		return fmt.Errorf("renderer is not initialized, cannot load texture")
	}

	texture, err := app.Renderer.CreateTextureFromSurface(sprite.Surface)
	if err != nil {
		return err
	}

	sprite.Tex = texture
	return nil
}

// Frees the sprite's SDL texture and surface.
func (sprite *Sprite) Free() {
	if sprite.Surface != nil {
		sprite.Surface.Free()
	}
	if sprite.Tex != nil {
		sprite.Tex.Destroy()
	}
}

// Sets the sprite's blend mode.
// Available blend modes: BLENDMODE_NONE, BLENDMODE_BLEND, BLENDMODE_ADD, BLENDMODE_MOD.
func (sprite *Sprite) SetBlendMode(blendMode BlendMode) *Sprite {
	sprite.BlendMode = sdl.BlendMode(blendMode)
	return sprite
}

// Draws the sprite at the given screen coordinates.
func (sprite *Sprite) Render(app *App, x, y int, entity *Entity) {
	src := &sdl.Rect{
		X: 0,
		Y: 0,
		W: sprite.Width,
		H: sprite.Height,
	}

	var flip sdl.RendererFlip
	switch entity.FlipDir {
	case FLIP_NONE:
		flip = sdl.FLIP_NONE
	case FLIP_HORIZONTAL:
		flip = sdl.FLIP_HORIZONTAL
	case FLIP_VERTICAL:
		flip = sdl.FLIP_VERTICAL
	}

	// FIXME: This does not work for negative scales.
	entity.Width, entity.Height = float64(sprite.Width)*entity.Scale.X, float64(sprite.Height)*entity.Scale.Y
	dst := &sdl.Rect{
		X: int32(float64(x)) + app.Width/2,
		Y: int32(float64(y)) + app.Height/2,
		W: int32(entity.Width * app.Camera.Zoom),
		H: int32(entity.Height * app.Camera.Zoom),
	}

	var pivot Vec2
	if entity.IsPivotCentered {
		pivot.X = float64(dst.W) / 2
		pivot.Y = float64(dst.H) / 2
	} else {
		pivot = entity.Pivot
	}

	// Don't draw sprite if it's over the screen
	if !app.isRectOnScreen(dst.X, dst.Y, dst.W, dst.H) {
		return
	}

	sprite.Tex.SetBlendMode(sprite.BlendMode)
	sprite.Tex.SetAlphaMod(uint8(entity.Opacity * 255))

	app.Renderer.CopyEx(
		sprite.Tex,
		src,
		dst,
		entity.Angle,
		&sdl.Point{X: int32(math.Round(pivot.X)), Y: int32(math.Round(pivot.Y))},
		flip,
	)
}
