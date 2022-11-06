package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"math"

	"github.com/zeozeozeo/fine"
)

var (
	gopher          *fine.Entity
	gopherFallSpeed float64
	gopherVelocity  float64
	gopherDirection string // "left" or "right"
	walkStartTime   float64
	cameraVelocity  float64
)

//go:embed gopher.png
var gopherSpriteData []byte

func main() {
	app := fine.NewApp("Demo", 1280, 720)
	app.SetUpdateFunc(update)
	createEntities(app)
	createLevel(app)
	app.Camera.Position = gopher.Position

	app.After(2, func(app *fine.App) {
		fmt.Println("Hello, World after 2 seconds!")
	})

	// Start app
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func update(dt float64, app *fine.App) {
	gopher.Position.Y += dt * gopherFallSpeed

	// Collision
	collision := gopher.Collide()
	if collision.TopPushY != 0 || collision.BottomPushY != 0 || collision.ForwardPushX != 0 || collision.BackwardPushX != 0 {
		// The entity is colliding, find the smallest collisions
		minX := math.Min(math.Abs(collision.ForwardPushX), math.Abs(collision.BackwardPushX))
		minY := math.Min(math.Abs(collision.TopPushY), math.Abs(collision.BottomPushY))

		if minX < minY {
			// The entity needs be pushed out on the X axis, figure out what side
			if math.Abs(collision.BackwardPushX) < math.Abs(collision.ForwardPushX) {
				gopher.Position.X += collision.BackwardPushX
			} else {
				gopher.Position.X += collision.ForwardPushX
			}
		} else {
			// The entity needs to be pushed on the Y axis, figure out what side
			if math.Abs(collision.TopPushY) < math.Abs(collision.BottomPushY) {
				gopher.Position.Y += collision.TopPushY
			} else {
				gopher.Position.Y += collision.BottomPushY
				gopherFallSpeed = 1
			}
		}
	}

	gopher.Position.Y++ // Because the gopher has an empty pixel when scaled -.-

	// Walking
	isWalking := false
	if app.IsKeyDown(fine.KEY_a) {
		gopherVelocity = 700
		gopherDirection = "left"
		isWalking = true
	} else if app.IsKeyDown(fine.KEY_d) {
		gopherVelocity = 700
		gopherDirection = "right"
		isWalking = true
	} else if gopherVelocity > 0 {
		gopherVelocity -= dt * 5000
	} else if gopherVelocity < 0 {
		gopherVelocity = 0
	}

	if app.IsKeyJustDown(fine.KEY_w) {
		walkStartTime = app.Time
	}

	switch gopherDirection {
	case "left":
		gopher.Position.X -= dt * gopherVelocity
		gopher.FlipNone()
	case "right":
		gopher.Position.X += dt * gopherVelocity
		gopher.FlipHorizontal()
	}

	if isWalking {
		gopher.Scale.Y = ((math.Sin((12*(app.Time-walkStartTime)+3)*math.Pi/2) + 1) / 2) + 5
	} else if gopher.Scale.Y > 5 {
		gopher.Scale.Y -= dt * 2
	} else if gopher.Scale.Y < 5 {
		gopher.Scale.Y = 5
	}
	if !isWalking {
		gopher.Scale.Y = ((math.Sin((2*(app.Time)+3)*math.Pi/2) + 1) / 8) + 5
	}

	// Falling
	if collision.IsFalling {
		gopherFallSpeed += dt * 5000
	}

	if gopherFallSpeed > 1500 {
		gopherFallSpeed = 1500
	}

	// Jumping
	if app.IsKeyJustDown(fine.KEY_SPACE) || app.IsKeyJustDown(fine.KEY_w) {
		gopherFallSpeed = -1000
		if gopher.Scale.X > 2 {
			gopher.Scale.X -= 1
		}
	}

	if gopher.Scale.X < 5 {
		gopher.Scale.X += dt * 5
	} else if gopher.Scale.X > 5 {
		gopher.Scale.X = 5
	}

	// Camera
	halfWidth, halfHeight := float64(app.Width)/2, float64(app.Height)/2
	camDistanceX := app.Camera.Position.X + halfWidth - gopher.Position.X + gopher.Width
	camDistanceY := app.Camera.Position.Y + halfHeight - gopher.Position.Y + gopher.Height
	isCameraMoving := false
	isMovingOnY := false

	if camDistanceX <= halfWidth-200 {
		app.Camera.Position.X += dt * cameraVelocity
		isCameraMoving = true
	}

	if camDistanceX >= halfWidth*2-200 {
		app.Camera.Position.X -= dt * cameraVelocity
		isCameraMoving = true
	}

	if camDistanceY <= 300 {
		app.Camera.Position.Y += dt * cameraVelocity
		isCameraMoving = true
		isMovingOnY = true
	}

	if camDistanceY >= halfHeight*2 {
		app.Camera.Position.Y -= dt * cameraVelocity
		isCameraMoving = true
		isMovingOnY = true
	}

	if !isCameraMoving {
		cameraVelocity = 0
	} else {
		cameraVelocity += dt * 2000
	}

	if !isMovingOnY && cameraVelocity > gopherVelocity {
		cameraVelocity = gopherVelocity
	} else if isMovingOnY && cameraVelocity > gopherFallSpeed {
		cameraVelocity = -gopherFallSpeed
	}
}

func createEntities(app *fine.App) {
	gopherSprite, err := app.NewSpriteFromData(gopherSpriteData)
	if err != nil {
		panic(err)
	}

	gopher = app.Entity(fine.NewVec2(1280/4, 720/4)).
		SetTexture(gopherSprite).
		SetScale(fine.NewVec2(5, 5)).
		SetPivotCentered(true)

	app.ClearColor = color.RGBA{54, 54, 54, 255}
}

func createLevel(app *fine.App) {
	// Floor
	app.Rect(fine.NewVec2(0, 700), 2000, 100, color.RGBA{255, 255, 255, 255}, true)

	// Box
	app.Rect(fine.NewVec2(100, 600), 50, 50, color.RGBA{128, 128, 128, 255}, true)
}
