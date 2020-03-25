package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// The main function needs to delegate to run()
// to ensure thread consistency.
func main() {
	pixelgl.Run(run)
}

// run is where all the code will be fired from.
func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Animation",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	// create the window
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	pop := NewPopulation(500, win.Bounds()).Seed(1)

	imd := imdraw.New(nil)

	for !win.Closed() {

		pop.DeltaTimeCompute()

		// Only move if running  ...
		if pop.running {
			pop.Update()
			pop.CollisionDetect(handlCollision)
			pop.Draw(imd)

			win.Clear(colornames.Skyblue)
			imd.Draw(win)
			FPSDisplay(win, pop)
		}

		win.Update() // keep refreshing to ensure windows remain responsive !

		// should we stop running ?
		if pop.running && pop.count[StateLive] <= 0 {
			pop.running = false
		}
	}
}

// FPSDisplay dispaly the FPS in the title bar.
func FPSDisplay(win *pixelgl.Window, pop *Population) {
	// counting FPS and displaying in title
	select {
	case <-pop.second:
		win.SetTitle(fmt.Sprintf("Time %ds, FPS: %d - Live %d, Touched %d, Cured %d, Dead %d",
			pop.elapsed,
			pop.frames,
			pop.count[StateLive],
			pop.count[StateTouched],
			pop.count[StateCured],
			pop.count[StateDead]))
		pop.frames = 0
		pop.elapsed++
	default:
		pop.frames++
	}
}
