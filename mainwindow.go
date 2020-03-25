package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

// MWindow object
type MWindow struct {
	win *pixelgl.Window
	imd *imdraw.IMDraw
	pop *Population
	err error
}

// NewMWindow constructor
func NewMWindow() *MWindow {

	mw := new(MWindow)

	cfg := pixelgl.WindowConfig{
		Title:  "Animation",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}

	// create the window
	mw.win, mw.err = pixelgl.NewWindow(cfg)
	if mw.err != nil {
		panic(mw.err)
	}

	// Prepare population window
	mw.win.SetSmooth(true)
	mw.pop = NewPopulation(100, mw.win.Bounds()).Seed(1)
	mw.imd = imdraw.New(nil)

	return mw
}

//Closed wrapper
func (mw *MWindow) Closed() bool {
	return mw.win.Closed()
}

// Update loop
func (mw *MWindow) Update() {
	mw.pop.DeltaTimeCompute()

	// Only move if running  ...
	if mw.pop.running {
		mw.pop.Update()
		mw.pop.CollisionDetect(handlCollision)
		mw.pop.Draw(mw.imd)

		mw.win.Clear(colornames.Skyblue)
		mw.imd.Draw(mw.win)
		mw.fpsDisplay()

	}

	// should we stop running ?
	if mw.pop.running && mw.pop.count[StateLive] <= 0 {
		mw.pop.running = false
	}

	// Update in anycase, to stay responsive
	mw.win.Update()
}

// fpsDisplay dispaly the FPS in the title bar,
// with the population status
func (mw *MWindow) fpsDisplay() {
	// counting FPS and displaying in title
	select {
	case <-mw.pop.second:
		mw.win.SetTitle(fmt.Sprintf("Time %ds, FPS: %d - Live %d, Touched %d, Cured %d, Dead %d",
			mw.pop.elapsed,
			mw.pop.frames,
			mw.pop.count[StateLive],
			mw.pop.count[StateTouched],
			mw.pop.count[StateCured],
			mw.pop.count[StateDead]))
		mw.pop.frames = 0
		mw.pop.elapsed++
	default:
		mw.pop.frames++
	}
}