package fine

import (
	"bytes"
	"errors"
	"io"
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/flac"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"
)

type AudioFormat int

const (
	AUDIO_MP3  AudioFormat = 0 // .mp3 files.
	AUDIO_FLAC AudioFormat = 1 // .flac files.
	AUDIO_WAV  AudioFormat = 2 // .wav files.
	AUDIO_OGG  AudioFormat = 3 // .ogg files.
)

var (
	errUnsupportedFormat error = errors.New("unsupported audio format")
)

type Audio struct {
	Buffer *beep.Buffer // Audio buffer.
	Format beep.Format  // The audio format.
}

// Loads an audio file from a ReadCloser. If the audio sample rate doesn't match
// the configured app sample rate, the file will be resampled.
func (app *App) LoadAudio(readCloser io.ReadCloser, inputFormat AudioFormat) (*Audio, error) {
	var streamer beep.StreamSeekCloser
	var format beep.Format
	var err error

	switch inputFormat {
	case AUDIO_MP3:
		streamer, format, err = mp3.Decode(readCloser)
	case AUDIO_FLAC:
		streamer, format, err = flac.Decode(readCloser)
	case AUDIO_WAV:
		streamer, format, err = wav.Decode(readCloser)
	case AUDIO_OGG:
		streamer, format, err = vorbis.Decode(readCloser)
	default:
		return nil, errUnsupportedFormat
	}

	if err != nil {
		return nil, err
	}

	// Resample if needed
	buffer := beep.NewBuffer(format)

	if int(format.SampleRate) != app.SampleRate {
		resampler := beep.Resample(
			app.ResamplingQuality,
			format.SampleRate,
			beep.SampleRate(app.SampleRate),
			streamer,
		)
		buffer.Append(resampler)
	} else {
		buffer.Append(streamer)
	}

	audio := &Audio{
		Buffer: buffer,
		Format: format,
	}

	return audio, nil
}

// Loads audio from bytes.
func (app *App) LoadAudioFromData(data []byte, inputFormat AudioFormat) (*Audio, error) {
	buf := bytes.NewBuffer(data)
	return app.LoadAudio(io.NopCloser(buf), inputFormat)
}

func (app *App) LoadAudioFromReader(reader io.Reader, inputFormat AudioFormat) (*Audio, error) {
	return app.LoadAudio(io.NopCloser(reader), inputFormat)
}

func (app *App) initAudio() {
	beepSampleRate := beep.SampleRate(app.SampleRate)
	// TODO: Custom buffer sizes
	err := speaker.Init(beepSampleRate, beepSampleRate.N(time.Second/10))
	if err != nil {
		log.Printf("[warn] failed to initialize audio speaker: %s", err)
	}
}

// Returns a brand new audio stream from the audio buffer.
func (audio *Audio) GetStream() beep.StreamSeeker {
	return audio.Buffer.Streamer(0, audio.Buffer.Len())
}

// Starts playing the audio. This is an asynchronous call.
func (audio *Audio) Play() {
	speaker.Play(audio.GetStream())
}

// Stops all playing audio.
func (app *App) StopAudio() {
	speaker.Clear()
}
