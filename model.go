package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

// Possible State of a Person
var (
	StateLive    = pixel.RGB(0.2, 0.2, 0.2)
	StateTouched = pixel.RGB(0.2, 0.4, 0.4)
	StateDead    = pixel.RGB(0, 0, 0)
	StateCured   = pixel.RGB(1, 0, 0)
)

// Person is the model for a single person
type Person struct {
	// State of the person, represented by a color
	State pixel.RGBA
	// ID is unique
	ID int
	// Position in the space
	Position pixel.Vec
	// Speed is the person speed
	Speed pixel.Vec
}

// Population is a slice of Person
type Population struct {
	people []Person
	size   int
	bounds pixel.Rect
	radius float64
	speed  float64
	frames int              // conting frames (reset every second ...)
	second <-chan time.Time // channel counting seconds
}

// NewPopulation generates a new population of the provided size.
func NewPopulation(nb int, bounds pixel.Rect) *Population {

	pop := new(Population)
	pop.size = nb
	pop.bounds = bounds
	pop.radius = 10
	pop.speed = 20.

	pop.people = make([]Person, nb, nb)
	for i := range pop.people {
		pop.people[i].ID = i
		pop.people[i].State = StateLive
		pop.people[i].Position.X = bounds.Size().X * rand.Float64()
		pop.people[i].Position.Y = bounds.Size().Y * rand.Float64()
		pop.people[i].Speed.X = pop.speed * (2*rand.Float64() - 1)
		pop.people[i].Speed.Y = pop.speed * (2*rand.Float64() - 1)
	}

	pop.second = time.Tick(time.Second)
	return pop
}

// Draw the population
func (pop Population) Draw(imd *imdraw.IMDraw) {
	for _, p := range pop.people {
		imd.Color = p.State
		imd.Push(p.Position)
		imd.Circle(pop.radius, 1) // radius, thickness
	}
}

// Update population since last update,
// return the start of update time.
func (pop *Population) Update(last time.Time) time.Time {
	dt := time.Since(last).Seconds()
	now := time.Now()
	for i, p := range pop.people {
		p.Position.X = math.Mod((p.Position.X + p.Speed.X*dt), pop.bounds.Size().X)
		p.Position.Y = math.Mod((p.Position.Y + p.Speed.Y*dt), pop.bounds.Size().Y)
		pop.people[i] = p
	}
	return now
}
