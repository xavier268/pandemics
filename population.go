package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

// State of a Person
type State int

// Possible States
const (
	StateLive State = iota
	StateTouched
	StateDead
	StateCured
	StateMarker
)

// StateColor provides the Color for the State
func StateColor(s State) color.RGBA {
	switch s {
	case StateLive:
		return colornames.Gray
	case StateCured:
		return colornames.Green
	case StateDead:
		return colornames.Red
	case StateTouched:
		return colornames.Blueviolet
	default:
		panic("State is not recognized")
	}
}

// Person is the model for a single person
type Person struct {
	// State of the person
	State State
	// ID is unique
	ID int
	// Position in the space
	Position pixel.Vec
	// Speed is the person speed
	Speed pixel.Vec
}

// Population is a slice of Person
type Population struct {
	people    []*Person
	size      int
	bounds    pixel.Rect
	radius    float64
	speed     float64
	frames    int              // conting frames (reset every second ...)
	second    <-chan time.Time // channel counting seconds
	count     [StateMarker]int // count by state
	elapsed   int              // elapsed time running in seconds
	running   bool             // set to false to freeze the display and update
	transProb float64          // transmission probability per second of contact
	last      time.Time        // last update cycle
	dt        float64          // delta time
	deathProb float64          // probability to die, per second
	curedProb float64          // probablity to be cured, per second
}

// NewPopulation generates a new population of the provided size.
func NewPopulation() *Population {

	pop := new(Population)
	pop.size = 200
	pop.bounds = pixel.R(0, 0, 1200, 800)
	pop.radius = 10
	pop.speed = 80.   // in pixel per second
	pop.transProb = 2 // probability to transfer, per second
	pop.running = true
	pop.last = time.Now()
	pop.curedProb, pop.deathProb = 0.096, 0.004

	for i := 0; i < pop.size; i++ {
		p := new(Person)
		p.ID = i
		p.State = StateLive
		pop.count[StateLive]++
		p.Position.X = pop.bounds.Min.X + (pop.bounds.Size().X-1)*rand.Float64()
		p.Position.Y = pop.bounds.Min.Y + (pop.bounds.Size().Y-1)*rand.Float64()
		p.Speed.X = pop.speed * (2.*rand.Float64() - 1.)
		p.Speed.Y = pop.speed * (2.*rand.Float64() - 1.)
		pop.people = append(pop.people, p)
	}

	pop.second = time.Tick(time.Second)
	return pop
}

// DeltaTimeCompute compute the delat time value.
func (pop *Population) DeltaTimeCompute() {
	pop.dt = time.Since(pop.last).Seconds()
	pop.last = time.Now()
}

// Seed the population with some infected cases
func (pop *Population) Seed(nbinfected int) *Population {
	for i := 0; i < pop.size && i < nbinfected; i++ {
		pop.people[i].State = StateTouched
		pop.count[StateTouched]++
		pop.count[StateLive]--
	}
	return pop
}

// Draw the population
func (pop Population) Draw(imd *imdraw.IMDraw) {
	imd.Clear()
	for _, p := range pop.people {
		imd.Color = StateColor(p.State)
		imd.Push(p.Position)
		imd.Circle(pop.radius, 0) // radius, thickness
		imd.Color = colornames.Black
		imd.Push(p.Position)
		imd.Circle(pop.radius, 1) // black border
	}
}

// Update population since delta time dt,
// return the start of update time.
func (pop *Population) Update() {

	for _, p := range pop.people {

		if p.State == StateTouched && rand.Float64() < pop.dt*pop.deathProb {
			pop.count[StateTouched]--
			pop.count[StateDead]++
			p.State = StateDead
		}
		if p.State == StateTouched && rand.Float64() < pop.dt*pop.curedProb {
			pop.count[StateTouched]--
			pop.count[StateCured]++
			p.State = StateCured
		}
		if p.State != StateDead {
			p.Position.X = p.Position.X + p.Speed.X*pop.dt
			p.Position.Y = p.Position.Y + p.Speed.Y*pop.dt
			pop.Reframe(p)
		}
	}
}

// Reframe handles the Person position to keep it inside the frame
func (pop *Population) Reframe(p *Person) {
	if pop.bounds.Contains(p.Position) {
		return
	}
	switch {
	case p.Position.X >= pop.bounds.Max.X:
		p.Position.X = pop.bounds.Max.X
		p.Speed.X = -p.Speed.X
	case p.Position.X <= pop.bounds.Min.X:
		p.Position.X = pop.bounds.Min.X
		p.Speed.X = -p.Speed.X
	case p.Position.Y >= pop.bounds.Max.Y:
		p.Position.Y = pop.bounds.Max.Y
		p.Speed.Y = -p.Speed.Y
	case p.Position.Y <= pop.bounds.Min.Y:
		p.Position.Y = pop.bounds.Min.Y
		p.Speed.Y = -p.Speed.Y
	default:
		panic("invalid switch case in Reframe")
	}
}

// CollisionDetect detects and triggers collision handling function.
func (pop *Population) CollisionDetect(hdlr func(p *Population, p1, p2 *Person)) {
	if hdlr == nil {
		return
	}
	for i := 0; i < pop.size; i++ {
		for j := 0; j < i; j++ {
			if d2(pop.people[i], pop.people[j]) <= 4*pop.radius*pop.radius {
				hdlr(pop, pop.people[i], pop.people[j])
			}
		}
	}
}

// d2 is the square distance between centers
func d2(p1, p2 *Person) float64 {
	x := p1.Position.X - p2.Position.X
	y := p1.Position.Y - p2.Position.Y
	return x*x + y*y
}

// handlCollision is a basic collision handler
func handlCollision(pop *Population, p1, p2 *Person) {
	if rand.Float64() < pop.transProb*pop.dt {
		switch {
		case p1.State == StateLive && p2.State == StateTouched:
			p1.State = StateTouched
			pop.count[StateLive]--
			pop.count[StateTouched]++
		case p2.State == StateLive && p1.State == StateTouched:
			p2.State = StateTouched
			pop.count[StateLive]--
			pop.count[StateTouched]++
		}
	}
}
