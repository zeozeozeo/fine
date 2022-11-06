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
	"github.com/zeozeozeo/gomodplay/pkg/mod"
)

type AudioFormat int

const (
	AUDIO_MP3  AudioFormat = 0 // .mp3 files.
	AUDIO_FLAC AudioFormat = 1 // .flac files.
	AUDIO_WAV  AudioFormat = 2 // .wav files.
	AUDIO_OGG  AudioFormat = 3 // .ogg files.
	AUDIO_MOD  AudioFormat = 4 // Amiga .mod files.
)

var (
	errUnsupportedFormat error = errors.New("unsupported audio format")
)

type Audio struct {
	Buffer     *beep.Buffer // Audio buffer.
	LastPlayed float64      // The time when the audio was last played.
	app        *App
}

// Loads an audio file from a ReadCloser. If the audio sample rate doesn't match
// the configured app sample rate, the file will be resampled.
func (app *App) LoadAudio(readCloser io.ReadCloser, inputFormat AudioFormat) (*Audio, error) {
	var streamer beep.StreamSeekCloser
	var modStreamer beep.Streamer
	isMod := false
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
	case AUDIO_MOD:
		modStreamer, format, err = app.loadMod(readCloser)
		isMod = true
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
		if !isMod {
			buffer.Append(streamer)
		} else {
			buffer.Append(modStreamer)
		}
	}

	audio := &Audio{
		Buffer: buffer,
		app:    app,
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
	audio.LastPlayed = audio.app.Time
	speaker.Play(audio.GetStream())
}

// Stops all playing audio.
func (app *App) StopAudio() {
	speaker.Clear()
}

// Returns the duration of the audio buffer in seconds.
func (audio *Audio) Duration() float64 {
	return float64(audio.Buffer.Len()) / float64(audio.Buffer.Format().SampleRate)
}

// Returns if the audio has stopped playing or not.
func (audio *Audio) Ended() bool {
	return (audio.app.Time - audio.LastPlayed) >= audio.Duration()
}

type modStreamer struct {
	player *mod.Player
}

func (s modStreamer) Stream(samples [][2]float64) (n int, ok bool) {
	if s.player.State.SongHasEnded || s.player.State.HasLooped {
		return 0, false
	}
	for idx := range samples {
		left, right := s.player.NextSample()
		samples[idx][0] = float64(left)
		samples[idx][1] = float64(right)
	}
	return len(samples), true
}

func (s modStreamer) Err() error {
	return nil
}

func (app *App) loadMod(readCloser io.ReadCloser) (beep.Streamer, beep.Format, error) {
	streamer := modStreamer{}
	streamer.player = mod.NewModPlayer(uint32(app.SampleRate))

	err := streamer.player.LoadModFile(readCloser)
	if err != nil {
		return nil, beep.Format{}, err
	}
	streamer.player.Play()

	return streamer, beep.Format{
		SampleRate:  beep.SampleRate(app.SampleRate),
		NumChannels: 2,
		Precision:   4,
	}, nil
}
