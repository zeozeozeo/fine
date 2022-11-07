// An example of playing audio.
package main

import (
	_ "embed"
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/zeozeozeo/fine"
)

//go:embed enigma.mod
var songData []byte
var audio *fine.Audio
var isPlaying bool
var rects []*fine.Entity

func main() {
	app := fine.NewApp("Audio", 1280, 720)
	app.SetUpdateFunc(update)

	fmt.Println("rendering audio...")
	start := time.Now()

	var err error
	audio, err = app.LoadAudioFromData(songData, fine.AUDIO_MOD)
	if err != nil {
		panic(err)
	}
	fmt.Printf("finished in %s (duration: %fmin)\n", time.Since(start), audio.Duration()/60)

	for i := float64(-15); i < 15; i++ {
		// TODO: Update app.Width and app.Height on every frame
		rect := app.Rect(fine.NewVec2(60*i, 8), 50, 50, color.RGBA{255, 255, 255, 255}, true)
		rects = append(rects, rect)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func update(dt float64, app *fine.App) {
	if !isPlaying {
		audio.Play() // Start the song

		// Close the app when the song is over
		app.After(audio.Duration(), func(app *fine.App) {
			app.Close()
		})

		isPlaying = true
	}

	for idx, rect := range rects {
		rect.Position.Y = (math.Sin(app.Time*3+(float64(idx)+1)) * 150)
	}
	app.Camera.Position.X = math.Sin(app.Time*2) * 150
	app.Camera.Position.Y = math.Sin(app.Time*4)*150 + 50
}
