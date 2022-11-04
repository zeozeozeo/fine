package fine

import "image/color"

type Entity struct {
	Position        Vec2       // The world position of the entity.
	Scene           *Scene     // The scene this entity is in.
	Texture         *Sprite    // The texture (sprite) of the entity.
	Shape           Shape      // The shape of the entity. This is ignored if a texture is set.
	Color           color.RGBA // The color of the entity. This is ignored if a texture is set.
	Scale           Vec2       // X and Y scale of the entity.
	Opacity         float64    // Opacity of the entity between 0 and 1.
	Visible         bool       // Specifies if the entity should be rendered or not.
	Angle           float64    // Rotation angle.
	Pivot           Vec2       // Rotation pivot point. This is 0,0 by default.
	IsPivotCentered bool       // Should the pivot point be centered?
}

// Entity shapes. It must implement Draw(), which will be called
// when the entity needs to be rendered to the screen.
type Shape interface {
	Draw()
}

// Creates a new entity on the scene.
func (app *App) Entity(position Vec2) *Entity {
	entity := &Entity{
		Position: position,
		Scene:    app.Scene,
		Scale:    NewVec2(1, 1),
		Visible:  true,
		Opacity:  1,
	}

	app.Scene.Entities = append(app.Scene.Entities, entity)
	return entity
}

// Removes this entity from the scene.
func (entity *Entity) Destroy() {
	for idx, sceneEntity := range entity.Scene.Entities {
		if sceneEntity == entity {
			entity.Scene.Entities = append(entity.Scene.Entities[:idx], entity.Scene.Entities[idx+1:]...)
			return
		}
	}
}

// Sets the texture of the entity.
func (entity *Entity) SetTexture(sprite *Sprite) *Entity {
	entity.Texture = sprite
	return entity
}

// Sets entity X and Y scale.
func (entity *Entity) SetScale(scale Vec2) *Entity {
	entity.Scale = scale
	return entity
}

// Set the pivot point to be centered.
func (entity *Entity) SetPivotCentered(state bool) *Entity {
	entity.IsPivotCentered = state
	return entity
}

// Set the pivot point of the entity.
func (entity *Entity) SetPivot(pivot Vec2) *Entity {
	entity.Pivot = pivot
	return entity
}

// Sets the opacity of the entity. The opacity is a value from 0 to 1.
func (entity *Entity) SetOpacity(opacity float64) *Entity {
	if opacity > 1 {
		entity.Opacity = 1
	} else if opacity < 0 {
		entity.Opacity = 0
	} else {
		entity.Opacity = opacity
	}
	return entity
}

// Sets the rotation angle of the entity.
func (entity *Entity) SetAngle(angle float64) *Entity {
	entity.Angle = angle
	return entity
}

// Sets the color of the entity. The color is ignored if the entity
// has a texture.
func (entity *Entity) SetColor(color color.RGBA) *Entity {
	entity.Color = color
	return entity
}
