// A simple example of making particles.
package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/zeozeozeo/fine"
)

type Particle struct {
	Created  float64
	Lifetime float64
	Speed    fine.Vec2
	Entity   *fine.Entity
}

var (
	particles        []*Particle
	spawnOrigin      fine.Vec2 = fine.NewVec2(0, 300)
	particleSize     float64   = 10
	particleLifetime float64   = 1
	lastSpawn        float64
	spawnDelay       float64
)

func main() {
	app := fine.NewApp("Particles", 1280, 720)
	app.SetUpdateFunc(update)

	// Start the application
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func update(dt float64, app *fine.App) {
	if app.Time-lastSpawn >= spawnDelay {
		emitParticle(app)
		lastSpawn = app.Time
	}

	// Remove decayed particles
	for idx, particle := range particles {
		if app.Time-particle.Created >= particle.Lifetime {
			particle.Entity.Destroy()
			particles = append(particles[:idx], particles[idx+1:]...)
			continue
		}

		particleUpdate(dt, app, particle)
	}
}

func emitParticle(app *fine.App) {
	sizeMul := rand.Float64() + 0.5
	col := 255 - uint8(rand.Float64()*128)
	particle := app.Rect(
		spawnOrigin,
		particleSize*sizeMul,
		particleSize*sizeMul,
		color.RGBA{col, col, col, 255},
		true,
	)
	xDir := math.Sin(app.Time*2) * 1500 * (rand.Float64()) // X direction of the particle

	particles = append(particles, &Particle{
		Created:  app.Time,
		Entity:   particle,
		Speed:    fine.NewVec2(xDir, -1000),
		Lifetime: particleLifetime,
	})
}

func particleUpdate(dt float64, app *fine.App, particle *Particle) {
	particle.Entity.Position.X += particle.Speed.X * dt
	particle.Entity.Position.Y += particle.Speed.Y * dt
}
