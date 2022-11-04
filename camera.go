package fine

import "math"

type Camera struct {
	Position Vec2    // World position of the camera.
	Zoom     float64 // Bigger value = bigger objects. Default: 1.
}

// Set the world position of the camera.
func (camera *Camera) SetPosition(position Vec2) *Camera {
	camera.Position = position
	return camera
}

// Set the zoom of the camera. Bigger value = bigger objects.
// FIXME: This zooms to 0,0 instead of the center of the screen.
func (camera *Camera) SetZoom(zoom float64) *Camera {
	camera.Zoom = zoom
	return camera
}

// Convert world position to screen posititon. Returns the X and Y positions on the screen.
func (camera *Camera) WorldToScreen(worldPos Vec2) (int, int) {
	sub := worldPos.Sub(camera.Position)
	return int(math.Round(sub.X)), int(math.Round(sub.Y))
}
