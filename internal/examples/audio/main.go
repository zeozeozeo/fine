// An example of playing audio.
package main

import (
	_ "embed"
	"math"

	"github.com/zeozeozeo/fine"
)

//go:embed lowsi_jeti_maah.mod
var songData []byte
var audio *fine.Audio
var isPlaying bool

func main() {
	app := fine.NewApp("Audio", 1280, 720)
	app.SetUpdateFunc(update)

	var err error
	audio, err = app.LoadAudioFromData(songData, fine.AUDIO_MOD)
	if err != nil {
		panic(err)
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

	c := uint8(math.Sin(app.Time) * 255)
	app.ClearColor.R = c
	app.ClearColor.G = c
	app.ClearColor.B = c
}
