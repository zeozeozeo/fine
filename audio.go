package fine

import (
	"bytes"
	"errors"
	"io"
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
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
	AUDIO_MOD  AudioFormat = 4 // .mod files.
)

var (
	errUnsupportedFormat error = errors.New("unsupported audio format")
)

type Audio struct {
	Buffer     *beep.Buffer // Audio buffer.
	LastPlayed float64      // The time when the audio was last played.

	// Volume of the audio. Use SetVolume to change the volume.
	// The volume is applied exponentially. Volume of 0 means no change.
	// Set Silent if you want to make the audio silent.
	Volume float64

	// Natural Base for the volume exponentiation. The default is 2.
	// In order to adjust volume along decibells, pick 10 as the Base and set Volume to dB/10.
	Base         float64
	Silent       bool // Specifies whether the audio should be silent.
	app          *App
	volumeEffect *effects.Volume
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
		// TODO: OGG support is poor, some audio sounds pitched. Find a better way to do it.
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
		Base:   2,
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

	err := speaker.Init(beepSampleRate, beepSampleRate.N(time.Duration(app.BufferNs)*time.Nanosecond))
	if err != nil {
		log.Printf("[warn] failed to initialize audio speaker: %s", err)
	}
}

// Returns a brand new audio stream from the audio buffer.
func (audio *Audio) GetStream() beep.Streamer {
	audio.volumeEffect = &effects.Volume{
		Streamer: audio.Buffer.Streamer(0, audio.Buffer.Len()),
		Volume:   audio.Volume,
		Base:     audio.Base,
		Silent:   audio.Silent,
	}
	return audio.volumeEffect
}

// Starts playing the audio. This is an asynchronous call.
func (audio *Audio) Play() *Audio {
	audio.LastPlayed = audio.app.Time
	speaker.Play(audio.GetStream())
	return audio
}

func (audio *Audio) SetVolume(vol float64) *Audio {
	audio.Volume = vol
	audio.volumeEffect.Volume = vol
	return audio
}

// Stops all playing audio.
func (app *App) StopAudio() *App {
	speaker.Clear()
	return app
}

// Sets the amount of nanoseconds in the audio buffer. Default: 48000000 (48ms)
func (app *App) SetAudioBufferNs(nanoseconds int64) *App {
	app.BufferNs = nanoseconds
	app.initAudio()
	return app
}

// Sets the sample rate of the audio. Default: 44100.
func (app *App) SetAudioSampleRate(sampleRate int) *App {
	app.SampleRate = sampleRate
	app.initAudio()
	return app
}

// Returns the duration of the audio buffer in seconds.
func (audio *Audio) Duration() float64 {
	return float64(audio.Buffer.Len()) / float64(audio.Buffer.Format().SampleRate)
}

// Returns if the audio has stopped playing or not.
func (audio *Audio) Ended() bool {
	return (audio.app.Time - audio.LastPlayed) >= audio.Duration()
}

// .mod file playback
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
