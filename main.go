package main

import (
	"fmt"
	"time"

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

	// create the window
	win, err := pixelgl.NewWindow(NewConfig())
	if err != nil {
		panic(err)
	}
	win.SetSmooth(true)

	pop := NewPopulation(100, win.Bounds())

	imd := imdraw.New(nil)
	pop.Draw(imd)

	last := time.Now()
	for !win.Closed() {

		last = pop.Update(last)

		imd := imdraw.New(nil)
		pop.Draw(imd)

		win.Clear(colornames.Skyblue)
		imd.Draw(win)

		win.Update()

		// counting FPS and displaying in title
		select {
		case <-pop.second:
			win.SetTitle(fmt.Sprintf("FPS: %d", pop.frames))
			pop.frames = 0
		default:
			pop.frames++
		}
	}
}
