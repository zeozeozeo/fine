// Creates a simple blank window.
package main

import "github.com/zeozeozeo/fine"

func main() {
	app := fine.NewApp("Window title", 1280, 720) // Window title, width, height
	// app.AddWindowFlags(fine.WINDOW_RESIZABLE) - makes the window resizeable
	// See all of the window flags in the documentation

	// Start the application
	if err := app.Run(); err != nil {
		panic(err)
	}
}
