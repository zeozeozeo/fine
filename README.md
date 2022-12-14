<p align="center">
  <a href="https://godoc.org/github.com/zeozeozeo/fine"><img src="https://pkg.go.dev/badge/github.com/zeozeozeo/fine" alt="Documentation"></img></a>
  <a href="https://goreportcard.com/report/github.com/zeozeozeo/fine"><img src="https://goreportcard.com/badge/github.com/zeozeozeo/fine" alt="Go Report Card"/></a>
  <a href="https://unlicense.org/"><img src="https://img.shields.io/github/license/zeozeozeo/fine" alt="License"/></a>

</p>
<p><h1 align="center">Fine - a super easy to use 2D game engine for Go.</h1></p>

Fine is a cross-platform OpenGL 2D engine (with [SDL2](https://github.com/veandco/go-sdl2)) written with Go. It has audio support, text rendering, collisons and much more.

![A platformer character](https://github.com/zeozeozeo/fine-examples/raw/main/platformer/platformer.gif)

# Getting Started

## [First of all, install the SDL2 dependencies.](https://github.com/veandco/go-sdl2#requirements)

The SDL2 `mixer` package is not necessary, Fine does not use it.

When distributing on Windows, make sure to move the SDL2 dynamic link libraries in the same folder as the executable.

## Install and import the package

`go get -v github.com/zeozeozeo/fine`

```go
package main

import (
    "github.com/zeozeozeo/fine"
)
```

## Start a new window

```go
func main() {
    // The arguments are the window title, window width and window height
	app := fine.NewApp("Title", 1280, 720)

	// Start the application
	if err := app.Run(); err != nil {
		panic(err)
	}
}
```

## You can see more examples (audio, text, platformer character, camera) in [the examples repository](https://github.com/zeozeozeo/fine-examples).

# Features

-   **Easy to use.**
-   **Simple entity system.** You can think of any game object as an entity.
-   **Cross-platform.** Your games will work on most platforms.
-   **Supports audio** (uses [Beep](https://github.com/faiface/beep) for mixing and audio playback).
    -   [MP3](https://en.wikipedia.org/wiki/MP3)
    -   [WAVE](https://en.wikipedia.org/wiki/WAV)
    -   [OGG (Vorbis)](https://en.wikipedia.org/wiki/Ogg)
    -   [FLAC](https://en.wikipedia.org/wiki/FLAC)
    -   [MOD](<https://en.wikipedia.org/wiki/MOD_(file_format)>)
-   **Text rendering.** You can import [TrueType](https://en.wikipedia.org/wiki/TrueType) fonts.

# TODO

-   More shapes
-   WebAssembly
-   Test it on more platforms
-   and more...
