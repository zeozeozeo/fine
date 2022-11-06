// An example with a platformer character and camera.
package main

import (
	_ "embed"
	"image/color"
	"math"

	"github.com/zeozeozeo/fine"
)

var (
	gopher          *fine.Entity
	gopherFallSpeed float64
	gopherVelocity  float64
	gopherDirection string // "left", "right" or ""
	walkStartTime   float64
	cameraVelocity  float64
	jumpSound       *fine.Audio
)

//go:embed gopher.png
var gopherSpriteData []byte

//go:embed jump.ogg
var gopherJumpSoundData []byte

func main() {
	app := fine.NewApp("Platformer", 1280, 720)
	app.SetUpdateFunc(update)
	// app.AddWindowFlags(fine.WINDOW_RESIZABLE) - makes the window resizeable
	createEntities(app)
	createLevel(app)

	var err error
	jumpSound, err = app.LoadAudioFromData(gopherJumpSoundData, fine.AUDIO_OGG)
	if err != nil {
		panic(err)
	}

	// Make the camera start at the gopher's position
	app.Camera.Position = gopher.Position

	// Start app
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// This function will be called every frame, and is used to update the game logic.
func update(dt float64, app *fine.App) {
	// Fall
	gopher.Position.Y += dt * gopherFallSpeed

	// Collision (TODO: Find a cleaner way to do this)
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
				gopherFallSpeed = 1 // Stop jumping
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
		// Stopped walking, start decreasing velocity
		gopherVelocity -= dt * 5000
	} else if gopherVelocity < 0 { // Make sure that the velocity is >= 0
		gopherVelocity = 0
	}

	// Check walking direction
	switch gopherDirection {
	case "left":
		gopher.Position.X -= dt * gopherVelocity
		gopher.FlipNone()
	case "right":
		gopher.Position.X += dt * gopherVelocity
		gopher.FlipHorizontal()
	}

	// The walking animation could start abruptly, so subtract the time when we started
	// walking from it.
	if app.IsKeyJustDown(fine.KEY_a) || app.IsKeyJustDown(fine.KEY_d) {
		walkStartTime = app.Time
	}

	// Walking animation. This scales a sine wave to be in the range of 0..1
	// and adds 5 to it (because the intial scale of the gopher is 5)
	if isWalking {
		gopher.Scale.Y = ((math.Sin((12*(app.Time-walkStartTime)+3)*math.Pi/2) + 1) / 2) + 5
	} else if gopher.Scale.Y > 5 {
		gopher.Scale.Y -= dt * 2
	} else if gopher.Scale.Y < 5 {
		gopher.Scale.Y = 5
	}

	// Idle animation, this is again just a scaled sine wave, but a little bit slower
	if !isWalking {
		gopher.Scale.Y = ((math.Sin((2*(app.Time)+3)*math.Pi/2) + 1) / 8) + 5
	}

	// Fall speed
	if collision.IsFalling {
		gopherFallSpeed += dt * 5000
	}

	// Make sure that the fall speed is <= 1500
	if gopherFallSpeed > 1500 {
		gopherFallSpeed = 1500
	}

	// Jumping
	if app.IsKeyJustDown(fine.KEY_SPACE) || app.IsKeyJustDown(fine.KEY_w) {
		gopherFallSpeed = -1000 // This just makes the fall velocity negative, which makes us go up.
		// Decrease X scale to animate the gopher
		if gopher.Scale.X > 2 {
			gopher.Scale.X -= 1
		}

		// Play sound effect
		jumpSound.Play()
	}

	// Jumping animation
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
	isMovingDown := false

	// Move camera to the right
	if camDistanceX <= halfWidth-200 {
		app.Camera.Position.X += dt * cameraVelocity
		isCameraMoving = true
	}

	// Move camera to the left
	if camDistanceX >= halfWidth*2-200 {
		app.Camera.Position.X -= dt * cameraVelocity
		isCameraMoving = true
	}

	// Move the camera down
	if camDistanceY <= 300 {
		app.Camera.Position.Y += dt * cameraVelocity
		isCameraMoving = true
		isMovingOnY = true
		isMovingDown = true
	}

	// Move the camera up
	if camDistanceY >= halfHeight*2 {
		app.Camera.Position.Y -= dt * cameraVelocity
		isCameraMoving = true
		isMovingOnY = true
	}

	// Increase camera velocity if it is moving
	if !isCameraMoving {
		cameraVelocity = 0
	} else {
		cameraVelocity += dt * 2000
	}

	// Make sure that the camera velocity is smaller than the walking velocity
	// if the camera is moving to the left or right, and smaller than the fall
	// velocity if the camera is moving to the top or bottom
	if !isMovingOnY && cameraVelocity > gopherVelocity {
		cameraVelocity = gopherVelocity
	} else if isMovingOnY && cameraVelocity > gopherFallSpeed {
		if isMovingDown {
			cameraVelocity = math.Abs(gopherFallSpeed)
		} else {
			cameraVelocity = -gopherFallSpeed
		}
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

	// This will be the background color
	app.ClearColor = color.RGBA{54, 54, 54, 255}
}

func createLevel(app *fine.App) {
	// Floor
	app.Rect(fine.NewVec2(0, 700), 2000, 100, color.RGBA{255, 255, 255, 255}, true)

	// Box
	app.Rect(fine.NewVec2(100, 600), 50, 50, color.RGBA{0, 255, 0, 255}, true)
}
